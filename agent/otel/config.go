package otel

import (
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confignet"
	"net/url"
	"time"
)

// Config defines configuration for simple prometheus receiver.
type Config struct {
	config.ReceiverSettings `mapstructure:",squash"`
	httpConfig              `mapstructure:",squash"`
	confignet.TCPAddr       `mapstructure:",squash"`
	// CollectionInterval is the interval at which metrics should be collected
	CollectionInterval time.Duration `mapstructure:"collection_interval"`
	// MetricsPath the path to the metrics endpoint.
	MetricsPath string `mapstructure:"metrics_path"`
	// Params the parameters to the metrics endpoint.
	Params url.Values `mapstructure:"params,omitempty"`
	// Whether or not to use pod service account to authenticate.
	UseServiceAccount bool `mapstructure:"use_service_account"`
}

// TODO: Move to a common package for use by other receivers and also pull
// in other utilities from
// https://github.com/signalfx/signalfx-agent/blob/main/pkg/core/common/httpclient/http.go.
type httpConfig struct {
	// Whether not TLS is enabled
	TLSEnabled bool      `mapstructure:"tls_enabled"`
	TLSConfig  tlsConfig `mapstructure:"tls_config"`
}

// tlsConfig holds common TLS config options
type tlsConfig struct {
	// Path to the CA cert that has signed the TLS cert.
	CAFile string `mapstructure:"ca_file"`
	// Path to the client TLS cert to use for TLS required connections.
	CertFile string `mapstructure:"cert_file"`
	// Path to the client TLS key to use for TLS required connections.
	KeyFile string `mapstructure:"key_file"`
	// Whether or not to verify the exporter's TLS cert.
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
}
