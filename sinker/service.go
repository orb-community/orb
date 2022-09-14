/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinker

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/go-redis/redis/v8"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	fleetpb "github.com/ns1labs/orb/fleet/pb"
	policiespb "github.com/ns1labs/orb/policies/pb"
	"github.com/ns1labs/orb/sinker/backend/pktvisor"
	"github.com/ns1labs/orb/sinker/config"
	"github.com/ns1labs/orb/sinker/otel"
	"github.com/ns1labs/orb/sinker/prometheus"
	sinkspb "github.com/ns1labs/orb/sinks/pb"
	promexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"go.opentelemetry.io/collector/component"
	otelconfig "go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

const (
	BackendMetricsTopic = "be.*.m.>"
	OtelMetricsTopic    = "otlp.*.m.>"
	MaxMsgPayloadSize   = 2048 * 100
)

var (
	ErrPayloadTooBig = errors.New("payload too big")
	ErrNotFound      = errors.New("non-existent entity")
)

type Service interface {
	// Start set up communication with the message bus to communicate with agents
	Start() error
	// Stop end communication with the message bus
	Stop() error
}

type sinkerService struct {
	pubSub             mfnats.PubSub
	otel               bool
	otelCancelFunct    context.CancelFunc
	otelMetricsChannel chan []byte

	sinkerCache config.ConfigRepo
	esclient    *redis.Client
	logger      *zap.Logger

	hbTicker *time.Ticker
	hbDone   chan bool

	promClient prometheus.Client

	policiesClient policiespb.PolicyServiceClient
	fleetClient    fleetpb.FleetServiceClient
	sinksClient    sinkspb.SinkServiceClient

	requestGauge   metrics.Gauge
	requestCounter metrics.Counter

	messageInputCounter metrics.Counter
	cancelAsyncContext  context.CancelFunc
	asyncContext        context.Context
}

func (svc sinkerService) Start() error {
	svc.asyncContext, svc.cancelAsyncContext = context.WithCancel(context.WithValue(context.Background(), "routine", "async"))
	topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
	if err := svc.pubSub.Subscribe(topic, svc.handleMsgFromAgent); err != nil {
		return err
	}
	//OtelMetricsTopic
	otelTopic := fmt.Sprintf("channels.*.%s", OtelMetricsTopic)
	if err := svc.pubSub.Subscribe(otelTopic, svc.handleOtelMsgFromAgent); err != nil {
		return err
	}
	svc.logger.Info("started metrics consumer", zap.String("topic", topic))

	svc.hbTicker = time.NewTicker(CheckerFreq)
	svc.hbDone = make(chan bool)
	svc.otelMetricsChannel = make(chan []byte)
	go svc.checkSinker()

	err := svc.startOtel(svc.asyncContext)
	if err != nil {
		return err
	}

	return nil
}

func (svc sinkerService) startOtel(ctx context.Context) error {
	var err error
	if svc.otel {
		svc.otelCancelFunct, err = otel.StartOtelComponents(ctx, svc.logger, &svc.otelMetricsChannel)
		if err != nil {
			svc.logger.Error("error during StartOtelComponents", zap.Error(err))
			return err
		}
	}
	return nil
}

func (svc sinkerService) Stop() error {
	topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
	if err := svc.pubSub.Unsubscribe(topic); err != nil {
		return err
	}
	svc.logger.Info("unsubscribed from agent metrics")

	svc.hbTicker.Stop()
	svc.hbDone <- true
	svc.cancelAsyncContext()

	return nil
}

// New instantiates the sinker service implementation.
func New(logger *zap.Logger,
	pubSub mfnats.PubSub,
	esclient *redis.Client,
	configRepo config.ConfigRepo,
	policiesClient policiespb.PolicyServiceClient,
	fleetClient fleetpb.FleetServiceClient,
	sinksClient sinkspb.SinkServiceClient,
	requestGauge metrics.Gauge,
	requestCounter metrics.Counter,
	inputCounter metrics.Counter,
) Service {

	pktvisor.Register(logger)
	return &sinkerService{
		logger:              logger,
		pubSub:              pubSub,
		esclient:            esclient,
		sinkerCache:         configRepo,
		policiesClient:      policiesClient,
		fleetClient:         fleetClient,
		sinksClient:         sinksClient,
		requestGauge:        requestGauge,
		requestCounter:      requestCounter,
		messageInputCounter: inputCounter,
		otel:                true,
	}
}

func createReceiver(ctx context.Context, logger *zap.Logger) (component.MetricsReceiver, error) {
	receiverFactory := otlpreceiver.NewFactory()

	set := component.ReceiverCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
			MetricsLevel:   configtelemetry.LevelDetailed,
		},
		BuildInfo: component.BuildInfo{},
	}
	metricsReceiver, err := receiverFactory.CreateMetricsReceiver(ctx, set,
		receiverFactory.CreateDefaultConfig(), consumertest.NewNop())
	return metricsReceiver, err
}

func createExporter(ctx context.Context, logger *zap.Logger) (component.MetricsExporter, error) {
	// 2. Create the Prometheus metrics exporter that'll receive and verify the metrics produced.
	exporterCfg := &promexporter.Config{
		ExporterSettings: otelconfig.NewExporterSettings(otelconfig.NewComponentID("pktvisor_prometheus_exporter")),
		HTTPServerSettings: confighttp.HTTPServerSettings{
			Endpoint: ":8787",
		},
		Namespace:        "test",
		SendTimestamps:   true,
		MetricExpiration: 2 * time.Hour,
	}
	exporterFactory := promexporter.NewFactory()
	set := component.ExporterCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	exporter, err := exporterFactory.CreateMetricsExporter(ctx, set, exporterCfg)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}
