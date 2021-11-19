/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/ns1labs/orb/fleet"
	"time"
)

var _ fleet.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     fleet.Service
}

func (m metricsMiddleware) ViewOwnerByChannelIDInternal(ctx context.Context, channelID string) (agent fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewOwnerByChannelIDInternal",
			"owner_id", agent.MFOwnerID,
			"service_id", agent.MFThingID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ViewOwnerByChannelIDInternal(ctx, channelID)
}

func (m metricsMiddleware) ViewAgentBackend(ctx context.Context, token string, name string) (interface{}, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentBackend",
			"owner_id", "",
			"service_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ViewAgentBackend(ctx, token, name)
}

func (m metricsMiddleware) ListAgentBackends(ctx context.Context, token string) ([]string, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "listAgentBackends",
			"owner_id", "",
			"service_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ListAgentBackends(ctx, token)
}

func (m metricsMiddleware) ViewAgentByIDInternal(ctx context.Context, ownerID string, thingID string) (fleet.Agent, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentByIDInternal",
			"owner_id", ownerID,
			"service_id", thingID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ViewAgentByIDInternal(ctx, ownerID, thingID)
}

func (m metricsMiddleware) ViewAgentByID(ctx context.Context, token string, thingID string) (a fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentByID",
			"owner_id", a.MFOwnerID,
			"service_id", thingID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ViewAgentByID(ctx, token, thingID)
}

func (m metricsMiddleware) EditAgent(ctx context.Context, token string, agent fleet.Agent) (a fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "editAgent",
			"owner_id", a.MFOwnerID,
			"service_id", a.MFThingID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.EditAgent(ctx, token, agent)
}

func (m metricsMiddleware) ViewAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentGroupByIDInternal",
			"owner_id", ownerID,
			"service_id", groupID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ViewAgentGroupByIDInternal(ctx, groupID, ownerID)
}

func (m metricsMiddleware) ViewAgentGroupByID(ctx context.Context, token string, groupID string) (group fleet.AgentGroup, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentGroupByID",
			"owner_id", group.MFOwnerID,
			"service_id", groupID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ViewAgentGroupByID(ctx, token, groupID)
}

func (m metricsMiddleware) ListAgentGroups(ctx context.Context, token string, pm fleet.PageMetadata) (groups fleet.PageAgentGroup, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "listAgentGroups",
		}
		if len(groups.AgentGroups) != 0{
			labels = append(labels,
				"owner_id", groups.AgentGroups[0].MFOwnerID,
				"service_id", "")
		}else {
		labels = append(labels,
			"owner_id", "",
			"service_id", "")
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ListAgentGroups(ctx, token, pm)
}

func (m metricsMiddleware) EditAgentGroup(ctx context.Context, token string, ag fleet.AgentGroup) (fleet.AgentGroup, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "editAgentGroup",
			"owner_id", ag.MFOwnerID,
			"service_id", ag.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.EditAgentGroup(ctx, token, ag)
}

func (m metricsMiddleware) ListAgents(ctx context.Context, token string, pm fleet.PageMetadata) (agents fleet.Page, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "listAgents",
		}
		if agents.Agents != nil{
			labels = append(labels,
				"owner_id", agents.Agents[0].MFOwnerID,
				"service_id", "")
		} else {
			labels = append(labels,
				"owner_id", "",
				"service_id", "")
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ListAgents(ctx, token, pm)
}

func (m metricsMiddleware) CreateAgent(ctx context.Context, token string, a fleet.Agent) (agent fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "createAgent",
			"owner_id", agent.MFOwnerID,
			"service_id", agent.MFThingID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.CreateAgent(ctx, token, a)
}

func (m metricsMiddleware) CreateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (group fleet.AgentGroup, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "createAgentGroup",
			"owner_id", group.MFOwnerID,
			"service_id", group.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.CreateAgentGroup(ctx, token, s)
}

func (m metricsMiddleware) RemoveAgentGroup(ctx context.Context, token string, groupID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "removeAgentGroup",
			"owner_id", "",
			"service_id", groupID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.RemoveAgentGroup(ctx, token, groupID)
}

func (m metricsMiddleware) ValidateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (group fleet.AgentGroup, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "validateAgentGroup",
			"owner_id", group.MFOwnerID,
			"service_id", group.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ValidateAgentGroup(ctx, token, s)
}

func (m metricsMiddleware) ValidateAgent(ctx context.Context, token string, a fleet.Agent) (agent fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "validateAgent",
			"owner_id", agent.MFOwnerID,
			"service_id", agent.MFThingID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.ValidateAgent(ctx, token, a)
}

func (m metricsMiddleware) RemoveAgent(ctx context.Context, token string, thingID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "removeAgent",
			"owner_id", "",
			"service_id", thingID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(time.Since(begin).Seconds())

	}(time.Now())

	return m.svc.RemoveAgent(ctx, token, thingID)
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc fleet.Service, counter metrics.Counter, latency metrics.Histogram) fleet.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
