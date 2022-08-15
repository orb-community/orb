package otlpexporter

import (
	"fmt"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type Config struct {
	config.ExporterSettings        `mapstructure:",squash"`
	exporterhelper.TimeoutSettings `mapstructure:",squash"`
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`
	configgrpc.GRPCClientSettings  `mapstructure:",squash"`
}

var _ config.Exporter = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if cfg.Endpoint == "" {
		return fmt.Errorf("at least one endpoint must be specified")
	}
	return nil
}


// ID gets the receiver name.
func (cfg *Config) ID() config.ComponentID {
	return cfg.ExporterSettings.ID()
}

// SetIDName sets the receiver name.
func (cfg *Config) SetIDName(idName string) {
	cfg.ExporterSettings.SetIDName(idName)
}
