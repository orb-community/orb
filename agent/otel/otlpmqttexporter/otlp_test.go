package otlpmqttexporter

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

func TestInvalidConfig(t *testing.T) {
	config := &Config{}
	f := NewFactory()
	set := componenttest.NewNopExporterCreateSettings()
	_, err := f.CreateTracesExporter(context.Background(), set, config)
	require.Error(t, err)
	_, err = f.CreateMetricsExporter(context.Background(), set, config)
	require.Error(t, err)
	_, err = f.CreateLogsExporter(context.Background(), set, config)
	require.Error(t, err)
}

func TestMetricsError(t *testing.T) {
	addr := "localhost"

	startMetricsReceiver(t, addr, consumertest.NewErr(errors.New("my_error")))
	exp := startMetricsExporter(t, "", fmt.Sprintf("http://%s/v1/metrics", addr))

	md := pmetric.NewMetrics()
	assert.Error(t, exp.ConsumeMetrics(context.Background(), md))
}

func TestMetricsRoundTrip(t *testing.T) {
	addr := "localhost"

	tests := []struct {
		name        string
		baseURL     string
		overrideURL string
	}{
		{
			name:        "wrongbase",
			baseURL:     "http://wronghostname",
			overrideURL: fmt.Sprintf("http://%s/v1/metrics", addr),
		},
		{
			name:        "onlybase",
			baseURL:     fmt.Sprintf("http://%s", addr),
			overrideURL: "",
		},
		{
			name:        "override",
			baseURL:     "",
			overrideURL: fmt.Sprintf("http://%s/v1/metrics", addr),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sink := new(consumertest.MetricsSink)
			startMetricsReceiver(t, addr, sink)
			exp := startMetricsExporter(t, test.baseURL, test.overrideURL)

			md := pmetric.NewMetrics()

			assert.NoError(t, exp.ConsumeMetrics(context.Background(), md))
			require.Eventually(t, func() bool {
				return sink.DataPointCount() > 0
			}, 1*time.Second, 10*time.Millisecond)
			allMetrics := sink.AllMetrics()
			require.Len(t, allMetrics, 1)
			assert.EqualValues(t, md, allMetrics[0])
		})
	}
}

func startMetricsExporter(t *testing.T, baseURL string, overrideURL string) component.MetricsExporter {
	factory := NewFactory()
	cfg := createExporterConfig(baseURL, factory.CreateDefaultConfig())
	cfg.MetricsTopic = overrideURL
	exp, err := factory.CreateMetricsExporter(context.Background(), componenttest.NewNopExporterCreateSettings(), cfg)
	require.NoError(t, err)
	startAndCleanup(t, exp)
	return exp
}

func createExporterConfig(baseURL string, defaultCfg config.Exporter) *Config {
	cfg := defaultCfg.(*Config)
	cfg.Endpoint = baseURL
	cfg.QueueSettings.Enabled = false
	cfg.RetrySettings.Enabled = false
	return cfg
}

func startMetricsReceiver(t *testing.T, addr string, next consumer.Metrics) {
	factory := otlpreceiver.NewFactory()
	cfg := createReceiverConfig(addr, factory.CreateDefaultConfig())
	recv, err := factory.CreateMetricsReceiver(context.Background(), componenttest.NewNopReceiverCreateSettings(), cfg, next)
	require.NoError(t, err)
	startAndCleanup(t, recv)
}

func createReceiverConfig(addr string, defaultCfg config.Receiver) *otlpreceiver.Config {
	cfg := defaultCfg.(*otlpreceiver.Config)
	cfg.HTTP.Endpoint = addr
	cfg.GRPC = nil
	return cfg
}

func startAndCleanup(t *testing.T, cmp component.Component) {
	require.NoError(t, cmp.Start(context.Background(), componenttest.NewNopHost()))
	t.Cleanup(func() {
		require.NoError(t, cmp.Shutdown(context.Background()))
	})
}

func TestUserAgent(t *testing.T) {
	addr := "localhost"
	set := componenttest.NewNopExporterCreateSettings()
	set.BuildInfo.Description = "Collector"
	set.BuildInfo.Version = "1.2.3test"

	tests := []struct {
		name       string
		headers    map[string]string
		expectedUA string
	}{
		{
			name:       "default_user_agent",
			expectedUA: "Collector/1.2.3test",
		},
		{
			name:       "custom_user_agent",
			headers:    map[string]string{"User-Agent": "My Custom Agent"},
			expectedUA: "My Custom Agent",
		},
		{
			name:       "custom_user_agent_lowercase",
			headers:    map[string]string{"user-agent": "My Custom Agent"},
			expectedUA: "My Custom Agent",
		},
	}

	t.Run("metrics", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				mux := http.NewServeMux()
				mux.HandleFunc("/v1/metrics", func(writer http.ResponseWriter, request *http.Request) {
					assert.Contains(t, request.Header.Get("user-agent"), test.expectedUA)
					writer.WriteHeader(200)
				})
				srv := http.Server{
					Addr:    addr,
					Handler: mux,
				}
				ln, err := net.Listen("tcp", addr)
				require.NoError(t, err)
				go func() {
					_ = srv.Serve(ln)
				}()

				cfg := &Config{
					ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
					Address:          fmt.Sprintf("http://%s/v1/metrics", addr),
					HTTPClientSettings: confighttp.HTTPClientSettings{
						Headers: test.headers,
					},
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
				err = exp.ConsumeMetrics(context.Background(), metrics)
				require.NoError(t, err)

				srv.Close()
			})
		}
	})
}
