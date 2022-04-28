package otlpmqttexporter

import (
	"context"
	"crypto/tls"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.uber.org/zap"
	"reflect"
	"testing"
	"time"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configtest"
)

func TestCreateDefaultConfig(t *testing.T) {
	t.Skip("Configuration is not done yet")
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, configtest.CheckConfigStruct(cfg))
	ocfg, ok := factory.CreateDefaultConfig().(*Config)
	assert.True(t, ok)
	assert.Equal(t, ocfg.Address, "localhost", "default address is localhost")
	assert.Equal(t, ocfg.Id, "uuid", "default id is uuid1")
	assert.Equal(t, ocfg.Key, "uuid", "default key uuid1")
	assert.Equal(t, ocfg.ChannelID, "agent_test_metrics", "default channel ID agent_test_metrics ")
	assert.Equal(t, ocfg.TLS, false, "default TLS is disabled")
	assert.Nil(t, ocfg.MetricsTopic, "default metrics topic is nil, only passed in the export function")
}

func TestCreateMetricsExporter(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)

	set := componenttest.NewNopExporterCreateSettings()
	oexp, err := factory.CreateMetricsExporter(context.Background(), set, cfg)
	require.Nil(t, err)
	require.NotNil(t, oexp)
}

func makeMQTTConnectedClient(t *testing.T) (client mqtt.Client, err error) {
	opts := mqtt.NewClientOptions().AddBroker("localhost:1889").SetClientID("1dad1121-4b05-4af8-9321-c541e252fe4b")
	opts.SetUsername("1dad1121-4b05-4af8-9321-c541e252fe4b")
	opts.SetPassword("2a2aabd8-927f-4c58-9dc4-2de784cf9644")
	opts.SetKeepAlive(10 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		t.Error("message on unknown channel, ignoring", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))
	})
	opts.SetPingTimeout(5 * time.Second)
	opts.SetAutoReconnect(true)

	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func TestCreateConfigClient(t *testing.T) {
	type args struct {
		client mqtt.Client
	}

	client, err := makeMQTTConnectedClient(t)
	require.Nil(t, err)

	tests := []struct {
		name string
		args args
		want config.Exporter
	}{
		{
			name: "ok client",
			args: args{
				client: client,
			},
			want: nil,
		},
		{
			name: "nil client",
			args: args{
				client: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateConfigClient(tt.args.client)
			assert.NoError(t, got.Validate())
		})
	}
}

func TestCreateDefaultSettings(t *testing.T) {
	type args struct {
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
		want component.ExporterCreateSettings
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateDefaultSettings(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDefaultSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}
