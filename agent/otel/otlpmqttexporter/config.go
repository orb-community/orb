package otlpmqttexporter

import (
	"fmt"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for OTLP/HTTP exporter.
type Config struct {
	config.ExporterSettings       `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	confighttp.HTTPClientSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings  `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings  `mapstructure:"retry_on_failure"`

	// Configuration to connect to MQTT
	Address      string `mapstructure:"address"`
	Id           string `mapstructure:"id"`
	Key          string `mapstructure:"key"`
	ChannelID    string `mapstructure:"channel_id"`
	TLS          bool   `mapstructure:"enable_tls"`
	MetricsTopic string `mapstructure:"metrics_topic"`
}

var _ config.Exporter = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if cfg.Address == "" && cfg.Id == "" && cfg.Key == "" && cfg.ChannelID == "" {
		return fmt.Errorf("invalid mqtt configuration")
	}
	return nil
}
