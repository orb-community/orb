package pktvisorreceiver

import (
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
	promconfig "github.com/prometheus/prometheus/config"
	"go.opentelemetry.io/collector/config/confignet"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	promexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
)

var (
	testLog, _ = zap.NewDevelopment()
)

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
			Namespace:        "test",
			Endpoint:         ":8787",
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

//TODO create a instance of pktvisor with mocked data to ensure the test
func TestEndToEndToPktvisor(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		if testing.Short() {
			t.Skip("This test can take a couple of seconds")
		}

		//1. Create the Prometheus scrape endpoint.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// 2. Create the Prometheus metrics exporter that'll receive and verify the metrics produced.
		exporterCfg := &promexporter.Config{
			ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
			Namespace:        "test",
			Endpoint:         ":8787",
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
		rcvCfg := &Config{
			ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
			TCPAddr: confignet.TCPAddr{
				Endpoint: defaultEndpoint,
			},
			MetricsPath:        defaultMetricsPath,
			CollectionInterval: 1 * time.Millisecond,
		}
		// 3.5 Create the Prometheus receiver and pass in the preivously created Prometheus exporter.
		pConfig, err := GetPrometheusConfig(rcvCfg)
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
		// We shall let the Prometheus exporter scrape the DropWizard mock server, at least 9 times.
		res, err := http.Get("http://localhost" + exporterCfg.Endpoint + "/metrics")
		if err != nil {
			t.Fatalf("Failed to scrape from the exporter: %v", err)
		}
		prometheusExporterScrape, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if len(prometheusExporterScrape) != 0 {
			t.Fatalf("Left-over unmatched Prometheus scrape content: %q\n", prometheusExporterScrape)
		}
	})
}

func TestOrbAgentContainer(t *testing.T) {
	file, err := os.CreateTemp("", "orb-agent-pktvisor-conf-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	data := []byte(`version: "1.0"
visor:
  taps:
    default_pcap:
      input_type: pcap
      config:
        iface: "mock"
        pcap_source: mock
  policies:
    # policy name and description
    mypolicy:
#      description: "a mock view of anycast traffic"
      kind: collection
      # input stream to create based on the given tap and optional filter config
      input:
        # this must reference a tap name, or application of the policy will fail
        tap: default_pcap
        input_type: pcap
        config:
          bpf: "tcp or udp"
      # stream handlers to attach to this input stream
      # these decide exactly which data to summarize and expose for collection
      handlers:
        # default configuration for the stream handlers
        window_config:
          num_periods: 5
          deep_sample_rate: 100
        modules:
          # the keys at this level are unique identifiers
          default_net:
            type: net
          default_dns:
            type: dns
#            window_config:
#              max_deep_sample: 75
          special_domain:
            type: dns
            config:
              qname_suffix: .mydomain.com
`)
	if _, err := file.Write(data); err != nil {
		log.Fatal(err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	//volume := fmt.Sprintf("%s:/etc/pktvisor/config.yaml", file.Name())
	//docker run --net=host -d ns1labs/pktvisor pktvisord eth0
	ro := dockertest.RunOptions{
		Repository: "ns1labs/pktvisor",
		Tag:        "latest-develop",
		Cmd: []string{
			//"-v",
			//volume,
			//"--rm",
			"--net=host",
			"-d",
			"pktvisord",
			"wlp0s20f3",
			//"--config",
			//"/etc/pktvisor/config.yaml",
			//"--admin-api",
		},
	}
	container, err := pool.RunWithOptions(&ro)
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	port := container.GetPort("10853/tcp")
	fmt.Sprintf(port)

	if err := pool.Retry(func() error {
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s/api/v1/policies/__all/metrics/prometheus", port), nil)

		if err != nil {
			fmt.Println(err)
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(string(body))
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	//if err := pool.Retry(func() error {
	//	url := fmt.Sprintf("host=localhost port=%s user=test dbname=test password=test sslmode=disable", port)
	//	db, err := sqlx.Open("postgres", url)
	//	if err != nil {
	//		return err
	//	}
	//	return db.Ping()
	//}); err != nil {
	//	log.Fatalf("Could not connect to docker: %s", err)
	//}

	//if db, err = postgres.Connect(dbConfig); err != nil {
	//	log.Fatalf("Could not setup test DB connection: %s", err)
	//}

	//testLog.Debug("connected to database")

	//code := t.Run()

	//db.Close()

	if err := pool.Purge(container); err != nil {
		log.Fatalf("Could not purge container: %s", err)
	}

	//os.Exit(code)
}

//func TestRespondsWithLove(t *testing.T) {
//
//	pool, err := dockertest.NewPool("")
//	require.NoError(t, err, "could not connect to Docker")
//
//	resource, err := pool.Run("docker-gs-ping", "latest", []string{})
//	require.NoError(t, err, "could not start container")
//
//	t.Cleanup(func() {
//		require.NoError(t, pool.Purge(resource), "failed to remove container")
//	})
//
//	var resp *http.Response
//
//	err = pool.Retry(func() error {
//		resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8080/tcp"), "/"))
//		if err != nil {
//			t.Log("container not ready, waiting...")
//			return err
//		}
//		return nil
//	})
//	require.NoError(t, err, "HTTP error")
//	defer resp.Body.Close()
//
//	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")
//
//	body, err := io.ReadAll(resp.Body)
//	require.NoError(t, err, "failed to read HTTP body")
//
//	// Finally, test the business requirement!
//	require.Contains(t, string(body), "<3", "does not respond with love?")
//}
