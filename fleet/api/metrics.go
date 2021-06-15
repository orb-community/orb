/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/ns1labs/orb/fleet"
)

var _ fleet.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     fleet.Service
}

func (m metricsMiddleware) CreateAgent(ctx context.Context, token string, a fleet.Agent) (fleet.Agent, error) {
	return m.svc.CreateAgent(ctx, token, a)
}

func (m metricsMiddleware) CreateSelector(ctx context.Context, token string, s fleet.Selector) (fleet.Selector, error) {
	return m.svc.CreateSelector(ctx, token, s)
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc fleet.Service, counter metrics.Counter, latency metrics.Histogram) fleet.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
