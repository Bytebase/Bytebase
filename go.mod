module github.com/bytebase/bytebase

go 1.23.4

// workaround mssql-docker default TLS cert negative serial number problem
// https://github.com/microsoft/mssql-docker/issues/895
godebug x509negativeserial=1

require (
	cloud.google.com/go/bigquery v1.65.0
	cloud.google.com/go/cloudsqlconn v1.13.2
	cloud.google.com/go/secretmanager v1.14.2
	cloud.google.com/go/spanner v1.73.0
	gitee.com/chunanyong/dm v1.8.17
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.8.0
	github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos v1.2.0
	github.com/ClickHouse/clickhouse-go/v2 v2.30.0
	github.com/alexmullins/zip v0.0.0-20180717182244-4affb64b04d0
	github.com/antlr4-go/antlr/v4 v4.13.1
	github.com/aws/aws-sdk-go-v2 v1.32.7
	github.com/aws/aws-sdk-go-v2/config v1.28.7
	github.com/aws/aws-sdk-go-v2/feature/rds/auth v1.5.2
	github.com/aws/aws-sdk-go-v2/service/licensemanager v1.29.8
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.34.8
	github.com/beltran/gohive v1.7.0
	github.com/blang/semver/v4 v4.0.0
	github.com/bytebase/google-sql-parser v0.0.0-20250116032737-689a327f9465
	github.com/bytebase/mysql-parser v0.0.0-20241224071214-cb9fd84811dd
	github.com/bytebase/plsql-parser v0.0.0-20241217100236-775726d725f9
	github.com/bytebase/postgresql-parser v0.0.0-20241108072816-e932f744004a
	github.com/bytebase/snowsql-parser v0.0.0-20241211024554-093949951501
	github.com/bytebase/tidb-parser v0.0.0-20240821091609-162fffe2a839
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/cockroachdb/cockroachdb-parser v0.23.2
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/databricks/databricks-sdk-go v0.54.0
	github.com/elastic/go-elasticsearch/v7 v7.17.10
	github.com/github/gh-ost v1.1.6
	github.com/go-ego/gse v0.80.3
	github.com/go-ldap/ldap/v3 v3.4.10
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/golang-sql/sqlexp v0.1.0
	github.com/google/cel-go v0.22.1
	github.com/google/go-cmp v0.6.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/gosimple/slug v1.15.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.24.0
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/hashicorp/vault/api v1.15.0
	github.com/hashicorp/vault/api/auth/approle v0.8.0
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/jackc/pgtype v1.14.4
	github.com/jackc/pgx/v5 v5.7.1
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
	github.com/labstack/echo-contrib v0.17.1
	github.com/labstack/echo/v4 v4.13.3
	github.com/lestrrat-go/jwx/v2 v2.1.3
	github.com/lib/pq v1.10.9
	github.com/lor00x/goldap v0.0.0-20240304151906-8d785c64d1c8
	github.com/mattn/go-oci8 v0.1.1
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/microsoft/go-mssqldb v1.8.0
	github.com/nyaruka/phonenumbers v1.4.4
	github.com/paulmach/orb v0.11.1
	github.com/pganalyze/pg_query_go/v5 v5.1.0
	github.com/pingcap/tidb v1.1.0-beta.0.20241125141335-ec8b81b98edc
	github.com/pingcap/tidb/pkg/parser v0.0.0-20241125141335-ec8b81b98edc
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.7.0
	github.com/sashabaranov/go-openai v1.36.0
	github.com/segmentio/analytics-go v3.1.0+incompatible
	github.com/shopspring/decimal v1.4.0
	github.com/sijms/go-ora/v2 v2.8.22
	github.com/snowflakedb/gosnowflake v1.12.1
	github.com/soheilhy/cmux v0.1.5
	github.com/sourcegraph/conc v0.3.0
	github.com/sourcegraph/go-lsp v0.0.0-20240223163137-f80c5dd31dfd
	github.com/sourcegraph/jsonrpc2 v0.2.0
	github.com/spf13/cobra v1.8.1
	github.com/stretchr/testify v1.10.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75
	github.com/vjeantet/ldapserver v1.0.1
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8
	github.com/xuri/excelize/v2 v2.9.0
	github.com/zeebo/xxh3 v1.0.2
	go.mongodb.org/mongo-driver v1.17.1
	go.uber.org/multierr v1.11.0
	golang.org/x/crypto v0.31.0
	golang.org/x/exp v0.0.0-20241217172543-b2144cdd0a67
	golang.org/x/oauth2 v0.24.0
	golang.org/x/text v0.21.0
	google.golang.org/api v0.214.0
	google.golang.org/genproto v0.0.0-20241230172942-26aa7a208def
	google.golang.org/genproto/googleapis/api v0.0.0-20241230172942-26aa7a208def
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241230172942-26aa7a208def
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cel.dev/expr v0.18.0 // indirect
	cloud.google.com/go/auth v0.13.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.6 // indirect
	cloud.google.com/go/monitoring v1.22.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.16.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.3.2 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.3.2 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.24.2 // indirect
	github.com/JohnCGriffin/overflow v0.0.0-20211019200055-46fa312c352c // indirect
	github.com/apache/arrow/go/v15 v15.0.2 // indirect
	github.com/apache/arrow/go/v16 v16.0.0 // indirect
	github.com/apache/thrift v0.19.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.48 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.57.1 // indirect
	github.com/bazelbuild/rules_go v0.49.0 // indirect
	github.com/beltran/gosasl v0.0.0-20240210185013-36d7ba6de436 // indirect
	github.com/beltran/gssapi v0.0.0-20200324152954-d86554db4bab // indirect
	github.com/biogo/store v0.0.0-20201120204734-aad293a2328f // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cloudfoundry/gosigar v1.3.6 // indirect
	github.com/cockroachdb/apd/v3 v3.1.0 // indirect
	github.com/cockroachdb/errors v1.9.0 // indirect
	github.com/cockroachdb/logtags v0.0.0-20211118104740-dabe8e521a4f // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/danieljoos/wincred v1.2.1 // indirect
	github.com/dave/dst v0.27.2 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.7.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/getsentry/sentry-go v0.12.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.7 // indirect
	github.com/go-jose/go-jose/v4 v4.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-zookeeper/zk v1.0.1 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang/geo v0.0.0-20210211234256-740aa86cb551 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.8 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.6 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/goidentity/v6 v6.0.1 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/onsi/gomega v1.27.10 // indirect
	github.com/petermattis/goid v0.0.0-20211229010228-4d14c490ee36 // indirect
	github.com/pierrre/geohash v1.0.0 // indirect
	github.com/pingcap/sysutil v1.0.1-0.20230407040306-fb007c5aff21 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/tiancaiamao/gp v0.0.0-20221230034425-4025bc8a4d4a // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/twmb/murmur3 v1.1.6 // indirect
	github.com/twpayne/go-geom v1.4.1 // indirect
	github.com/twpayne/go-kml v1.5.2 // indirect
	github.com/vcaesar/cedar v0.20.2 // indirect
	github.com/xuri/efp v0.0.0-20240408161823-9ad904a10d6d // indirect
	github.com/xuri/nfp v0.0.0-20240318013403-ab9948c2c4a7 // indirect
	go.etcd.io/etcd/api/v3 v3.5.10 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.10 // indirect
	go.etcd.io/etcd/client/v3 v3.5.10 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.31.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.54.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/sdk v1.31.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.31.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.27.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

require (
	cloud.google.com/go v0.117.0 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	cloud.google.com/go/iam v1.3.0 // indirect
	cloud.google.com/go/longrunning v0.6.3 // indirect
	github.com/ClickHouse/ch-go v0.61.5 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.3 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.26 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.26 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.38.1
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.3
	github.com/aws/smithy-go v1.22.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bytebase/partiql-parser v0.0.0-20240531101102-1962ff456f2c
	github.com/bytebase/tsql-parser v0.0.0-20241114064553-239d32311c0b
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20240905190251-b4127c9b8d78 // indirect
	github.com/cockroachdb/cockroach-go/v2 v2.3.8
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/danjacques/gofslock v0.0.0-20220131014315-6e321f4509c8 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/envoyproxy/go-control-plane v0.13.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.4 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-mysql-org/go-mysql v1.6.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v24.3.25+incompatible // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20230326075908-cb1d2100619a // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/openark/golib v0.0.0-20210531070646-355f37940af8
	github.com/opentracing/basictracer-go v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pingcap/errors v0.11.5-0.20221009092201-b66cddb77c32 // indirect
	github.com/pingcap/failpoint v0.0.0-20240528011301-b51a646c7c86 // indirect
	github.com/pingcap/kvproto v0.0.0-20231122054644-fb0f5c2a0a10 // indirect
	github.com/pingcap/log v1.1.1-0.20230317032135-a0d097d16e22 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/pquerna/cachecontrol v0.2.0 // indirect
	github.com/pquerna/otp v1.4.0
	github.com/prometheus/client_golang v1.19.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.13.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/backo-go v1.0.1 // indirect
	github.com/shirou/gopsutil/v3 v3.23.12 // indirect
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726 // indirect
	github.com/siddontang/go-log v0.0.0-20190221022429-1e957dd83bed // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tikv/client-go/v2 v2.0.8-0.20231116051730-1c2351c28173 // indirect
	github.com/tikv/pd/client v0.0.0-20231127075044-9f4803d8bd05 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.31.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	nhooyr.io/websocket v1.8.10 // indirect
)

// copied from pingcap/tidb
// fix potential security issue(CVE-2020-26160) introduced by indirect dependency.
replace github.com/dgrijalva/jwt-go => github.com/form3tech-oss/jwt-go v3.2.6-0.20210809144907-32ab6a8243d7+incompatible

replace github.com/github/gh-ost => github.com/bytebase/gh-ost2 v0.0.0-20240731091141-e7a12a9401ab

replace github.com/pingcap/tidb => github.com/bytebase/tidb2 v0.0.0-20231129002249-5bbb6bb83940

replace github.com/pingcap/tidb/pkg/parser => github.com/bytebase/tidb2/pkg/parser v0.0.0-20231129002249-5bbb6bb83940

replace github.com/pganalyze/pg_query_go/v5 => github.com/bytebase/pg_query_go2/v5 v5.0.0-20241024030058-bbe7104c5b4b

replace github.com/mattn/go-oci8 => github.com/bytebase/go-obo v0.0.0-20231026081615-705a7fffbfd2

replace github.com/antlr4-go/antlr/v4 => github.com/bytebase/antlr/v4 v4.0.0-20240827034948-8c385f108920

replace github.com/sourcegraph/go-lsp => github.com/bytebase/go-lsp v0.0.0-20240130071507-c04b5c75010c

replace (
	github.com/beltran/gohive => github.com/bytebase/gohive v0.0.0-20240422092929-d76993a958a4
	github.com/beltran/gosasl => github.com/bytebase/gosasl v0.0.0-20240422091407-6b7481e86f08
)

replace github.com/microsoft/go-mssqldb => github.com/bytebase/go-mssqldb v0.0.0-20240801091126-3ff3ca07d898

replace github.com/youmark/pkcs8 => github.com/bytebase/pkcs8 v0.0.0-20240612095628-fcd0a7484c94

replace github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos => github.com/bytebase/azure-sdk-for-go/sdk/data/azcosmos v0.0.0-20250109032656-87cf24d45689
