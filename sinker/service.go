/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinker

import (
	"context"
	"fmt"
	"time"

	"github.com/orb-community/orb/sinker/redis/consumer"
	"github.com/orb-community/orb/sinker/redis/producer"

	"github.com/go-kit/kit/metrics"
	"github.com/go-redis/redis/v8"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	fleetpb "github.com/orb-community/orb/fleet/pb"
	policiespb "github.com/orb-community/orb/policies/pb"
	"github.com/orb-community/orb/sinker/otel"
	"github.com/orb-community/orb/sinker/otel/bridgeservice"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
)

const (
	OtelMetricsTopic = "otlp.*.m.>"
)

type Service interface {
	// Start set up communication with the message bus to communicate with agents
	Start() error
	// Stop end communication with the message bus
	Stop() error
}

type SinkerService struct {
	pubSub                 mfnats.PubSub
	otel                   bool
	otelMetricsCancelFunct context.CancelFunc
	otelLogsCancelFunct    context.CancelFunc
	otelKafkaUrl           string

	inMemoryCacheExpiration time.Duration
	streamClient            *redis.Client
	cacheClient             *redis.Client
	sinkTTLSvc              producer.SinkerKeyService
	sinkActivitySvc         producer.SinkActivityProducer
	logger                  *zap.Logger

	hbTicker *time.Ticker
	hbDone   chan bool

	policiesClient policiespb.PolicyServiceClient
	fleetClient    fleetpb.FleetServiceClient
	sinksClient    sinkspb.SinkServiceClient

	requestGauge   metrics.Gauge
	requestCounter metrics.Counter

	messageInputCounter metrics.Counter
	cancelAsyncContext  context.CancelFunc
	asyncContext        context.Context
}

func (svc SinkerService) Start() error {
	ctx := context.WithValue(context.Background(), "routine", "async")
	ctx = context.WithValue(ctx, "cache_expiry", svc.inMemoryCacheExpiration)
	svc.asyncContext, svc.cancelAsyncContext = context.WithCancel(ctx)

	svc.sinkTTLSvc = producer.NewSinkerKeyService(svc.logger, svc.cacheClient)
	svc.sinkActivitySvc = producer.NewSinkActivityProducer(svc.logger, svc.streamClient, svc.sinkTTLSvc)
	// Create Handle and Listener to Redis Key Events
	sinkerIdleProducer := producer.NewSinkIdleProducer(svc.logger, svc.streamClient)
	sinkerKeyExpirationListener := consumer.NewSinkerKeyExpirationListener(svc.logger, svc.cacheClient, sinkerIdleProducer)
	err := sinkerKeyExpirationListener.SubscribeToKeyExpiration(svc.asyncContext)
	if err != nil {
		svc.logger.Error("error on starting otel, exiting")
		ctx.Done()
		svc.cancelAsyncContext()
		return err
	}

	err = svc.startOtel(svc.asyncContext)
	if err != nil {
		svc.logger.Error("error on starting otel, exiting")
		ctx.Done()
		svc.cancelAsyncContext()
		return err
	}

	return nil
}

func (svc SinkerService) startOtel(ctx context.Context) error {
	if svc.otel {
		var err error

		bridgeService := bridgeservice.NewBridgeService(svc.logger, svc.inMemoryCacheExpiration, svc.sinkActivitySvc,
			svc.policiesClient, svc.sinksClient, svc.fleetClient, svc.messageInputCounter)
		svc.otelMetricsCancelFunct, err = otel.StartOtelMetricsComponents(ctx, &bridgeService, svc.logger, svc.otelKafkaUrl, svc.pubSub)

		// starting Otel Logs components
		svc.otelLogsCancelFunct, err = otel.StartOtelLogsComponents(ctx, &bridgeService, svc.logger, svc.otelKafkaUrl, svc.pubSub)

		if err != nil {
			svc.logger.Error("error during StartOtelComponents", zap.Error(err))
			return err
		}
	}
	return nil
}

func (svc SinkerService) Stop() error {
	otelTopic := fmt.Sprintf("channels.*.%s", OtelMetricsTopic)
	if err := svc.pubSub.Unsubscribe(otelTopic); err != nil {
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
	streamsClient *redis.Client,
	cacheClient *redis.Client,
	policiesClient policiespb.PolicyServiceClient,
	fleetClient fleetpb.FleetServiceClient,
	sinksClient sinkspb.SinkServiceClient,
	otelKafkaUrl string,
	enableOtel bool,
	requestGauge metrics.Gauge,
	requestCounter metrics.Counter,
	inputCounter metrics.Counter,
	defaultCacheExpiration time.Duration,
) Service {
	return &SinkerService{
		inMemoryCacheExpiration: defaultCacheExpiration,
		logger:                  logger,
		pubSub:                  pubSub,
		streamClient:            streamsClient,
		cacheClient:             cacheClient,
		policiesClient:          policiesClient,
		fleetClient:             fleetClient,
		sinksClient:             sinksClient,
		requestGauge:            requestGauge,
		requestCounter:          requestCounter,
		messageInputCounter:     inputCounter,
		otel:                    enableOtel,
		otelKafkaUrl:            otelKafkaUrl,
	}
}
