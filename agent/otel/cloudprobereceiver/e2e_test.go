package cloudprobereceiver_test

import (
	"context"
	"fmt"
	"github.com/ns1labs/orb/agent/otel/cloudprobereceiver"
	promconfig "github.com/prometheus/prometheus/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/confignet"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	promexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
)

// Test with mocked server
func TestEndToEndSummarySupport(t *testing.T) {
	t.Run("test", func(t *testing.T) {

		if testing.Short() {
			t.Skip("This test can take a couple of seconds")
		}

		//1. Create the Prometheus scrape endpoint.
		waitForScrape := make(chan bool, 1)
		shutdown := make(chan bool, 1)
		dropWizardServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			select {
			case <-shutdown:
				return
			case waitForScrape <- true:
				// Serve back the metrics as if they were from DropWizard.
				_, err := rw.Write([]byte(dropWizardResponse))
				require.NoError(t, err)
			}
		}))
		defer dropWizardServer.Close()
		defer close(shutdown)

		srvURL, err := url.Parse(dropWizardServer.URL)
		if err != nil {
			t.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// 2. Create the Prometheus metrics exporter that'll receive and verify the metrics produced.
		exporterCfg := &promexporter.Config{
			ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
			HTTPServerSettings: confighttp.HTTPServerSettings{
				Endpoint: ":8788",
			},
			Namespace:        "test_cloudprober",
			SendTimestamps:   true,
			MetricExpiration: 2 * time.Hour,
		}
		exporterFactory := promexporter.NewFactory()
		set := componenttest.NewNopExporterCreateSettings()
		exporter, err := exporterFactory.CreateMetricsExporter(ctx, set, exporterCfg)
		if err != nil {
			t.Fatal(err)
		}
		if err = exporter.Start(ctx, nil); err != nil {
			t.Fatalf("Failed to start the Prometheus receiver: %v", err)
		}
		t.Cleanup(func() { require.NoError(t, exporter.Shutdown(ctx)) })

		//3. Create the Prometheus receiver scraping from the DropWizard mock server and
		//it'll feed scraped and converted metrics then pass them to the Prometheus exporter.
		yamlConfig := []byte(fmt.Sprintf(`
       global:
         scrape_interval: 2ms
    
       scrape_configs:
           - job_name: 'otel-collector'
             scrape_interval: 2ms
             static_configs:
               - targets: ['%s']
       `, srvURL.Host))
		receiverConfig := new(promconfig.Config)
		if err = yaml.Unmarshal(yamlConfig, receiverConfig); err != nil {
			t.Fatal(err)
		}

		receiverFactory := prometheusreceiver.NewFactory()
		receiverCreateSet := componenttest.NewNopReceiverCreateSettings()
		rcvCfg := &prometheusreceiver.Config{
			PrometheusConfig: receiverConfig,
			ReceiverSettings: config.NewReceiverSettings(config.NewComponentID("prometheus")),
		}
		// 3.5 Create the Prometheus receiver and pass in the preivously created Prometheus exporter.
		prometheusReceiver, err := receiverFactory.CreateMetricsReceiver(ctx, receiverCreateSet, rcvCfg, exporter)
		if err != nil {
			t.Fatal(err)
		}

		if err = prometheusReceiver.Start(ctx, nil); err != nil {
			t.Fatalf("Failed to start the Prometheus receiver: %v", err)
		}
		t.Cleanup(func() { require.NoError(t, prometheusReceiver.Shutdown(ctx)) })

		// 4. Scrape from the Prometheus exporter to ensure that we export summary metrics
		// We shall let the Prometheus exporter scrape the DropWizard mock server, at least 9 times.
		for i := 0; i < 8; i++ {
			<-waitForScrape
		}

		res, err := http.Get("http://localhost" + exporterCfg.Endpoint + "/metrics")
		if err != nil {
			t.Fatalf("Failed to scrape from the exporter: %v", err)
		}
		prometheusExporterScrape, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if len(prometheusExporterScrape) == 0 {
			t.Fatalf("Left-over unmatched Prometheus scrape content: %q\n", prometheusExporterScrape)
		}
	})
}

const dropWizardResponse = `
# HELP jvm_memory_pool_bytes_used Used bytes of a given JVM memory pool.
# TYPE jvm_memory_pool_bytes_used gauge
jvm_memory_pool_bytes_used{pool="CodeHeap 'non-nmethods'",} 1277952.0
jvm_memory_pool_bytes_used{pool="Metaspace",} 2.6218176E7
jvm_memory_pool_bytes_used{pool="CodeHeap 'profiled nmethods'",} 6871168.0
jvm_memory_pool_bytes_used{pool="Compressed Class Space",} 2751312.0
jvm_memory_pool_bytes_used{pool="G1 Eden Space",} 4.4040192E7
jvm_memory_pool_bytes_used{pool="G1 Old Gen",} 4385408.0
jvm_memory_pool_bytes_used{pool="G1 Survivor Space",} 8388608.0
jvm_memory_pool_bytes_used{pool="CodeHeap 'non-profiled nmethods'",} 2869376.0
# HELP jvm_info JVM version info
# TYPE jvm_info gauge
jvm_info{version="9.0.4+11",vendor="Oracle Corporation",} 1.0
# HELP jvm_gc_collection_seconds Time spent in a given JVM garbage collector in seconds.
# TYPE jvm_gc_collection_seconds summary
jvm_gc_collection_seconds_count{gc="G1 Young Generation",} 9.0
jvm_gc_collection_seconds_sum{gc="G1 Young Generation",} 0.229
jvm_gc_collection_seconds_count{gc="G1 Old Generation",} 0.0
jvm_gc_collection_seconds_sum{gc="G1 Old Generation",} 0.0`

func TestEndToEndToCloudprober(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		if testing.Short() {
			t.Skip("This test can take a couple of seconds")
		}

		defaultEndpoint := container.GetHostPort("9313/tcp")

		//1. Create the Prometheus scrape endpoint.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// 2. Create the Prometheus metrics exporter that'll receive and verify the metrics produced.
		exporterCfg := &promexporter.Config{
			ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
			HTTPServerSettings: confighttp.HTTPServerSettings{
				Endpoint: ":8788",
			},
			Namespace:        "test",
			SendTimestamps:   true,
			MetricExpiration: 2 * time.Hour,
		}
		exporterFactory := promexporter.NewFactory()
		set := componenttest.NewNopExporterCreateSettings()
		exporter, err := exporterFactory.CreateMetricsExporter(ctx, set, exporterCfg)
		if err != nil {
			t.Fatal(err)
		}
		if err = exporter.Start(ctx, nil); err != nil {
			t.Fatalf("Failed to start the Prometheus receiver: %v", err)
		}
		t.Cleanup(func() { require.NoError(t, exporter.Shutdown(ctx)) })

		//3. Create the Prometheus receiver scraping from the DropWizard mock server and
		//it'll feed scraped and converted metrics then pass them to the Prometheus exporter.
		receiverFactory := prometheusreceiver.NewFactory()
		receiverCreateSet := componenttest.NewNopReceiverCreateSettings()
		rcvCfg := &cloudprobereceiver.Config{
			ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
			TCPAddr: confignet.TCPAddr{
				Endpoint: defaultEndpoint,
			},
			MetricsPath:        defaultMetricsPath,
			CollectionInterval: 1 * time.Second,
		}
		// 3.5 Create the Prometheus receiver and pass in the preivously created Prometheus exporter.
		pConfig, err := cloudprobereceiver.GetPrometheusConfig(rcvCfg)
		if err != nil {
			t.Fatalf("failed to create prometheus receiver config: %v", err)
		}
		prometheusReceiver, err := receiverFactory.CreateMetricsReceiver(ctx, receiverCreateSet, pConfig, exporter)
		if err != nil {
			t.Fatal(err)
		}

		if err = prometheusReceiver.Start(ctx, nil); err != nil {
			t.Fatalf("Failed to start the Prometheus receiver: %v", err)
		}
		t.Cleanup(func() { require.NoError(t, prometheusReceiver.Shutdown(ctx)) })

		// 4. Scrape from the Prometheus exporter to ensure that we export summary metrics
		// We shall let the Prometheus exporter scrape the pktvisor mock server, at least after 5 seconds.
		var res *http.Response
		var prometheusExporterScrape []byte
		var backoffSchedule = []time.Duration{
			5 * time.Second,
			10 * time.Second,
			15 * time.Second,
		}
		for _, backoff := range backoffSchedule {
			res, err = http.Get("http://localhost" + exporterCfg.Endpoint + "/metrics")
			if err != nil {
				t.Fatalf("Failed to scrape from the exporter: %v", err)
			}
			prometheusExporterScrape, err = ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if len(prometheusExporterScrape) == 0 {
				time.Sleep(backoff)
			}
		}
		if len(prometheusExporterScrape) == 0 {
			t.Fatalf("Left-over unmatched Prometheus scrape content: %q\n", prometheusExporterScrape)
		}
	})
}
