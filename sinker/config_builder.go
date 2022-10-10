package sinker

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

// ReturnConfigYamlFromSink this is the main method, which will generate the YAML file from the
func ReturnConfigYamlFromSink(_ context.Context, kafkaUrlConfig, sinkId, sinkUrl, sinkUsername, sinkPassword string) (string, error) {

	config := OtelConfigFile{
		Receivers: Receivers{
			Kafka: KafkaReceiver{
				Brokers:         []string{kafkaUrlConfig},
				Topic:           fmt.Sprintf("otlp_metrics-%s", sinkId),
				ProtocolVersion: "2.0.0", // Leaving default of over 2.0.0
			},
		},
		Extensions: &Extensions{
			PProf: &PProfExtension{
				Endpoint: ":1888", // Leaving default for now, will need to change with more processes
			},
			BasicAuth: &BasicAuthenticationExtension{
				ClientAuth: &struct {
					Username string `json:"username" yaml:"username"`
					Password string `json:"password" yaml:"password"`
				}{Username: sinkUsername, Password: sinkPassword},
			},
		},
		Exporters: Exporters{
			PrometheusRemoteWrite: &PrometheusRemoteWriteExporterConfig{
				Endpoint: sinkUrl,
				auth: struct {
					Authenticator string `json:"authenticator" yaml:"authenticator"`
				}{Authenticator: "basicauth/exporter"},
			},
		},
		Service: ServiceConfig{
			Extensions: []string{"pprof", "basicauth/exporter"},
			Pipelines: struct {
				Metrics struct {
					Receivers  []string `json:"receivers" yaml:"receivers"`
					Processors []string `json:"processors,omitempty" yaml:"processors,omitempty"`
					Exporters  []string `json:"exporters" yaml:"exporters"`
				} `json:"metrics" yaml:"metrics"`
			}{
				Metrics: struct {
					Receivers  []string `json:"receivers" yaml:"receivers"`
					Processors []string `json:"processors,omitempty" yaml:"processors,omitempty"`
					Exporters  []string `json:"exporters" yaml:"exporters"`
				}{
					Receivers: []string{"kafka"},
					Exporters: []string{"prometheusremotewrite"},
				},
			},
		},
	}
	marshal, err := yaml.Marshal(&config)
	if err != nil {
		return "", err
	}
	returnedString := "---\n" + string(marshal)
	returnedString = strings.TrimRight(returnedString, `'`)
	return returnedString, nil

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
	HealthCheckExtConfig *HealthCheckExtension `json:"health_check,omitempty" yaml:"health_check,omitempty"`
	PProf                *PProfExtension       `json:"pprof,omitempty" yaml:"pprof,omitempty"`
	ZPages               *ZPagesExtension      `json:"zpages,omitempty" yaml:"zpages,omitempty"`
	// Exporters Authentication
	BasicAuth *BasicAuthenticationExtension `json:"basicauth/exporter,omitempty" yaml:"basicauth/exporter,omitempty"`
}

type HealthCheckExtension struct {
	CollectorPipeline *struct {
		Enabled          bool   `json:"enabled" yaml:"enabled"`
		Interval         string `json:"interval" yaml:"interval"`
		FailureThreshold int32  `json:"exporter_failure_threshold" yaml:"exporter_failure_threshold"`
	} `json:"check_collector_pipeline,omitempty" yaml:"check_collector_pipeline,omitempty"`
}

type PProfExtension struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type ZPagesExtension struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type BasicAuthenticationExtension struct {
	ClientAuth *struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	} `json:"client_auth" yaml:"client_auth"`
}

type Exporters struct {
	PrometheusRemoteWrite *PrometheusRemoteWriteExporterConfig `json:"prometheusremotewrite,omitempty" yaml:"prometheusremotewrite,omitempty"`
}

type PrometheusRemoteWriteExporterConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	auth     struct {
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
}
