package pktvisorreceiver

import (
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confignet"
	"net/url"
	"time"
)

type Config struct {
	config.ReceiverSettings `mapstructure:",squash"`
	confignet.TCPAddr       `mapstructure:",squash"`
	CollectionInterval      time.Duration `mapstructure:"collection_interval"`

	// MetricsPath the path to the metrics endpoint.
	MetricsPath string `mapstructure:"metrics_path"`
	// Params the parameters to the metrics endpoint.
	Params url.Values `mapstructure:"params,omitempty"`
	// Whether to use pod service account to authenticate.
	UseServiceAccount bool `mapstructure:"use_service_account"`
}
