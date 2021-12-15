package pktvisorreceiver_test

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"testing"
)

const (
	// The value of "type" key in configuration.
	typeStr            = "pktvisor_prometheus"
	defaultMetricsPath = "/api/v1/policies/__all/metrics/prometheus"
)

var (
	testLog, _ = zap.NewDevelopment()
	container  *dockertest.Resource
)

func TestMain(m *testing.M) {
	file, err := os.CreateTemp("", "pktvisor-conf-")
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
	ro := dockertest.RunOptions{
		Repository:   "ns1labs/pktvisor",
		Tag:          "latest-develop",
		Name:         "pktvisord",
		Hostname:     "pktvisor",
		ExposedPorts: []string{"10853/tcp"},
		Mounts:       []string{fmt.Sprintf("%s:/etc/pktvisor/config.yaml", file.Name())},
		Cmd: []string{
			"pktvisord",
			"--config",
			"/etc/pktvisor/config.yaml",
			"--admin-api",
			"-l",
			"pktvisor"},
	}
	err = pool.RemoveContainerByName("pktvisord")
	if err != nil {
		log.Fatalf("Could not remove existing container: %s", err)
	}

	container, err = pool.RunWithOptions(&ro)
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	host := container.GetHostPort("10853/tcp")

	var res *http.Response
	var client *http.Client
	if err := pool.Retry(func() error {
		client = &http.Client{}
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/api/v1/policies/__all/metrics/prometheus", host), nil)

		if err != nil {
			fmt.Println(err)
			return err
		}
		res, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil

	}); err != nil {
		if err := pool.Purge(container); err != nil {
			log.Fatalf("Could not purge container: %s", err)
		}
		log.Fatalf("Could not connect to docker: %s", err)
	}
	defer res.Body.Close()

	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(body))

	testLog.Debug("pktvisor running")

	code := m.Run()

	if err := pool.Purge(container); err != nil {
		log.Fatalf("Could not purge container: %s", err)
	}

	os.Exit(code)
}
