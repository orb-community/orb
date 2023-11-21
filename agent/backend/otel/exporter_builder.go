package otel

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"strconv"
)

type ExporterBuilder interface {
	GetStructFromYaml(yamlString string) (openTelemetryConfig, error)
	MergeDefaultValueWithPolicy(config openTelemetryConfig, policyName string) (openTelemetryConfig, error)
}

type openTelemetryConfig struct {
	Receivers  map[string]interface{} `yaml:"receivers"`
	Processors map[string]interface{} `yaml:"processors,omitempty"`
	Extensions map[string]interface{} `yaml:"extensions,omitempty"`
	Exporters  map[string]interface{} `yaml:"exporters"`
	Service    *service               `yaml:"service"`
}

type defaultOtlpExporter struct {
	Endpoint string `yaml:"endpoint"`
	Tls      *tls   `yaml:"tls"`
}

type tls struct {
	Insecure bool `yaml:"insecure"`
}

type service struct {
	Pipelines *pipelines `yaml:"pipelines"`
	Telemetry *telemetry `yaml:"telemetry,omitempty"`
}

type telemetry struct {
	Metrics *metrics `yaml:"metrics,omitempty"`
	Logs    *logs    `yaml:"logs,omitempty"`
	Traces  *traces  `yaml:"traces,omitempty"`
}

type metrics struct {
	Level   string `yaml:"level,omitempty"`
	Address string `yaml:"address,omitempty"`
}

type traces struct {
	Enabled bool `yaml:"enabled"`
}

type logs struct {
	Enabled bool `yaml:"enabled"`
}

type pipelines struct {
	Metrics *pipeline `yaml:"metrics,omitempty"`
	Traces  *pipeline `yaml:"traces,omitempty"`
	Logs    *pipeline `yaml:"logs,omitempty"`
}

type pipeline struct {
	Exporters  []string `yaml:"exporters,omitempty"`
	Receivers  []string `yaml:"receivers,omitempty"`
	Processors []string `yaml:"processors,omitempty"`
}

func getExporterBuilder(logger *zap.Logger, host string, port int) *exporterBuilder {
	return &exporterBuilder{logger: logger, host: host, port: port}
}

type exporterBuilder struct {
	logger *zap.Logger
	host   string
	port   int
}

func (e *exporterBuilder) GetStructFromYaml(yamlString string) (openTelemetryConfig, error) {
	var config openTelemetryConfig
	if err := yaml.Unmarshal([]byte(yamlString), &config); err != nil {
		e.logger.Error("failed to unmarshal yaml string", zap.Error(err))
		return config, err
	}
	return config, nil
}

func (e *exporterBuilder) MergeDefaultValueWithPolicy(config openTelemetryConfig, policyId string, policyName string) (openTelemetryConfig, error) {
	endpoint := e.host + ":" + strconv.Itoa(e.port)
	defaultOtlpExporter := defaultOtlpExporter{
		Endpoint: endpoint,
		Tls: &tls{
			Insecure: true,
		},
	}

	// Override any openTelemetry exporter that may come, to connect to agent's otlp receiver
	config.Exporters = map[string]interface{}{
		"otlp": &defaultOtlpExporter,
	}
	if config.Processors == nil {
		config.Processors = make(map[string]interface{})
	}
	config.Processors["transform/policy_data"] = map[string]interface{}{
		"metric_statements": map[string]interface{}{
			"context": "scope",
			"statements": []string{
				`set(attributes["policy_id"], "` + policyId + `")`,
				`set(attributes["policy_name"], "` + policyName + `")`,
			},
		},
	}
	if config.Extensions == nil {
		config.Extensions = make(map[string]interface{})
	}
	tel := &telemetry{
		Metrics: &metrics{Level: "none"},
	}
	config.Service.Telemetry = tel
	// Override metrics exporter and append attributes/policy_data processor
	if config.Service.Pipelines.Metrics != nil {
		config.Service.Pipelines.Metrics.Exporters = []string{"otlp"}
		config.Service.Pipelines.Metrics.Processors = append(config.Service.Pipelines.Metrics.Processors, "transform/policy_data")
	}
	if config.Service.Pipelines.Traces != nil {
		config.Service.Pipelines.Traces.Exporters = []string{"otlp"}
		config.Service.Pipelines.Traces.Processors = append(config.Service.Pipelines.Traces.Processors, "transform/policy_data")
	}
	if config.Service.Pipelines.Logs != nil {
		config.Service.Pipelines.Logs.Exporters = []string{"otlp"}
		config.Service.Pipelines.Logs.Processors = append(config.Service.Pipelines.Logs.Processors, "transform/policy_data")
	}
	return config, nil
}

func (o *openTelemetryBackend) buildDefaultExporterAndProcessor(policyYaml string, policyId string, policyName string, telemetryPort int) (openTelemetryConfig, error) {
	defaultPolicyYaml, err := yaml.Marshal(policyYaml)
	if err != nil {
		o.logger.Warn("yaml policy marshal failure", zap.String("policy_id", policyId))
		return openTelemetryConfig{}, err
	}
	defaultPolicyString := string(defaultPolicyYaml)
	builder := getExporterBuilder(o.logger, o.otelReceiverHost, o.otelReceiverPort)
	defaultPolicyStruct, err := builder.GetStructFromYaml(defaultPolicyString)
	if err != nil {
		return openTelemetryConfig{}, err
	}
	defaultPolicyStruct, err = builder.MergeDefaultValueWithPolicy(
		defaultPolicyStruct,
		policyId,
		policyName)
	if err != nil {
		return openTelemetryConfig{}, err
	}
	return defaultPolicyStruct, nil
}
