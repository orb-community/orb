/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/go-kit/kit/metrics"
	"github.com/ns1labs/orb/pkg/fleet"
)

var _ fleet.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     fleet.Service
}

func (m metricsMiddleware) Add() (fleet.Agent, error) {
	panic("implement me")
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc fleet.Service, counter metrics.Counter, latency metrics.Histogram) fleet.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
