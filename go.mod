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
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-version v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.9.0
	github.com/mainflux/mainflux v0.12.0
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/mitchellh/mapstructure v1.4.3
	github.com/opentracing/opentracing-go v1.2.0
	github.com/ory/dockertest/v3 v3.6.0
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/prometheus v1.8.2-0.20210621150501-ff58416a0b02
	github.com/rubenv/sql-migrate v0.0.0-20200616145509-8d140a17f351
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.1
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	go.opentelemetry.io/collector/model v0.48.0
	go.uber.org/zap v1.21.0
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/GehirnInc/crypt v0.0.0-20200316065508-bb7000b8a962 // indirect
	github.com/knadh/koanf v1.4.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/basicauthextension v0.48.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/tg123/go-htpasswd v1.2.0 // indirect
	go.opentelemetry.io/otel v1.6.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20210915214749-c084706c2272 // indirect
)

//These libs are used to allow orb extend opentelemetry features
require (
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.40.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.40.0
	github.com/prometheus/common v0.32.1
	go.opentelemetry.io/collector v0.48.0
	go.opentelemetry.io/otel/metric v0.28.0
	go.opentelemetry.io/otel/trace v1.6.1
	k8s.io/client-go v0.22.4
)
