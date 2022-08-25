package otlpmqttexporter

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/gob"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ns1labs/orb/fleet"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"net/http"
	"net/url"
	"runtime"
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
}

// Crete new exporter.
func newExporter(cfg config.Exporter, set component.ExporterCreateSettings) (*exporter, error) {
	oCfg := cfg.(*Config)

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
		config:    oCfg,
		logger:    set.Logger,
		userAgent: userAgent,
		settings:  set.TelemetrySettings,
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

func (e *exporter) pushTraces(_ context.Context, _ ptrace.Traces) error {
	return fmt.Errorf("not implemented")
}

// pushMetrics Exports metrics
func (e *exporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	tr := pmetricotlp.NewRequest()
	tr.SetMetrics(md)
	request, err := tr.MarshalProto()
	if err != nil {
		return consumererror.NewPermanent(err)
	}
	metricRPC := fleet.AgentMetricsRPC{
		SchemaVersion: fleet.CurrentRPCSchemaVersion,
		Func:          fleet.AgentMetricsRPCFunc,
		Payload: []fleet.AgentMetricsRPCPayload{
			{
				Format:    "otlp",
				BEVersion: e.config.PktVisorVersion,
				Data:      request,
			},
		},
	}
	return e.export(ctx, e.config.MetricsTopic, metricRPC)
}

func (e *exporter) pushLogs(_ context.Context, _ plog.Logs) error {
	return fmt.Errorf("not implemented")
}

func (e *exporter) export(_ context.Context, metricsTopic string, request fleet.AgentMetricsRPC) error {
	// convert metrics to interface
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(request.Payload); err != nil {
		e.logger.Error("Failed to encode metrics", zap.Error(err))
	}

	if token := e.config.Client.Publish(metricsTopic, 1, false, buf.Bytes()); token.Wait() && token.Error() != nil {
		e.logger.Error("error sending metrics RPC", zap.String("topic", metricsTopic), zap.Error(token.Error()))
		return token.Error()
	}
	e.logger.Info("scraped and published metrics", zap.String("topic", metricsTopic), zap.Int("payload_size_b", 0), zap.Int("batch_count", 0))

	return nil
}
