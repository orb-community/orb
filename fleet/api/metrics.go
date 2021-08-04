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

func (m metricsMiddleware) RetrieveAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	return m.svc.RetrieveAgentGroupByIDInternal(ctx, groupID, ownerID)
}

func (m metricsMiddleware) RetrieveAgentGroupByID(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	return m.svc.RetrieveAgentGroupByID(ctx, groupID, ownerID)
}

func (m metricsMiddleware) ListAgents(ctx context.Context, token string, pm fleet.PageMetadata) (fleet.Page, error) {
	return m.svc.ListAgents(ctx, token, pm)
}

func (m metricsMiddleware) CreateAgent(ctx context.Context, token string, a fleet.Agent) (fleet.Agent, error) {
	return m.svc.CreateAgent(ctx, token, a)
}

func (m metricsMiddleware) CreateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (fleet.AgentGroup, error) {
	return m.svc.CreateAgentGroup(ctx, token, s)
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc fleet.Service, counter metrics.Counter, latency metrics.Histogram) fleet.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
