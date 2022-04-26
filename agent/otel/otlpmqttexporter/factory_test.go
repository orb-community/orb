package otlpmqttexporter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

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
