/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prom

import (
	"fmt"

	mfconsumers "github.com/mainflux/mainflux/consumers"
	"github.com/mainflux/mainflux/logger"
	"github.com/ns1labs/orb/pkg/mainflux/transformers/passthrough"
	"github.com/ns1labs/orb/pkg/promremotewrite"

	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/pkg/sinks/writer"
)

type promSinkService struct {
	mfsdk      mfsdk.SDK
	mfconsumer mfconsumers.Consumer
	pWriterMgr promremotewrite.PromRemoteWriter
}

func (p promSinkService) Run() error {
	t := passthrough.New()
	if err = mfconsumers.Start(pubSub, p.mfconsumer, t, cfg.ConfigPath, logger); err != nil {
		logger.Error(fmt.Sprintf("Failed to create promsink writer: %s", err))
	}
}

// New instantiates the prom sink service implementation.
func New() writer.Service {
	return &promSinkService{
		pWriterMgr: promremotewrite.New(promremotewrite.PromRemoteConfig{}),
	}
}
