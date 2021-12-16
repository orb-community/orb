package otlpexporter

import (
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type Config struct {
	config.ExporterSettings        `mapstructure:",squash"`
	exporterhelper.TimeoutSettings `mapstructure:",squash"`
	exporterhelper.QueueSettings   `mapstructure:",sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:",retry_on_failure"`

	configgrpc.GRPCClientSettings `mapstructure:",squash"`
}

var _ config.Exporter = (*Config)(nil)
