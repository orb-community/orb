package sinkOrchestrator

import (
	"github.com/ns1labs/orb/sinker/config"
	"strings"
)

const configuration = "receivers:\n  kafka:\n    brokers:\n      - kafka1:19092\n    topic: otlp_metrics\n    protocol_version: 2.0.0\n\nprocessors:\n  batch:\n\nextensions:\n  health_check:\n    check_collector_pipeline:\n      enabled: true\n      interval: \"1m\"\n      exporter_failure_threshold: 5\n\n  pprof:\n    endpoint: :1888\n\n  zpages:\n    endpoint: :55679\n\n  basicauth/client:\n    client_auth:\n      username: \n      password: \n\nexporters:\n  prometheusremotewrite:\n    endpoint: ${PROM_URL}\n    auth:\n      authenticator: basicauth/client\n  logging:\n    loglevel: debug\n\nservice:\n  extensions: [pprof, zpages, health_check, basicauth/client]\n  pipelines:\n    metrics:\n      receivers: [kafka]\n      exporters: [prometheusremotewrite]\n"

type ColConfig struct {
	Receivers  map[string]interface{} `json:"receivers" mapstructure:"receivers"`
	Processors map[string]interface{} `json:"processors" mapstructure:"processors"`
	Extensions map[string]interface{} `json:"extensions" mapstructure:"extensions"`
	Pprof      map[string]interface{} `json:"pprof" mapstructure:"pprof"`
	Zpages     map[string]interface{} `json:"zpages" mapstructure:"zpages"`
	BasicAuth  basicAuth              `json:"basicauth" mapstructure:"basicauth"`
	Exporters  exporters              `json:"exporters" mapstructure:"exporters"`
	Service    map[string]interface{} `json:"service" mapstructure:"service"`
}

type basicAuth struct {
	clientAuth
}

type clientAuth struct {
	username string
	password string
}

type exporters struct {
	prometheusremotewrite
	Logging map[string]interface{} `json:"logging" mapstructure:"logging"`
}

type prometheusremotewrite struct {
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`
}

func UpdateSink(sinkConfig config.SinkConfig) string {

	configFile := BuildConfigFile(sinkConfig)
	return configFile
}

func BuildConfigFile(sinkConfig config.SinkConfig) string {
	var config strings.Builder

	config.WriteString(
		"receivers:\n  kafka:\n    brokers:\n      - kafka1:19092\n    topic: otlp_metrics\n    protocol_version: 2.0.0\n\nprocessors:\n  batch:\n\nextensions:\n  health_check:\n    check_collector_pipeline:\n      enabled: true\n      interval: \"1m\"\n      exporter_failure_threshold: 5\n\n  pprof:\n    endpoint: :1888\n\n  zpages:\n    endpoint: :55679\n\n  basicauth/client:\n    client_auth:\n")
	config.WriteString("      username: \"" + sinkConfig.User + "\"\n")
	config.WriteString("      password: \"" + sinkConfig.Password + "\"\n")
	config.WriteString(
		"\nexporters:\n  prometheusremotewrite:\n    endpoint: \"" + sinkConfig.Url + "\"\n")
	config.WriteString(
		"    auth:\\n      authenticator: basicauth/client\\n  logging:\\n    loglevel: debug\\n\\nservice:\\n  extensions: [pprof, zpages, health_check, basicauth/client]\\n  pipelines:\\n    metrics:\\n      receivers: [kafka]\\n      exporters: [prometheusremotewrite]\\n\"")

	return config.String()
}
