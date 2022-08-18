package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bytebase/bytebase/api"
	"github.com/bytebase/bytebase/common"
	"github.com/bytebase/bytebase/common/log"
	"github.com/bytebase/bytebase/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	// Register clickhouse driver.
	_ "github.com/bytebase/bytebase/plugin/db/clickhouse"
	// Register mysql driver.
	_ "github.com/bytebase/bytebase/plugin/db/mysql"
	// Register postgres driver.
	_ "github.com/bytebase/bytebase/plugin/db/pg"
	// Register snowflake driver.
	_ "github.com/bytebase/bytebase/plugin/db/snowflake"
	// Register sqlite driver.
	_ "github.com/bytebase/bytebase/plugin/db/sqlite"

	// Register pingcap parser driver.
	_ "github.com/pingcap/tidb/types/parser_driver"
	// Register fake advisor.
	_ "github.com/bytebase/bytebase/plugin/advisor/fake"
	// Register mysql advisor.
	_ "github.com/bytebase/bytebase/plugin/advisor/mysql"
	// Register postgresql advisor.
	_ "github.com/bytebase/bytebase/plugin/advisor/pg"

	// Register postgres parser driver.
	_ "github.com/bytebase/bytebase/plugin/parser/engine/pg"
)

// -----------------------------------Global constant BEGIN----------------------------------------.
const (

	// greetingBanner is the greeting banner.
	// http://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=Bytebase
	greetingBanner = `
██████╗ ██╗   ██╗████████╗███████╗██████╗  █████╗ ███████╗███████╗
██╔══██╗╚██╗ ██╔╝╚══██╔══╝██╔════╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██████╔╝ ╚████╔╝    ██║   █████╗  ██████╔╝███████║███████╗█████╗
██╔══██╗  ╚██╔╝     ██║   ██╔══╝  ██╔══██╗██╔══██║╚════██║██╔══╝
██████╔╝   ██║      ██║   ███████╗██████╔╝██║  ██║███████║███████╗
╚═════╝    ╚═╝      ╚═╝   ╚══════╝╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝

%s
___________________________________________________________________________________________

`
	// byeBanner is the bye banner.
	// http://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=BYE
	byeBanner = `
██████╗ ██╗   ██╗███████╗
██╔══██╗╚██╗ ██╔╝██╔════╝
██████╔╝ ╚████╔╝ █████╗
██╔══██╗  ╚██╔╝  ██╔══╝
██████╔╝   ██║   ███████╗
╚═════╝    ╚═╝   ╚══════╝

`
)

// -----------------------------------Global constant END------------------------------------------

// -----------------------------------Command Line Config BEGIN------------------------------------.
var (
	flags struct {
		// Used for Bytebase command line config
		host         string
		port         int
		frontendHost string
		frontendPort int
		dataDir      string
		// When we are running in readonly mode:
		// - The data file will be opened in readonly mode, no applicable migration or seeding will be applied.
		// - Requests other than GET will be rejected
		// - Any operations involving mutation will not start (e.g. Background schema syncer, task scheduler)
		readonly bool
		// demo is a flag to seed the database with demo data.
		demo bool
		// demoName is the name of the demo. It is only used when --demo is set,
		// and should be one of the subpath name in the ./store/demo/ directory.
		demoName string
		debug    bool
		// pgURL must follow PostgreSQL connection URIs pattern.
		// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
		pgURL string

		// Cloud backup configs
		backupRegion      string
		backupBucket      string
		backupCredentials string
	}
	rootCmd = &cobra.Command{
		Use:   "bytebase",
		Short: "Bytebase is a database schema change and version control tool",
		Run: func(_ *cobra.Command, _ []string) {
			if flags.frontendHost == "" {
				flags.frontendHost = flags.host
			}
			if flags.frontendPort == 0 {
				flags.frontendPort = flags.port
			}

			start()

			fmt.Print(byeBanner)
		},
	}
)

type backupMeta struct {
	storageBackend  api.BackupStorageBackend
	region          string
	bucket          string
	credentialsFile string
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flags.host, "host", "http://localhost", "host where Bytebase backend is accessed from, must start with http:// or https://. This is used by Bytebase to create the webhook callback endpoint for VCS integration")
	rootCmd.PersistentFlags().IntVar(&flags.port, "port", 80, "port where Bytebase backend is accessed from. This is also used by Bytebase to create the webhook callback endpoint for VCS integration")
	rootCmd.PersistentFlags().StringVar(&flags.frontendHost, "frontend-host", "", "host where Bytebase frontend is accessed from, must start with http:// or https://. This is used by Bytebase to compose the frontend link when posting the webhook event. Default is the same as --host")
	rootCmd.PersistentFlags().IntVar(&flags.frontendPort, "frontend-port", 0, "port where Bytebase frontend is accessed from. This is used by Bytebase to compose the frontend link when posting the webhook event. Default is the same as --port")
	rootCmd.PersistentFlags().StringVar(&flags.dataDir, "data", ".", "directory where Bytebase stores data. If relative path is supplied, then the path is relative to the directory where Bytebase is under")
	rootCmd.PersistentFlags().BoolVar(&flags.readonly, "readonly", false, "whether to run in read-only mode")
	rootCmd.PersistentFlags().BoolVar(&flags.demo, "demo", false, "whether to run using demo data")
	rootCmd.PersistentFlags().StringVar(&flags.demoName, "demo-name", "", "name of the demo to use when running in demo mode")
	rootCmd.PersistentFlags().BoolVar(&flags.debug, "debug", false, "whether to enable debug level logging")
	rootCmd.PersistentFlags().StringVar(&flags.pgURL, "pg", "", "optional external PostgreSQL instance connection url(must provide dbname); for example postgresql://user:secret@masterhost:5432/dbname?sslrootcert=cert")

	// Cloud backup related flags.
	// TODO(dragonly): Add GCS usages when it's supported.
	rootCmd.PersistentFlags().StringVar(&flags.backupRegion, "backup-region", "", "region of the backup bucket, e.g., us-west-2 for AWS S3.")
	rootCmd.PersistentFlags().StringVar(&flags.backupBucket, "backup-bucket", "", "bucket where Bytebase stores backup data, e.g., s3://example-bucket. When provided, Bytebase will store data to the S3 bucket.")
	rootCmd.PersistentFlags().StringVar(&flags.backupCredentials, "backup-credentials", "", "credentials file to use for the backup bucket. It should be the same format as the AWS credential files.")
}

// -----------------------------------Command Line Config END--------------------------------------

// -----------------------------------Main Entry Point---------------------------------------------

func checkDataDir() error {
	// Convert to absolute path if relative path is supplied.
	if !filepath.IsAbs(flags.dataDir) {
		absDir, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/" + flags.dataDir)
		if err != nil {
			return err
		}
		flags.dataDir = absDir
	}

	// Trim trailing / in case user supplies
	flags.dataDir = strings.TrimRight(flags.dataDir, "/")

	if _, err := os.Stat(flags.dataDir); err != nil {
		return errors.Wrapf(err, "unable to access --data directory %s", flags.dataDir)
	}

	return nil
}

func start() {
	if flags.debug {
		log.SetLevel(zap.DebugLevel)
	}
	defer log.Sync()

	// check flags
	if !common.HasPrefixes(flags.host, "http://", "https://") {
		log.Error(fmt.Sprintf("--host %s must start with http:// or https://", flags.host))
		return
	}
	if err := checkDataDir(); err != nil {
		log.Error(err.Error())
		return
	}

	var profile server.Profile
	// This enables backup to cloud, and all backup data will be stored in the supported cloud storage.
	if flags.backupBucket != "" {
		if flags.backupCredentials == "" {
			log.Error("Must specify --backup-credentials when --backup-bucket is present.")
			return
		}
		bucketMeta, err := parseBucketURI(flags.backupBucket)
		if err != nil {
			log.Error("failed to parse backup bucket", zap.Error(err))
			return
		}
		if bucketMeta.storageBackend == api.BackupStorageBackendS3 {
			if flags.backupRegion == "" {
				log.Error("Must specify --backup-region for AWS S3 backup.")
				return
			}
			bucketMeta.region = flags.backupRegion
		}
		bucketMeta.credentialsFile = flags.backupCredentials
		profile = activeProfile(flags.dataDir, bucketMeta)
	} else {
		profile = activeProfile(flags.dataDir, backupMeta{storageBackend: api.BackupStorageBackendLocal})
	}

	var s *server.Server
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	// Trigger graceful shutdown on SIGINT or SIGTERM.
	// The default signal sent by the `kill` command is SIGTERM,
	// which is taken as the graceful shutdown signal for many systems, eg., Kubernetes, Gunicorn.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		log.Info(fmt.Sprintf("%s received.", sig.String()))
		if s != nil {
			_ = s.Shutdown(ctx)
		}
		cancel()
	}()

	s, err := server.NewServer(ctx, profile)
	if err != nil {
		log.Error("Cannot new server", zap.Error(err))
		return
	}
	fmt.Printf(greetingBanner, fmt.Sprintf("Version %s has started at %s:%d", profile.Version, profile.BackendHost, profile.BackendPort))
	// Execute program.
	if err := s.Run(ctx); err != nil {
		if err != http.ErrServerClosed {
			log.Error(err.Error())
			_ = s.Shutdown(ctx)
			cancel()
		}
	}

	// Wait for CTRL-C.
	<-ctx.Done()
}

// Examples:
//   s3:us-west-2//dev-bytebase-backup
//   gcs://dev-bytebase-backup
func parseBucketURI(uri string) (backupMeta, error) {
	parts := strings.Split(uri, "://")
	if len(parts) != 2 {
		return backupMeta{}, errors.Errorf("invalid bucket URI %q, expected format is s3://${BUCKET_NAME}", uri)
	}

	backend, bucket := parts[0], parts[1]
	if strings.Contains(bucket, "/") {
		return backupMeta{}, errors.Errorf("invalid bucket URI %q, expecting no / in the BUCKET_NAME", uri)
	}

	switch strings.ToUpper(backend) {
	case string(api.BackupStorageBackendS3):
		return backupMeta{
			storageBackend: api.BackupStorageBackendS3,
			bucket:         bucket,
		}, nil
	default:
		return backupMeta{}, errors.Errorf("unsupported storage backend %q", backend)
	}
}
