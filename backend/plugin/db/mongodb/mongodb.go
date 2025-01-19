// Package mongodb is the plugin for MongoDB driver.
package mongodb

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/plugin/db"
	"github.com/bytebase/bytebase/backend/resources/mongoutil"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
)

var _ db.Driver = (*Driver)(nil)

func init() {
	db.Register(storepb.Engine_MONGODB, newDriver)
}

// Driver is the MongoDB driver.
type Driver struct {
	dbBinDir      string
	connectionCtx db.ConnectionContext
	connCfg       db.ConnectionConfig
	client        *mongo.Client
	databaseName  string
}

func newDriver(dc db.DriverConfig) db.Driver {
	return &Driver{dbBinDir: dc.DbBinDir}
}

// Open opens a MongoDB driver.
func (driver *Driver) Open(ctx context.Context, _ storepb.Engine, connCfg db.ConnectionConfig) (db.Driver, error) {
	connectionURI := getBasicMongoDBConnectionURI(connCfg)
	opts := options.Client().ApplyURI(connectionURI)
	tlsConfig, err := connCfg.TLSConfig.GetSslConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get SSL config")
	}
	if tlsConfig != nil {
		// TODO(zp): User uses ssh tunnel?
		tlsConfig.InsecureSkipVerify = true
		opts.SetTLSConfig(tlsConfig)
	}
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create MongoDB client")
	}
	driver.client = client
	driver.connectionCtx = connCfg.ConnectionContext
	driver.connCfg = connCfg
	driver.databaseName = connCfg.Database
	return driver, nil
}

// Close closes the MongoDB driver.
func (driver *Driver) Close(ctx context.Context) error {
	if err := driver.client.Disconnect(ctx); err != nil {
		return errors.Wrap(err, "failed to disconnect MongoDB")
	}
	return nil
}

// Ping pings the database.
func (driver *Driver) Ping(ctx context.Context) error {
	if err := driver.client.Ping(ctx, nil); err != nil {
		return errors.Wrap(err, "failed to ping MongoDB")
	}
	return nil
}

// GetDB gets the database.
func (*Driver) GetDB() *sql.DB {
	return nil
}

// Execute executes a statement, always returns 0 as the number of rows affected because we execute the statement by mongosh, it's hard to catch the row effected number.
func (driver *Driver) Execute(ctx context.Context, statement string, _ db.ExecuteOptions) (int64, error) {
	connectionURI := getBasicMongoDBConnectionURI(driver.connCfg)
	// For MongoDB, we execute the statement in mongosh, which is a shell for MongoDB.
	// There are some ways to execute the statement in mongosh:
	// 1. Use the --eval option to execute the statement.
	// 2. Use the --file option to execute the statement from a file.
	// We choose the second way with the following reasons:
	// 1. The statement may too long to be executed in the command line.
	// 2. We cannot catch the error from the --eval option.
	mongoshArgs := []string{
		connectionURI,
		// DocumentDB do not support retryWrites, so we set it to false.
		"--retryWrites",
		"false",
		"--quiet",
	}

	if driver.connCfg.TLSConfig.UseSSL {
		mongoshArgs = append(mongoshArgs, "--tls")
		mongoshArgs = append(mongoshArgs, "--tlsAllowInvalidHostnames")

		uuid := uuid.New().String()
		if driver.connCfg.TLSConfig.SslCA == "" {
			mongoshArgs = append(mongoshArgs, "--tlsUseSystemCA")
		} else {
			// Write the tlsCAFile to a temporary file, and use the temporary file as the value of --tlsCAFile.
			// The reason is that the --tlsCAFile option of mongosh does not support the value of the certificate directly.
			caFileName := fmt.Sprintf("mongodb-tls-ca-%s-%s", driver.connCfg.Database, uuid)
			defer func() {
				// While error occurred in mongosh, the temporary file may not created, so we ignore the error here.
				_ = os.Remove(caFileName)
			}()
			if err := os.WriteFile(caFileName, []byte(driver.connCfg.TLSConfig.SslCA), 0400); err != nil {
				return 0, errors.Wrap(err, "failed to write tlsCAFile to temporary file")
			}
			mongoshArgs = append(mongoshArgs, "--tlsCAFile", caFileName)
		}

		if driver.connCfg.TLSConfig.SslKey != "" && driver.connCfg.TLSConfig.SslCert != "" {
			clientCertName := fmt.Sprintf("mongodb-tls-client-cert-%s-%s", driver.connCfg.Database, uuid)
			defer func() {
				// While error occurred in mongosh, the temporary file may not created, so we ignore the error here.
				_ = os.Remove(clientCertName)
			}()
			var sb strings.Builder
			if _, err := sb.WriteString(driver.connCfg.TLSConfig.SslKey); err != nil {
				return 0, errors.Wrapf(err, "failed to write ssl key into string builder")
			}
			if _, err := sb.WriteString("\n"); err != nil {
				return 0, errors.Wrapf(err, "failed to write new line into string builder")
			}
			if _, err := sb.WriteString(driver.connCfg.TLSConfig.SslCert); err != nil {
				return 0, errors.Wrapf(err, "failed to write ssl cert into string builder")
			}
			if err := os.WriteFile(clientCertName, []byte(sb.String()), 0400); err != nil {
				return 0, errors.Wrap(err, "failed to write tlsCAFile to temporary file")
			}
			mongoshArgs = append(mongoshArgs, "--tlsCertificateKeyFile", clientCertName)
		}
	}

	// First, we create a temporary file to store the statement.
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "mongodb-statement")
	if err != nil {
		return 0, errors.Wrap(err, "failed to create temporary file")
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.WriteString(statement); err != nil {
		return 0, errors.Wrap(err, "failed to write statement to temporary file")
	}
	if err := tempFile.Close(); err != nil {
		return 0, errors.Wrap(err, "failed to close temporary file")
	}
	mongoshArgs = append(mongoshArgs, "--file", tempFile.Name())

	mongoshCmd := exec.CommandContext(ctx, mongoutil.GetMongoshPath(driver.dbBinDir), mongoshArgs...)
	var errContent bytes.Buffer
	mongoshCmd.Stderr = &errContent
	if err := mongoshCmd.Run(); err != nil {
		return 0, errors.Wrapf(err, "failed to execute statement in mongosh: %s", errContent.String())
	}
	return 0, nil
}

// Dump dumps the database.
func (*Driver) Dump(_ context.Context, _ io.Writer, _ *storepb.DatabaseSchemaMetadata) error {
	return nil
}

// getBasicMongoDBConnectionURI returns the basic MongoDB connection URI, the following fields are excluded:
// - TLS related
// https://www.mongodb.com/docs/manual/reference/connection-string/
func getBasicMongoDBConnectionURI(connConfig db.ConnectionConfig) string {
	u := &url.URL{
		Scheme: "mongodb",
		// In RFC, there can be no tailing slash('/') in the path if the path is empty and the query is not empty.
		// For mongosh, it can handle this case correctly, but for driver, it will throw the error likes "error parsing uri: must have a / before the query ?".
		Path: "/",
	}
	if connConfig.SRV {
		u.Scheme = "mongodb+srv"
	}
	if connConfig.Username != "" {
		u.User = url.UserPassword(connConfig.Username, connConfig.Password)
	}
	u.Host = connConfig.Host
	if connConfig.Port != "" {
		u.Host = fmt.Sprintf("%s:%s", u.Host, connConfig.Port)
	}
	for _, additionalAddress := range connConfig.AdditionalAddresses {
		address := additionalAddress.Host
		if additionalAddress.Port != "" {
			address = fmt.Sprintf("%s:%s", address, additionalAddress.Port)
		}
		u.Host = fmt.Sprintf("%s,%s", u.Host, address)
	}
	if connConfig.Database != "" {
		u.Path = connConfig.Database
	}
	authDatabase := "admin"
	if connConfig.AuthenticationDatabase != "" {
		authDatabase = connConfig.AuthenticationDatabase
	}

	values := u.Query()
	values.Add("authSource", authDatabase)
	if connConfig.ReplicaSet != "" {
		values.Add("replicaSet", connConfig.ReplicaSet)
	}
	values.Add("appName", "bytebase")
	if connConfig.DirectConnection {
		values.Add("directConnection", "true")
	}
	u.RawQuery = values.Encode()

	return u.String()
}

// QueryConn queries a SQL statement in a given connection.
func (driver *Driver) QueryConn(ctx context.Context, _ *sql.Conn, statement string, queryContext db.QueryContext) ([]*v1pb.QueryResult, error) {
	if queryContext.Explain {
		return nil, errors.New("MongoDB does not support EXPLAIN")
	}

	statement = strings.Trim(statement, " \t\n\r\f;")
	simpleStatement := isMongoStatement(statement)
	startTime := time.Now()
	connectionURI := getBasicMongoDBConnectionURI(driver.connCfg)
	// For MongoDB query, we execute the statement in mongosh with flag --eval for the following reasons:
	// 1. Query always short, so it's safe to execute in the command line.
	// 2. We cannot catch the output if we use the --file option.

	evalArg := statement
	if simpleStatement {
		limit := ""
		if queryContext.Limit > 0 {
			limit = fmt.Sprintf(".slice(0, %d)", queryContext.Limit)
		}
		evalArg = fmt.Sprintf(`a = %s; if (typeof a.toArray === 'function') {a.toArray()%s;} else {a;}`, statement, limit)
	}
	// We will use single quotes for the evalArg, so we need to escape the single quotes in the statement.
	evalArg = strings.ReplaceAll(evalArg, `'`, `'"'`)
	evalArg = fmt.Sprintf(`'%s'`, evalArg)

	mongoshArgs := []string{
		mongoutil.GetMongoshPath(driver.dbBinDir),
		// quote the connectionURI because we execute the mongosh via sh, and the multi-queries part contains '&', which will be translated to the background process.
		fmt.Sprintf(`"%s"`, connectionURI),
		"--quiet",
		"--json",
		"canonical",
		"--eval",
		evalArg,
		// DocumentDB do not support retryWrites, so we set it to false.
		"--retryWrites",
		"false",
	}

	if driver.connCfg.TLSConfig.UseSSL {
		mongoshArgs = append(mongoshArgs, "--tls")
		mongoshArgs = append(mongoshArgs, "--tlsAllowInvalidHostnames")

		uuid := uuid.New().String()
		if driver.connCfg.TLSConfig.SslCA == "" {
			mongoshArgs = append(mongoshArgs, "--tlsUseSystemCA")
		} else {
			// Write the tlsCAFile to a temporary file, and use the temporary file as the value of --tlsCAFile.
			// The reason is that the --tlsCAFile option of mongosh does not support the value of the certificate directly.
			caFileName := fmt.Sprintf("mongodb-tls-ca-%s-%s", driver.connCfg.Database, uuid)
			defer func() {
				// While error occurred in mongosh, the temporary file may not created, so we ignore the error here.
				_ = os.Remove(caFileName)
			}()
			if err := os.WriteFile(caFileName, []byte(driver.connCfg.TLSConfig.SslCA), 0400); err != nil {
				return nil, errors.Wrap(err, "failed to write tlsCAFile to temporary file")
			}
			mongoshArgs = append(mongoshArgs, "--tlsCAFile", caFileName)
		}

		if driver.connCfg.TLSConfig.SslKey != "" && driver.connCfg.TLSConfig.SslCert != "" {
			clientCertName := fmt.Sprintf("mongodb-tls-client-cert-%s-%s", driver.connCfg.Database, uuid)
			defer func() {
				// While error occurred in mongosh, the temporary file may not created, so we ignore the error here.
				_ = os.Remove(clientCertName)
			}()
			var sb strings.Builder
			if _, err := sb.WriteString(driver.connCfg.TLSConfig.SslKey); err != nil {
				return nil, errors.Wrapf(err, "failed to write ssl key into string builder")
			}
			if _, err := sb.WriteString("\n"); err != nil {
				return nil, errors.Wrapf(err, "failed to write new line into string builder")
			}
			if _, err := sb.WriteString(driver.connCfg.TLSConfig.SslCert); err != nil {
				return nil, errors.Wrapf(err, "failed to write ssl cert into string builder")
			}
			if err := os.WriteFile(clientCertName, []byte(sb.String()), 0400); err != nil {
				return nil, errors.Wrap(err, "failed to write tlsCAFile to temporary file")
			}
			mongoshArgs = append(mongoshArgs, "--tlsCertificateKeyFile", clientCertName)
		}
	}

	queryResultFileName := fmt.Sprintf("mongodb-query-%s-%s", driver.connCfg.Database, uuid.New().String())
	defer func() {
		// While error occurred in mongosh, the temporary file may not created, so we ignore the error here.
		_ = os.Remove(queryResultFileName)
	}()
	mongoshArgs = append(mongoshArgs, ">", queryResultFileName)

	shellArgs := []string{
		"-c",
		strings.Join(mongoshArgs, " "),
	}
	shCmd := exec.CommandContext(ctx, "sh", shellArgs...)
	var errContent bytes.Buffer
	var outContent bytes.Buffer
	shCmd.Stderr = &errContent
	shCmd.Stdout = &outContent
	if err := shCmd.Run(); err != nil {
		return nil, errors.Wrapf(err, "failed to execute statement in mongosh: \n stdout: %s\n stderr: %s", outContent.String(), errContent.String())
	}

	f, err := os.OpenFile(queryResultFileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", queryResultFileName)
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if int64(fileInfo.Size()) > driver.connCfg.MaximumSQLResultSize {
		return []*v1pb.QueryResult{{
			Latency:   durationpb.New(time.Since(startTime)),
			Statement: statement,
			Error:     common.FormatMaximumSQLResultSizeMessage(driver.connCfg.MaximumSQLResultSize),
		}}, nil
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file: %s", queryResultFileName)
	}

	if simpleStatement {
		// We make best-effort attempt to parse the content and fallback to single bulk result on failure.
		result, err := getSimpleStatementResult(content)
		if err != nil {
			slog.Error("failed to get simple statement result", slog.String("content", string(content)), log.BBError(err))
		} else {
			result.Latency = durationpb.New(time.Since(startTime))
			result.Statement = statement
			return []*v1pb.QueryResult{result}, nil
		}
	}

	return []*v1pb.QueryResult{{
		ColumnNames:     []string{"result"},
		ColumnTypeNames: []string{"TEXT"},
		Rows: []*v1pb.QueryRow{{
			Values: []*v1pb.RowValue{{
				Kind: &v1pb.RowValue_StringValue{StringValue: string(content)},
			}},
		}},
		Latency:   durationpb.New(time.Since(startTime)),
		Statement: statement,
		RowsCount: 0, /* unknown */
	}}, nil
}

func getSimpleStatementResult(data []byte) (*v1pb.QueryResult, error) {
	rows, err := convertRows(data)
	if err != nil {
		return nil, err
	}

	result := &v1pb.QueryResult{
		ColumnNames:     []string{"result"},
		ColumnTypeNames: []string{"TEXT"},
	}

	for _, v := range rows {
		r, err := bson.MarshalExtJSONIndent(v, true, false, "", "  ")
		if err != nil {
			return nil, err
		}
		result.Rows = append(result.Rows, &v1pb.QueryRow{
			Values: []*v1pb.RowValue{
				{Kind: &v1pb.RowValue_StringValue{StringValue: string(r)}},
			},
		})
	}
	return result, nil
}

func convertRows(data []byte) ([]any, error) {
	var a any
	if err := bson.UnmarshalExtJSON(data, true, &a); err != nil {
		return nil, err
	}

	if aa, ok := a.(bson.A); ok {
		return []any(aa), nil
	}
	return []any{a}, nil
}

func isMongoStatement(statement string) bool {
	statement = strings.ToLower(statement)
	if strings.HasPrefix(statement, "db.") {
		return true
	}
	return strings.HasPrefix(statement, `db["`)
}
