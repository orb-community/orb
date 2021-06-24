/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/ns1labs/orb/policies"
)

var _ policies.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     policies.Service
}

func (m metricsMiddleware) CreatePolicy(ctx context.Context, token string, p policies.Policy, format string, policyData string) (policies.Policy, error) {
	return m.svc.CreatePolicy(ctx, token, p, format, policyData)
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc policies.Service, counter metrics.Counter, latency metrics.Histogram) policies.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
