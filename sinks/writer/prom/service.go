/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prom

import (
	"github.com/go-redis/redis"
	mfconsumers "github.com/mainflux/mainflux/consumers"
	mflog "github.com/mainflux/mainflux/logger"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	esconsume "github.com/ns1labs/orb/sinks/redis/consumer"

	"github.com/ns1labs/orb/pkg/mainflux/transformers/passthrough"
	"github.com/ns1labs/orb/promremotewrite"
	"github.com/ns1labs/orb/sinks/writer"
	"go.uber.org/zap"
)

type promSinkService struct {
	mflogger   mflog.Logger
	mfconsumer mfconsumers.Consumer
	pubSub     mfnats.PubSub

	esclient *redis.Client

	natSubjectConfigPath string
	logger               *zap.Logger
	pWriterMgr           promremotewrite.PromRemoteWriter
}

func (p promSinkService) Run() error {
	t := passthrough.New()
	if err := mfconsumers.Start(p.pubSub, p.mfconsumer, t, p.natSubjectConfigPath, p.mflogger); err != nil {
		p.logger.Error("Failed to create promsink writer", zap.Error(err))
	}
	p.logger.Info("started nats consumer")
	eventStore := esconsume.NewEventStore(p, p.esclient, "prom-sink", p.logger)
	p.logger.Info("Subscribed to Redis Event Store")
	if err := eventStore.Subscribe("orb.policy"); err != nil {
		p.logger.Warn("orb prometheus sync service failed to subscribe to event sourcing: %s", zap.Error(err))
	}
	return nil
}

// New instantiates the prom sink service implementation.
func New(logger *zap.Logger, mflogger mflog.Logger, mfconsumer mfconsumers.Consumer, pubSub mfnats.PubSub, esclient *redis.Client, natSubjectConfigPath string) writer.Service {
	return &promSinkService{
		mflogger:             mflogger,
		mfconsumer:           mfconsumer,
		logger:               logger,
		pubSub:               pubSub,
		esclient:             esclient,
		natSubjectConfigPath: natSubjectConfigPath,
		pWriterMgr:           promremotewrite.New(promremotewrite.PromRemoteConfig{}),
	}
}
