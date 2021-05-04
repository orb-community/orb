/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package promsink

import (
	"github.com/mainflux/mainflux/consumers"
	"github.com/mainflux/mainflux/pkg/messaging"
	"go.uber.org/zap"
)

type prometheusRepo struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) consumers.Consumer {
	logger.Info("created promsink")
	return &prometheusRepo{logger: logger}
}

func (p prometheusRepo) Consume(messages interface{}) error {
	p.logger.Info("promsink consume", zap.String("subtopic", messages.(messaging.Message).Subtopic))
	return nil
}
