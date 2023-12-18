package otlpmqttexporter

import (
	"go.opentelemetry.io/collector/component"
)

type CollectorConfig struct {
	Receivers  map[component.ID]component.Config `mapstructure:"receivers"`
	Extensions map[string]interface{}            `mapstructure:"extensions,omitempty"`
	Exporters  map[string]interface{}            `mapstructure:"exporters,omitempty"`
	Service    map[string]interface{}            `mapstructure:"service,omitempty"`
}

type HTTPCheckReceiver struct {
	CollectionInterval string            `mapstructure:"collection_interval"`
	Targets            []HTTPCheckTarget `mapstructure:"targets"`
}

type HTTPCheckTarget struct {
	Endpoint string            `mapstructure:"endpoint"`
	Method   string            `mapstructure:"method"`
	Tags     map[string]string `mapstructure:"tags,omitempty"`
}
