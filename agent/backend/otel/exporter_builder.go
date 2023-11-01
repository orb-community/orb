package otel

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type ExporterBuilder interface {
	GetStructFromYaml(yamlString string) (openTelemetryConfig, error)
	MergeDefaultValueWithPolicy(config openTelemetryConfig, policyName string) (openTelemetryConfig, error)
}

type openTelemetryConfig struct {
	Receivers  map[string]interface{} `yaml:"receivers"`
	Processors map[string]interface{} `yaml:"processors"`
	Extensions map[string]interface{} `yaml:"extensions"`
	Exporters  *exporters             `yaml:"exporters"`
	Service    *service               `yaml:"service"`
}

type exporters struct {
	Otlp *defaultOtlpExporter `yaml:"otlp"`
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
}

type pipelines struct {
	Metrics *pipeline `yaml:"metrics"`
	Traces  *pipeline `yaml:"traces"`
	Logs    *pipeline `yaml:"logs"`
}

type pipeline struct {
	Exporters  []string `yaml:"exporters"`
	Receivers  []string `yaml:"receivers"`
	Processors []string `yaml:"processors"`
}

func getExporterBuilder(logger *zap.Logger) *exporterBuilder {
	return &exporterBuilder{logger: logger}
}

type exporterBuilder struct {
	logger *zap.Logger
}

func (e *exporterBuilder) GetStructFromYaml(yamlString string) (openTelemetryConfig, error) {
	var config openTelemetryConfig
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		e.logger.Error("failed to unmarshal yaml string", zap.Error(err))
		return openTelemetryConfig{}, err
	}
	return config, nil
}

func (e *exporterBuilder) MergeDefaultValueWithPolicy(config openTelemetryConfig, policyId string, policyName string) (openTelemetryConfig, error) {
	defaultOtlpExporter := defaultOtlpExporter{
		Endpoint: "localhost:4317",
		Tls: &tls{
			Insecure: true,
		},
	}
	// Override any openTelemetry exporter that may come, to connect to agent's otlp receiver
	config.Exporters = &exporters{&defaultOtlpExporter}
	if config.Processors == nil {
		config.Processors = make(map[string]interface{})
	}
	config.Processors["attributes/policy_data"] = map[string]interface{}{
		"actions": []struct {
			Key    string `yaml:"key"`
			Value  string `yaml:"value"`
			Action string `yaml:"action"`
		}{
			{Key: "policy_id", Value: policyId, Action: "insert"},
			{Key: "policy_name", Value: policyName, Action: "insert"},
		},
	}
	// Override metrics exporter and append attributes/policy_data processor
	if config.Service.Pipelines.Metrics != nil {
		config.Service.Pipelines.Metrics.Exporters = []string{"otlp"}
		config.Service.Pipelines.Metrics.Processors = append(config.Service.Pipelines.Metrics.Processors, "attributes/policy_data")
	}
	if config.Service.Pipelines.Traces != nil {
		config.Service.Pipelines.Traces.Exporters = []string{"otlp"}
		config.Service.Pipelines.Traces.Processors = append(config.Service.Pipelines.Traces.Processors, "attributes/policy_data")
	}
	if config.Service.Pipelines.Logs != nil {
		config.Service.Pipelines.Logs.Exporters = []string{"otlp"}
		config.Service.Pipelines.Logs.Processors = append(config.Service.Pipelines.Logs.Processors, "attributes/policy_data")
	}
	return config, nil
}

func (o *openTelemetryBackend) buildDefaultExporterAndProcessor(policyYaml string, policyId string, policyName string) (openTelemetryConfig, error) {
	defaultPolicyYaml, err := yaml.Marshal(policyYaml)
	if err != nil {
		o.logger.Warn("yaml policy marshal failure", zap.String("policy_id", policyId))
		return openTelemetryConfig{}, err
	}
	defaultPolicyString := string(defaultPolicyYaml)
	builder := getExporterBuilder(o.logger)
	defaultPolicyStruct, err := builder.GetStructFromYaml(defaultPolicyString)
	if err != nil {
		return openTelemetryConfig{}, err
	}
	defaultPolicyStruct, err = builder.MergeDefaultValueWithPolicy(defaultPolicyStruct, policyId, policyName)
	if err != nil {
		return openTelemetryConfig{}, err
	}
	return defaultPolicyStruct, nil
}
