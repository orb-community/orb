package otlpmqttexporter

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type exporter struct {
	// Input configuration.
	config     *Config
	client     *http.Client
	tracesURL  string
	metricsURL string
	logsURL    string
	logger     *zap.Logger
	settings   component.TelemetrySettings
	// Default user-agent header.
	userAgent string
	// Policy handled by this exporter
	policyID   string
	policyName string
}

func (e *exporter) compressBrotli(data []byte) []byte {
	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, brotli.BestCompression)
	_, err := w.Write(data)
	if err != nil {
		return nil
	}
	err = w.Close()
	if err != nil {
		return nil
	}
	return b.Bytes()
}

// Crete new exporter.
func newExporter(cfg config.Exporter, set component.ExporterCreateSettings, ctx context.Context) (*exporter, error) {
	oCfg := cfg.(*Config)
	policyID := ctx.Value("policy_id").(string)
	policyName := ctx.Value("policy_name").(string)
	if oCfg.Address != "" {
		_, err := url.Parse(oCfg.Address)
		if err != nil {
			return nil, errors.New("address must be a valid mqtt server")
		}
	}

	userAgent := fmt.Sprintf("%s/%s (%s/%s)",
		set.BuildInfo.Description, set.BuildInfo.Version, runtime.GOOS, runtime.GOARCH)

	// Client construction is deferred to start
	return &exporter{
		config:     oCfg,
		logger:     set.Logger,
		userAgent:  userAgent,
		settings:   set.TelemetrySettings,
		policyID:   policyID,
		policyName: policyName,
	}, nil
}

// start actually creates the MQTT client.
func (e *exporter) start(_ context.Context, _ component.Host) error {
	token := e.config.Client
	if token == nil {
		opts := mqtt.NewClientOptions().AddBroker(e.config.Address).SetClientID(e.config.Id)
		opts.SetUsername(e.config.Id)
		opts.SetPassword(e.config.Key)
		opts.SetKeepAlive(10 * time.Second)
		opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
			e.logger.Info("message on unknown channel, ignoring", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))
		})
		opts.SetPingTimeout(5 * time.Second)
		opts.SetAutoReconnect(true)

		if e.config.TLS {
			opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		e.config.Client = client
	}

	return nil
}

// extractAttribute extract attribute from metricsRequest metrics
func (e *exporter) extractAttribute(metricsRequest pmetricotlp.Request, attribute string) string {
	metrics := metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
	for i := 0; i < metrics.Len(); i++ {
		metricItem := metrics.At(i)
		if metricItem.Name() == "dns_wire_packets_tcp" || metricItem.Name() == "packets_ipv4" || metricItem.Name() == "dhcp_wire_packets_ack" || metricItem.Name() == "flow_in_udp_bytes" {
			p, _ := metricItem.Gauge().DataPoints().At(0).Attributes().Get(attribute)
			if p.AsString() != "" {
				return p.AsString()
			}
		}
	}
	return ""
}

// inject attribute on all metricsRequest metrics
func (e *exporter) injectAttribute(metricsRequest pmetricotlp.Request, attribute string, value string) pmetricotlp.Request {
	metrics := metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
	for i := 0; i < metrics.Len(); i++ {
		metricItem := metrics.At(i)

		if metricItem.Type().String() == "Gauge" {
			metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Gauge().DataPoints().At(0).Attributes().PutStr(attribute, value)
		} else if metricItem.Type().String() == "Summary" {
			metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Summary().DataPoints().At(0).Attributes().PutStr(attribute, value)
		} else if metricItem.Type().String() == "Histogram" {
			metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Histogram().DataPoints().At(0).Attributes().PutStr(attribute, value)
		} else if metricItem.Type().String() == "ExponentialHistogram" {
			metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr(attribute, value)
		} else {
			e.logger.Error("Unknown metric type: " + metricItem.Type().String())
		}
	}
	return metricsRequest
}

func (e *exporter) pushTraces(_ context.Context, _ ptrace.Traces) error {
	return fmt.Errorf("not implemented")
}

// pushMetrics Exports metrics
func (e *exporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	tr := pmetricotlp.NewRequestFromMetrics(md)

	agentData, err := e.config.OrbAgentService.RetrieveAgentInfoByPolicyID(e.policyID)
	if err != nil {
		defer ctx.Done()
		return consumererror.NewPermanent(err)
	}
	// injecting policy ID attribute on metrics
	tr = e.injectAttribute(tr, "policy_id", e.policyID)
	tr = e.injectAttribute(tr, "dataset_ids", agentData.Datasets)
	tr = e.injectAttribute(tr, "orb_tags", agentData.OrbTags)
	tr = e.injectAttribute(tr, "agent_tags", agentData.AgentTags)

	request, err := tr.MarshalProto()
	if err != nil {
		defer ctx.Done()
		return consumererror.NewPermanent(err)
	}

	e.logger.Info("request metrics count: " + strconv.Itoa(md.MetricCount()) + ", policyID: " + e.policyID)

	err = e.export(ctx, e.config.MetricsTopic, request)
	if err != nil {
		defer ctx.Done()
		return err
	}
	return err
}

func (e *exporter) pushLogs(_ context.Context, _ plog.Logs) error {
	return fmt.Errorf("not implemented")
}

func (e *exporter) export(_ context.Context, metricsTopic string, request []byte) error {
	compressedPayload := e.compressBrotli(request)
	if token := e.config.Client.Publish(metricsTopic, 1, false, compressedPayload); token.Wait() && token.Error() != nil {
		e.logger.Error("error sending metrics RPC", zap.String("topic", metricsTopic), zap.Error(token.Error()))
		return token.Error()
	}
	e.logger.Info("scraped and published metrics", zap.String("topic", metricsTopic), zap.Int("payload_size_b", len(request)), zap.Int("compressed_payload_size_b", len(compressedPayload)))

	return nil
}
