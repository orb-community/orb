/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/ns1labs/orb/sinks"
	"go.uber.org/zap"
	"time"
)

var _ sinks.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *zap.Logger
	svc    sinks.Service
}

func (l loggingMiddleware) ListSinks(ctx context.Context, token string, pm sinks.PageMetadata) (_ sinks.Page, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: list_sinks",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: list_sinks",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ListSinks(ctx, token, pm)
}

func (l loggingMiddleware) CreateSink(ctx context.Context, token string, s sinks.Sink) (_ sinks.Sink, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: create_sink",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: create_sink",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.CreateSink(ctx, token, s)
}

func (l loggingMiddleware) ListBackends(ctx context.Context, token string)(_ []string, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: list_backends",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: list_backends",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ListBackends(ctx, token)
}

func (l loggingMiddleware) DeleteSink(ctx context.Context, token string, id string) (err error)  {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: delete_sink",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: delete_sink",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())

	return l.svc.DeleteSink(ctx, token, id)
}

func NewLoggingMiddleware(svc sinks.Service, logger *zap.Logger) sinks.Service {
	return &loggingMiddleware{logger, svc}
}
