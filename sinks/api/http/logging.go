/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"github.com/ns1labs/orb/sinks"
	"github.com/ns1labs/orb/sinks/backend"
	"go.uber.org/zap"
	"time"
)

var _ sinks.SinkService = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *zap.Logger
	svc    sinks.SinkService
}

func (l loggingMiddleware) ChangeSinkStateInternal(ctx context.Context, sinkID string, msg string, ownerID string, state sinks.State) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: change_sink_state_internal",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: change_sink_state_internal",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ChangeSinkStateInternal(ctx, sinkID, msg, ownerID, state)
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

func (l loggingMiddleware) UpdateSink(ctx context.Context, token string, s sinks.Sink) (sink sinks.Sink, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: edit_sink",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: edit_sink",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.UpdateSink(ctx, token, s)
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

func (l loggingMiddleware) ListBackends(ctx context.Context, token string) (_ []string, err error) {
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

func (l loggingMiddleware) ViewBackend(ctx context.Context, token string, key string) (_ backend.Backend, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_backend",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: view_backend",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewBackend(ctx, token, key)
}

func (l loggingMiddleware) ViewSink(ctx context.Context, token string, key string) (_ sinks.Sink, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_sink",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Warn("method call: view_sink",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewSink(ctx, token, key)
}

func (l loggingMiddleware) ViewSinkInternal(ctx context.Context, ownerID string, key string) (_ sinks.Sink, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_sink_internal",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Warn("method call: view_sink_internal",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewSinkInternal(ctx, ownerID, key)
}

func (l loggingMiddleware) DeleteSink(ctx context.Context, token string, key string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: delete_sink",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Warn("method call: delete_sink",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.DeleteSink(ctx, token, key)
}

func (l loggingMiddleware) ValidateSink(ctx context.Context, token string, s sinks.Sink) (_ sinks.Sink, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: validate_sink",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: validate_sink",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ValidateSink(ctx, token, s)
}

func NewLoggingMiddleware(svc sinks.SinkService, logger *zap.Logger) sinks.SinkService {
	return &loggingMiddleware{logger, svc}
}
