/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prom

import (
	mfconsumers "github.com/mainflux/mainflux/consumers"
	mflog "github.com/mainflux/mainflux/logger"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"

	"github.com/ns1labs/orb/pkg/mainflux/transformers/passthrough"
	"github.com/ns1labs/orb/pkg/promremotewrite"
	"github.com/ns1labs/orb/pkg/sinks/writer"
	"go.uber.org/zap"
)

type promSinkService struct {
	mflogger   mflog.Logger
	mfsdk      mfsdk.SDK
	mfconsumer mfconsumers.Consumer
	pubSub     mfnats.PubSub

	natSubjectConfigPath string
	logger               *zap.Logger
	pWriterMgr           promremotewrite.PromRemoteWriter
}

func (p promSinkService) Run() error {
	t := passthrough.New()
	if err := mfconsumers.Start(p.pubSub, p.mfconsumer, t, p.natSubjectConfigPath, p.mflogger); err != nil {
		p.logger.Error("Failed to create promsink writer", zap.Error(err))
	}
	return nil
}

// New instantiates the prom sink service implementation.
func New(logger *zap.Logger) writer.Service {
	return &promSinkService{
		logger:     logger,
		pWriterMgr: promremotewrite.New(promremotewrite.PromRemoteConfig{}),
	}
}
