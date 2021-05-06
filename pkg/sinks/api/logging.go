/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/sinks"
	"go.uber.org/zap"
)

var _ sinks.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *zap.Logger
	svc    sinks.Service
}

func (l loggingMiddleware) Add() (sinks.Sink, error) {
	panic("implement me")
}

func NewLoggingMiddleware(svc sinks.Service, logger *zap.Logger) sinks.Service {
	return &loggingMiddleware{logger, svc}
}
