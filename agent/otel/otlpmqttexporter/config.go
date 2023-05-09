package otlpmqttexporter

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/orb-community/orb/agent/otel"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for OTLP/HTTP exporter.
type Config struct {
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	// Add Client directly to only re-use an existing connection - requires "github.com/eclipse/paho.mqtt.golang"
	Client *mqtt.Client

	// Configuration to connect to MQTT
	Address   string `mapstructure:"address"`
	Id        string `mapstructure:"id"`
	Key       string `mapstructure:"key"`
	ChannelID string `mapstructure:"channel_id"`
	TLS       bool   `mapstructure:"enable_tls"`
	Topic     string `mapstructure:"topic"`

	// Specific for ORB Agent
	PktVisorVersion string `mapstructure:"pktvisor_version"`
	OrbAgentService otel.AgentBridgeService
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if ((cfg.Address != "" && cfg.Id != "" && cfg.Key != "" && cfg.ChannelID != "") ||
		cfg.Client != nil) && cfg.Topic != "" {
		return nil
	}
	return fmt.Errorf("invalid mqtt configuration")
}
