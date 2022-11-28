/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-redis/redis/v8"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	fleetpb "github.com/ns1labs/orb/fleet/pb"
	policiespb "github.com/ns1labs/orb/policies/pb"
	"github.com/ns1labs/orb/sinker/backend/pktvisor"
	"github.com/ns1labs/orb/sinker/config"
	"github.com/ns1labs/orb/sinker/otel"
	"github.com/ns1labs/orb/sinker/otel/bridgeservice"
	"github.com/ns1labs/orb/sinker/prometheus"
	sinkspb "github.com/ns1labs/orb/sinks/pb"
	"go.uber.org/zap"
)

const (
	BackendMetricsTopic = "be.*.m.>"
	OtelMetricsTopic    = "otlp.*.m.>"
	MaxMsgPayloadSize   = 1048 * 1000
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

type SinkerService struct {
	pubSub          mfnats.PubSub
	otel            bool
	otelCancelFunct context.CancelFunc
	otelKafkaUrl    string

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

func (svc SinkerService) Start() error {
	svc.asyncContext, svc.cancelAsyncContext = context.WithCancel(context.WithValue(context.Background(), "routine", "async"))
	if !svc.otel {
		topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
		if err := svc.pubSub.Subscribe(topic, svc.handleMsgFromAgent); err != nil {
			return err
		}
		svc.logger.Info("started metrics consumer", zap.String("topic", topic))
	}

	svc.hbTicker = time.NewTicker(CheckerFreq)
	svc.hbDone = make(chan bool)
	go svc.checkSinker()

	err := svc.startOtel(svc.asyncContext)
	if err != nil {
		svc.logger.Error("error on starting otel, exiting")
		return err
	}

	return nil
}

func (svc SinkerService) startOtel(ctx context.Context) error {
	if svc.otel {
		var err error
		bridgeService := bridgeservice.NewBridgeService(svc.logger, svc.sinkerCache, svc.policiesClient, svc.fleetClient)
		svc.otelCancelFunct, err = otel.StartOtelComponents(ctx, &bridgeService, svc.logger, svc.otelKafkaUrl, svc.pubSub)
		if err != nil {
			svc.logger.Error("error during StartOtelComponents", zap.Error(err))
			return err
		}
	}
	return nil
}

func (svc SinkerService) Stop() error {
	if svc.otel {
		otelTopic := fmt.Sprintf("channels.*.%s", OtelMetricsTopic)
		if err := svc.pubSub.Unsubscribe(otelTopic); err != nil {
			return err
		}
	} else {
		topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
		if err := svc.pubSub.Unsubscribe(topic); err != nil {
			return err
		}
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
	otelKafkaUrl string,
	enableOtel bool,
	requestGauge metrics.Gauge,
	requestCounter metrics.Counter,
	inputCounter metrics.Counter,
) Service {

	pktvisor.Register(logger)
	return &SinkerService{
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
		otel:                enableOtel,
		otelKafkaUrl:        otelKafkaUrl,
	}
}
