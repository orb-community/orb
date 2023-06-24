package config

import (
	"database/sql/driver"
	"time"

	"github.com/orb-community/orb/pkg/types"
)

type SinkData struct {
	SinkID          string          `json:"sink_id"`
	OwnerID         string          `json:"owner_id"`
	Backend         string          `json:"backend"`
	Config          types.Metadata  `json:"config"`
	State           PrometheusState `json:"state,omitempty"`
	Msg             string          `json:"msg,omitempty"`
	LastRemoteWrite time.Time       `json:"last_remote_write,omitempty"`
}

type SinkConfigData struct {
	Authentication struct {
		Type     string `json:"type"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"authentication"`
	Exporter struct {
		Endpoint string `json:"endpoint"`
	} `json:"exporter"`
	Headers struct {
		OrgID string `json:"X-Scope-OrgID"`
	} `json:"headers"`
	Opentelemetry string `json:"opentelemetry"`
}


const (
	Unknown PrometheusState = iota
	Active
	Error
	Idle
	Warning
)

type PrometheusState int

var promStateMap = [...]string{
	"unknown",
	"active",
	"error",
	"idle",
	"warning",
}

var promStateRevMap = map[string]PrometheusState{
	"unknown": Unknown,
	"active":  Active,
	"error":   Error,
	"idle":    Idle,
	"warning": Warning,
}

func (p PrometheusState) String() string {
	return promStateMap[p]
}

func (p *PrometheusState) SetFromString(value string) error {
	*p = promStateRevMap[value]
	return nil
}

func (p PrometheusState) Value() (driver.Value, error) {
	return p.String(), nil
}

type OtelConfigFile struct {
	Receivers  Receivers     `json:"receivers" yaml:"receivers"`
	Processors *Processors   `json:"processors,omitempty" yaml:"processors,omitempty"`
	Extensions *Extensions   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Exporters  Exporters     `json:"exporters" yaml:"exporters"`
	Service    ServiceConfig `json:"service" yaml:"service"`
}

// Receivers will receive only with Kafka for now
type Receivers struct {
	Kafka KafkaReceiver `json:"kafka" yaml:"kafka"`
}

type KafkaReceiver struct {
	Brokers         []string `json:"brokers" yaml:"brokers"`
	Topic           string   `json:"topic" yaml:"topic"`
	ProtocolVersion string   `json:"protocol_version" yaml:"protocol_version"`
}

type Processors struct {
}

type Extensions struct {
	HealthCheckExtConfig *HealthCheckExtension `json:"health_check,omitempty" yaml:"health_check,omitempty" :"health_check_ext_config"`
	PProf                *PProfExtension       `json:"pprof,omitempty" yaml:"pprof,omitempty" :"p_prof"`
	ZPages               *ZPagesExtension      `json:"zpages,omitempty" yaml:"zpages,omitempty" :"z_pages"`
	// Exporters Authentication
	BasicAuth *BasicAuthenticationExtension `json:"basicauth/exporter,omitempty" yaml:"basicauth/exporter,omitempty" :"basic_auth"`
	//BearerAuth *BearerAuthExtension          `json:"bearerauth/exporter,omitempty" yaml:"bearerauth/exporter,omitempty" :"bearer_auth"`
}

type HealthCheckExtension struct {
	Endpoint          string                      `json:"endpoint" yaml:"endpoint"`
	Path              string                      `json:"path" yaml:"path"`
	CollectorPipeline *CollectorPipelineExtension `json:"check_collector_pipeline,omitempty" yaml:"check_collector_pipeline,omitempty"`
}

type CollectorPipelineExtension struct {
	Enabled          string `json:"enabled" yaml:"enabled"`
	Interval         string `json:"interval" yaml:"interval"`
	FailureThreshold int32  `json:"exporter_failure_threshold" yaml:"exporter_failure_threshold"`
}

type PProfExtension struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type ZPagesExtension struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type ClientAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type BasicAuthenticationExtension struct {
	ClientAuth *ClientAuth `json:"client_auth" yaml:"client_auth"`
}

type BearerAuthExtension struct {
	BearerAuth *struct {
		Token string `json:"token" yaml:"token"`
	} `json:"client_auth" yaml:"client_auth"`
}

type Exporters struct {
	PrometheusRemoteWrite *PrometheusRemoteWriteExporterConfig `json:"prometheusremotewrite,omitempty" yaml:"prometheusremotewrite,omitempty"`
	OTLPExporter          *OTLPExporterConfig                  `json:"otlphttp,omitempty" yaml:"otlphttp,omitempty"`
	LoggingExporter       *LoggingExporterConfig               `json:"logging,omitempty" yaml:"logging,omitempty"`
}

type LoggingExporterConfig struct {
	Verbosity          string `json:"verbosity,omitempty" yaml:"verbosity,omitempty"`
	SamplingInitial    int    `json:"sampling_initial,omitempty" yaml:"sampling_initial,omitempty"`
	SamplingThereAfter int    `json:"sampling_thereafter,omitempty" yaml:"sampling_thereafter,omitempty"`
}

type OTLPExporterConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Auth     struct {
		Authenticator string `json:"authenticator" yaml:"authenticator"`
	}
}

type Auth struct {
	Authenticator string `json:"authenticator" yaml:"authenticator"`
}

type PrometheusRemoteWriteExporterConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Auth     struct {
		Authenticator string `json:"authenticator" yaml:"authenticator"`
	}
}

type ServiceConfig struct {
	Extensions []string `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Pipelines  struct {
		Metrics struct {
			Receivers  []string `json:"receivers" yaml:"receivers"`
			Processors []string `json:"processors,omitempty" yaml:"processors,omitempty"`
			Exporters  []string `json:"exporters" yaml:"exporters"`
		} `json:"metrics" yaml:"metrics"`
	} `json:"pipelines" yaml:"pipelines"`
	Telemetry struct {
		Logs struct {
			Level            string   `json:"level,omitempty" yaml:"level,omitempty"`
			Encoding         string   `json:"encoding,omitempty" yaml:"encoding,omitempty"`
			OutputPaths      []string `json:"output_paths,omitempty" yaml:"output_paths,omitempty"`
			ErrorOutputPaths []string `json:"error_output_paths,omitempty" yaml:"error_output_paths,omitempty"`
		} `json:"logs" yaml:"logs,omitempty"`
	} `json:"telemetry,omitempty" yaml:"telemetry,omitempty"`
}
