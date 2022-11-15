package otlpmqttexporter

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ns1labs/orb/agent/otel"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for OTLP/HTTP exporter.
type Config struct {
	config.ExporterSettings      `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	exporterhelper.QueueSettings `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings `mapstructure:"retry_on_failure"`

	// Add Client directly to only re-use an existing connection - requires "github.com/eclipse/paho.mqtt.golang"
	Client *mqtt.Client

	// Configuration to connect to MQTT
	Address      string `mapstructure:"address"`
	Id           string `mapstructure:"id"`
	Key          string `mapstructure:"key"`
	ChannelID    string `mapstructure:"channel_id"`
	TLS          bool   `mapstructure:"enable_tls"`
	MetricsTopic string `mapstructure:"metrics_topic"`

	// Specific for ORB Agent
	PktVisorVersion string `mapstructure:"pktvisor_version"`
	OrbAgentService otel.AgentBridgeService
}

var _ config.Exporter = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if ((cfg.Address != "" && cfg.Id != "" && cfg.Key != "" && cfg.ChannelID != "") ||
		cfg.Client != nil) && cfg.MetricsTopic != "" {
		return nil
	}
	return fmt.Errorf("invalid mqtt configuration")
}
