package otlpmqttexporter

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.uber.org/zap"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	testedCfg, ok := factory.CreateDefaultConfig().(*Config)
	assert.True(t, ok)
	assert.Equal(t, "localhost:1883", testedCfg.Address, "default address is localhost")
	assert.Equal(t, "uuid1", testedCfg.Id, "default id is uuid1")
	assert.Equal(t, "uuid2", testedCfg.Key, "default key uuid1")
	assert.Equal(t, "channels/uuid1/messages", testedCfg.ChannelID, "default channel ID agent_test_metrics ")
	assert.False(t, testedCfg.TLS, "default TLS is disabled")
	assert.Equal(t, "channels/uuid1/messages/otlp/pktvisor", testedCfg.Topic, "default metrics topic is nil, only passed in the export function")
}

func TestCreateMetricsExporter(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)

	set := exportertest.NewNopCreateSettings()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "policy_name", "test")
	ctx = context.WithValue(ctx, "policy_id", "test")
	ctx = context.WithValue(ctx, "all", false)
	oexp, err := factory.CreateMetricsExporter(ctx, set, cfg)
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
	t.Skip("This test requires a local mqtt broker, unskip it locally")
	type args struct {
		client       mqtt.Client
		metricsTopic string
	}

	client, err := makeMQTTConnectedClient(t)
	require.Nil(t, err)

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "ok client",
			args: args{
				client:       client,
				metricsTopic: "topic",
			},
			want: nil,
		},
		{
			name: "nil client",
			args: args{
				client:       nil,
				metricsTopic: "",
			},
			want: fmt.Errorf("invalid mqtt configuration"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tt.args.client
			got := CreateConfigClient(c, tt.args.metricsTopic, " 1.0", nil)
			assert.Equal(t, tt.want, component.ValidateConfig(got), "expected %s but got %s", tt.want, component.ValidateConfig(got))
		})
	}
}

func TestCreateConfig(t *testing.T) {
	t.Skip(" only run this if local mqtt is installed locally at port 1889")
	type args struct {
		addr    string
		id      string
		key     string
		channel string
	}
	tests := []struct {
		name string
		args args
		want component.Config
	}{
		{
			name: "local mqtt",
			args: args{
				addr:    "localhost:1889",
				id:      "uuid1",
				key:     "uuid1",
				channel: "channels/uuid1/channel/metrics",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CreateConfig(tt.args.addr, tt.args.id, tt.args.key, tt.args.channel, "1.0", "metricstopic", nil), "CreateConfig(%v, %v, %v, %v)", tt.args.addr, tt.args.id, tt.args.key, tt.args.channel)
		})
	}
}
