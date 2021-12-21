module github.com/ns1labs/orb

go 1.17

require (
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/fatih/structs v1.1.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-cmd/cmd v1.3.0
	github.com/go-co-op/gocron v1.9.0
	github.com/go-kit/kit v0.11.0
	github.com/go-redis/redis/v8 v8.11.0
	github.com/go-zoo/bone v1.3.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4
	github.com/hashicorp/go-version v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.9.0
	github.com/mainflux/mainflux v0.12.0
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/mitchellh/mapstructure v1.4.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/ory/dockertest/v3 v3.6.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/prometheus v1.8.2-0.20210621150501-ff58416a0b02
	github.com/rubenv/sql-migrate v0.0.0-20200616145509-8d140a17f351
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	go.uber.org/zap v1.19.1
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

//These libs are used to allow orb extend opentelemetry features
require (
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.40.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.40.0
	github.com/prometheus/common v0.32.1
	go.opentelemetry.io/collector v0.40.0
	go.opentelemetry.io/collector/model v0.40.0
	go.opentelemetry.io/otel/metric v0.25.0
	go.opentelemetry.io/otel/trace v1.2.0
	google.golang.org/genproto v0.0.0-20210917145530-b395a37504d4
	k8s.io/client-go v0.22.4
)
