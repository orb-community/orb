/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package consumer

import (
	mfconsumers "github.com/mainflux/mainflux/consumers"
	"github.com/mainflux/mainflux/pkg/messaging"
	"go.uber.org/zap"
)

var _ mfconsumers.Consumer = (*fleetConsumer)(nil)

type fleetConsumer struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) mfconsumers.Consumer {
	logger.Info("created fleet agent nats consumer")
	return &fleetConsumer{logger: logger}
}

func (p fleetConsumer) Consume(messages interface{}) error {
	p.logger.Info("fleet agent consume", zap.String("subtopic", messages.(messaging.Message).Subtopic))
	return nil
}
