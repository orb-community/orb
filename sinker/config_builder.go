package sinker

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
)

// ReturnConfigYamlFromSink this is the main method, which will generate the YAML file from the
func ReturnConfigYamlFromSink(_ context.Context, kafkaUrlConfig, sinkId, sinkUrl, sinkUsername, sinkPassword string) (string, error) {
	config := OtelConfigFile{
		Receivers: Receivers{
			Kafka: KafkaReceiver{
				Brokers:         []string{kafkaUrlConfig},
				Topic:           fmt.Sprintf("otlp_metrics/%s", sinkId),
				ProtocolVersion: "2.0.0", // Leaving default of over 2.0.0
			},
		},
		Processors: Processors{},
		Extensions: Extensions{
			HealthCheckExtConfig: HealthCheckExtension{},
			PProf: PProfExtension{
				Endpoint: ":1888", // Leaving default for now, will need to change with more processes
			},
			ZPages: ZPagesExtension{
				Endpoint: ":55679",
			},
			basicAuth: BasicAuthenticationExtension{
				ClientAuth: struct {
					Username string
					Password string
				}{Username: sinkUsername, Password: sinkPassword},
			},
		},
		Exporters: Exporters{
			PrometheusRemoteWrite: PrometheusRemoteWriteExporterConfig{
				Endpoint: sinkUrl,
				auth: struct {
					authenticator string `yaml:"authenticator"`
				}{authenticator: "basic_auth/exporter"},
			},
		},
		Service: ServiceConfig{},
	}
	marshal, err := yaml.Marshal(&config)
	if err != nil {
		return "", err
	}
	return string(marshal), nil

}

type OtelConfigFile struct {
	Receivers  Receivers     `yaml:"receivers"`
	Processors Processors    `yaml:"processors"`
	Extensions Extensions    `yaml:"extensions"`
	Exporters  Exporters     `yaml:"exporters"`
	Service    ServiceConfig `yaml:"service"`
}

// Receivers will receive only with Kafka for now
type Receivers struct {
	Kafka KafkaReceiver `yaml:"kafka"`
}

type KafkaReceiver struct {
	Brokers         []string `yaml:"brokers"`
	Topic           string   `yaml:"topic"`
	ProtocolVersion string   `yaml:"protocol_version"`
}

type Processors struct {
}

type Extensions struct {
	HealthCheckExtConfig HealthCheckExtension `yaml:"health_check"`
	PProf                PProfExtension       `yaml:"pprof"`
	ZPages               ZPagesExtension      `yaml:"zpages"`
	// Exporters Authentication
	basicAuth BasicAuthenticationExtension `yaml:"basic_auth/exporter"`
}

type HealthCheckExtension struct {
	CollectorPipeline struct {
		Enabled          bool   `yaml:"enabled"`
		interval         string `yaml:"interval"`
		FailureThreshold int32  `yaml:"exporter_failure_threshold"`
	} `yaml:"check_collector_pipeline"`
}

type PProfExtension struct {
	Endpoint string `yaml:"endpoint"`
}

type ZPagesExtension struct {
	Endpoint string `yaml:"endpoint"`
}

type BasicAuthenticationExtension struct {
	ClientAuth struct {
		Username string
		Password string
	} `yaml:"client_auth"`
}

type Exporters struct {
	PrometheusRemoteWrite PrometheusRemoteWriteExporterConfig `yaml:"prometheusremotewrite"`
}

type PrometheusRemoteWriteExporterConfig struct {
	Endpoint string `yaml:"endpoint"`
	auth     struct {
		authenticator string `yaml:"authenticator"`
	}
}

type ServiceConfig struct {
	Extensions string `yaml:"extensions"`
	Pipelines  struct {
		Metrics struct {
			Receivers string `yaml:"receivers"`
			Exporters string `yaml:"exporters"`
		} `yaml:"metrics"`
	} `yaml:"pipelines"`
}
