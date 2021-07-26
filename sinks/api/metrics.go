/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/ns1labs/orb/sinks"
	"github.com/ns1labs/orb/sinks/backend"
)

var _ sinks.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     sinks.Service
}

func (m metricsMiddleware) CreateSink(ctx context.Context, token string, s sinks.Sink) (sinks.Sink, error) {
	return m.svc.CreateSink(ctx, token, s)
}

func (m metricsMiddleware) UpdateSink(ctx context.Context, token string, s sinks.Sink) (err error) {
	return m.svc.UpdateSink(ctx, token, s)
}

func (m metricsMiddleware) ListSinks(ctx context.Context, token string, pm sinks.PageMetadata) (_ sinks.Page, err error) {
	return m.svc.ListSinks(ctx, token, pm)
}

func (m metricsMiddleware) ListBackends(ctx context.Context, token string) (_ []string, err error) {
	return m.svc.ListBackends(ctx, token)
}

func (m metricsMiddleware) ViewBackend(ctx context.Context, token string, key string) (_ backend.Backend, err error) {
	return m.svc.ViewBackend(ctx, token, key)
}

func (m metricsMiddleware) ViewSink(ctx context.Context, token string, key string) (_ sinks.Sink, err error) {
	return m.svc.ViewSink(ctx, token, key)
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc sinks.Service, counter metrics.Counter, latency metrics.Histogram) sinks.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
