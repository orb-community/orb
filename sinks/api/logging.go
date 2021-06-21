/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/ns1labs/orb/sinks"
	"go.uber.org/zap"
)

var _ sinks.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *zap.Logger
	svc    sinks.Service
}

func (l loggingMiddleware) CreateAgent(ctx context.Context, token string, s sinks.Sink) (sinks.Sink, error) {
	panic("implement me")
}

func NewLoggingMiddleware(svc sinks.Service, logger *zap.Logger) sinks.Service {
	return &loggingMiddleware{logger, svc}
}
