package mongodb

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/plugin/db"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
)

func TestGetMongoDBConnectionURL(t *testing.T) {
	tests := []struct {
		connConfig db.ConnectionConfig
		want       string
	}{
		{
			connConfig: db.ConnectionConfig{
				Host:                 "localhost",
				Port:                 "27017",
				Username:             "",
				Password:             "",
				MaximumSQLResultSize: common.DefaultMaximumSQLResultSize,
			},
			want: "mongodb://localhost:27017/?appName=bytebase&authSource=admin",
		},
		{
			connConfig: db.ConnectionConfig{
				Host:                 "localhost",
				Port:                 "27017",
				Username:             "",
				Password:             "",
				DirectConnection:     true,
				MaximumSQLResultSize: common.DefaultMaximumSQLResultSize,
			},
			want: "mongodb://localhost:27017/?appName=bytebase&authSource=admin&directConnection=true",
		},
		{
			connConfig: db.ConnectionConfig{
				Host:                 "localhost",
				Port:                 "27017",
				Username:             "",
				Password:             "",
				Database:             "sampleDB",
				MaximumSQLResultSize: common.DefaultMaximumSQLResultSize,
			},
			want: "mongodb://localhost:27017/sampleDB?appName=bytebase&authSource=admin",
		},
		{
			connConfig: db.ConnectionConfig{
				Host:                 "cluster0.sample.mongodb.net",
				Username:             "bytebase",
				Password:             "passwd",
				Database:             "sampleDB",
				SRV:                  true,
				MaximumSQLResultSize: common.DefaultMaximumSQLResultSize,
			},
			want: "mongodb+srv://bytebase:passwd@cluster0.sample.mongodb.net/sampleDB?appName=bytebase&authSource=admin",
		},
		{
			connConfig: db.ConnectionConfig{
				Host:                   "cluster0.sample.mongodb.net",
				Username:               "bytebase",
				Password:               "passwd",
				Database:               "sampleDB",
				AuthenticationDatabase: "",
				SRV:                    true,
				MaximumSQLResultSize:   common.DefaultMaximumSQLResultSize,
			},
			want: "mongodb+srv://bytebase:passwd@cluster0.sample.mongodb.net/sampleDB?appName=bytebase&authSource=admin",
		},
		{
			connConfig: db.ConnectionConfig{
				Host:                   "cluster0.sample.mongodb.net",
				Username:               "bytebase",
				Password:               "passwd",
				Database:               "sampleDB",
				AuthenticationDatabase: "admin",
				SRV:                    true,
				MaximumSQLResultSize:   common.DefaultMaximumSQLResultSize,
			},
			want: "mongodb+srv://bytebase:passwd@cluster0.sample.mongodb.net/sampleDB?appName=bytebase&authSource=admin",
		},
		{
			connConfig: db.ConnectionConfig{
				Host:                   "node1.cluster0.sample.mongodb.net",
				Port:                   "27017",
				Username:               "bytebase",
				Password:               "passwd",
				Database:               "sampleDB",
				AuthenticationDatabase: "admin",
				SRV:                    false,
				AdditionalAddresses: []*storepb.DataSourceOptions_Address{
					{Host: "node2.cluster0.sample.mongodb.net", Port: "27017"},
					{Host: "node3.cluster0.sample.mongodb.net", Port: "27017"},
				},
				ReplicaSet:           "rs0",
				MaximumSQLResultSize: common.DefaultMaximumSQLResultSize,
			},
			want: "mongodb://bytebase:passwd@node1.cluster0.sample.mongodb.net:27017,node2.cluster0.sample.mongodb.net:27017,node3.cluster0.sample.mongodb.net:27017/sampleDB?appName=bytebase&authSource=admin&replicaSet=rs0",
		},
	}

	a := require.New(t)
	for _, tt := range tests {
		got := getBasicMongoDBConnectionURI(tt.connConfig)
		a.Equal(tt.want, got)
	}
}

func TestIsMongoStatement(t *testing.T) {
	tests := []struct {
		statement string
		want      bool
	}{
		{
			statement: `show collections`,
			want:      false,
		},
		{
			statement: `db.cpl_station_info.find().limit(100)`,
			want:      true,
		},
	}

	a := require.New(t)
	for _, tt := range tests {
		got := isMongoStatement(tt.statement)
		a.Equal(tt.want, got, tt.statement)
	}
}

func TestGetSimpleStatementResult(t *testing.T) {
	testData1 := `{
  "_id": {
    "$oid": "66f62cad7195ccc0dbdfafbb"
  },
  "a": {
    "$numberLong": "1546786128982089728"
  }
}`

	testData2 := `{
  "_id": {
    "$oid": "66f6758c30daae815ac8784f"
  },
  "name": "dannyyy",
  "groups": [
    "123",
    "222"
  ]
}`

	tests := []struct {
		data string
		want *v1pb.QueryResult
	}{
		{
			data: fmt.Sprintf(`[%s, %s]`, testData1, testData2),
			want: &v1pb.QueryResult{
				ColumnNames:     []string{"result"},
				ColumnTypeNames: []string{"TEXT"},
				Rows: []*v1pb.QueryRow{
					{
						Values: []*v1pb.RowValue{
							{Kind: &v1pb.RowValue_StringValue{StringValue: testData1}},
						},
					},
					{
						Values: []*v1pb.RowValue{
							{Kind: &v1pb.RowValue_StringValue{StringValue: testData2}},
						},
					},
				},
			},
		},
	}

	a := require.New(t)
	for _, tt := range tests {
		result, err := getSimpleStatementResult([]byte(tt.data))
		a.NoError(err)
		diff := cmp.Diff(tt.want, result, protocmp.Transform(), protocmp.IgnoreMessages(&durationpb.Duration{}))
		a.Empty(diff)
	}
}
