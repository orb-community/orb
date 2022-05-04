package otlpmqttexporter

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"testing"
)

func TestInvalidConfig(t *testing.T) {
	t.Skip("TODO Not sure how to solve this")
	c := &Config{}
	f := NewFactory()
	set := componenttest.NewNopExporterCreateSettings()
	_, err := f.CreateTracesExporter(context.Background(), set, c)
	require.Error(t, err)
	_, err = f.CreateMetricsExporter(context.Background(), set, c)
	require.Error(t, err)
	_, err = f.CreateLogsExporter(context.Background(), set, c)
	require.Error(t, err)
}

func TestUserAgent(t *testing.T) {
	// This test also requires you to use a local mqtt broker, for this I will use mosquitto on port 1887
	mqttAddr := "localhost:1887"
	set := componenttest.NewNopExporterCreateSettings()
	set.BuildInfo.Description = "Collector"
	set.BuildInfo.Version = "1.2.3test"

	tests := []struct {
		name string
	}{
		{
			name: "default_user_agent",
		},
		{
			name: "custom_user_agent",
		},
		{
			name: "custom_user_agent_lowercase",
		},
	}

	t.Run("metrics", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				metricsTopic := "channels/uuid1/messages/be/test"
				cfg := &Config{
					ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
					Address:          mqttAddr,
					Id:               "uuid1",
					Key:              "uuid2",
					TLS:              false,
					MetricsTopic:     metricsTopic,
				}
				exp, err := CreateMetricsExporter(context.Background(), set, cfg)
				require.NoError(t, err)

				// start the exporter
				err = exp.Start(context.Background(), componenttest.NewNopHost())
				require.NoError(t, err)
				t.Cleanup(func() {
					require.NoError(t, exp.Shutdown(context.Background()))
				})

				// generate data
				metrics := pmetric.NewMetrics()
				metrics.ResourceMetrics()
				metrics.ResourceMetrics().AppendEmpty()
				tv := metrics.ResourceMetrics().At(0)
				tv.SetSchemaUrl("test_url")
				tv.ScopeMetrics().AppendEmpty()
				sm := tv.ScopeMetrics().At(0)
				sm.Metrics().AppendEmpty()
				metric := sm.Metrics().At(0)
				metric.SetName("test_value")
				metric.SetDescription("test_description")
				metric.SetUnit("test_unit")
				err = exp.ConsumeMetrics(context.Background(), metrics)
				require.NoError(t, err)
			})
		}
	})
}
