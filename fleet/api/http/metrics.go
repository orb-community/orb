/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/errors"
	"time"
)

var _ fleet.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     fleet.Service
	auth    mainflux.AuthServiceClient
}

func (m metricsMiddleware) ResetAgent(ct context.Context, token string, agentID string) error {
	ownerID, err := m.identify(token)
	if err != nil {
		return err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "resetAgent",
			"owner_id", ownerID,
			"agent_id", agentID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ResetAgent(ct, token, agentID)
}

func (m metricsMiddleware) ViewAgentInfoByChannelIDInternal(ctx context.Context, channelID string) (agent fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentInfoByChannelIDInternal",
			"owner_id", agent.MFOwnerID,
			"agent_id", agent.MFThingID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewAgentInfoByChannelIDInternal(ctx, channelID)
}

func (m metricsMiddleware) ViewAgentBackend(ctx context.Context, token string, name string) (interface{}, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return nil, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentBackend",
			"owner_id", ownerID,
			"agent_id", "",
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewAgentBackend(ctx, token, name)
}

func (m metricsMiddleware) ListAgentBackends(ctx context.Context, token string) ([]string, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return nil, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "listAgentBackends",
			"owner_id", ownerID,
			"agent_id", "",
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ListAgentBackends(ctx, token)
}

func (m metricsMiddleware) ViewAgentByIDInternal(ctx context.Context, ownerID string, thingID string) (fleet.Agent, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentByIDInternal",
			"owner_id", ownerID,
			"agent_id", thingID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewAgentByIDInternal(ctx, ownerID, thingID)
}

func (m metricsMiddleware) ViewAgentByID(ctx context.Context, token string, thingID string) (a fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentByID",
			"owner_id", a.MFOwnerID,
			"agent_id", thingID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewAgentByID(ctx, token, thingID)
}

func (m metricsMiddleware) ViewAgentMatchingGroupsByID(ctx context.Context, token string, thingID string) (a fleet.MatchingGroups, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentMatchingGroupsByID",
			"owner_id", a.OwnerID,
			"agent_id", thingID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewAgentMatchingGroupsByID(ctx, token, thingID)
}

func (m metricsMiddleware) EditAgent(ctx context.Context, token string, agent fleet.Agent) (a fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "editAgent",
			"owner_id", a.MFOwnerID,
			"agent_id", a.MFThingID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.EditAgent(ctx, token, agent)
}

func (m metricsMiddleware) ViewAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentGroupByIDInternal",
			"owner_id", ownerID,
			"agent_id", "",
			"group_id", groupID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewAgentGroupByIDInternal(ctx, groupID, ownerID)
}

func (m metricsMiddleware) ViewAgentGroupByID(ctx context.Context, token string, groupID string) (group fleet.AgentGroup, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewAgentGroupByID",
			"owner_id", group.MFOwnerID,
			"agent_id", "",
			"group_id", groupID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewAgentGroupByID(ctx, token, groupID)
}

func (m metricsMiddleware) ListAgentGroups(ctx context.Context, token string, pm fleet.PageMetadata) (groups fleet.PageAgentGroup, _ error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return fleet.PageAgentGroup{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "listAgentGroups",
			"owner_id", ownerID,
			"agent_id", "",
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ListAgentGroups(ctx, token, pm)
}

func (m metricsMiddleware) EditAgentGroup(ctx context.Context, token string, ag fleet.AgentGroup) (fleet.AgentGroup, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "editAgentGroup",
			"owner_id", ag.MFOwnerID,
			"agent_id", "",
			"group_id", ag.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.EditAgentGroup(ctx, token, ag)
}

func (m metricsMiddleware) ListAgents(ctx context.Context, token string, pm fleet.PageMetadata) (agents fleet.Page, _ error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return fleet.Page{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "listAgents",
			"owner_id", ownerID,
			"agent_id", "",
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ListAgents(ctx, token, pm)
}

func (m metricsMiddleware) CreateAgent(ctx context.Context, token string, a fleet.Agent) (agent fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "createAgent",
			"owner_id", agent.MFOwnerID,
			"agent_id", agent.MFThingID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.CreateAgent(ctx, token, a)
}

func (m metricsMiddleware) CreateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (group fleet.AgentGroup, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "createAgentGroup",
			"owner_id", group.MFOwnerID,
			"agent_id", "",
			"group_id", group.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.CreateAgentGroup(ctx, token, s)
}

func (m metricsMiddleware) RemoveAgentGroup(ctx context.Context, token string, groupID string) error {
	ownerID, err := m.identify(token)
	if err != nil {
		return err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "removeAgentGroup",
			"owner_id", ownerID,
			"agent_id", "",
			"group_id", groupID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.RemoveAgentGroup(ctx, token, groupID)
}

func (m metricsMiddleware) ValidateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (group fleet.AgentGroup, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "validateAgentGroup",
			"owner_id", group.MFOwnerID,
			"agent_id", "",
			"group_id", group.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ValidateAgentGroup(ctx, token, s)
}

func (m metricsMiddleware) ValidateAgent(ctx context.Context, token string, a fleet.Agent) (agent fleet.Agent, _ error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "validateAgent",
			"owner_id", agent.MFOwnerID,
			"agent_id", agent.MFThingID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ValidateAgent(ctx, token, a)
}

func (m metricsMiddleware) RemoveAgent(ctx context.Context, token string, thingID string) error {
	ownerID, err := m.identify(token)
	if err != nil {
		return err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "removeAgent",
			"owner_id", ownerID,
			"agent_id", thingID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.RemoveAgent(ctx, token, thingID)
}

func (m metricsMiddleware) GetPoliciesState(ctx context.Context, agent fleet.Agent) (map[string]interface{}, error) {

	defer func(begin time.Time) {
		labels := []string{
			"method", "getPoliciesState",
			"owner_id", agent.MFOwnerID,
			"agent_id", agent.MFThingID,
			"group_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.GetPoliciesState(ctx, agent)
}

func (m metricsMiddleware) identify(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return res.GetId(), nil
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(auth mainflux.AuthServiceClient, svc fleet.Service, counter metrics.Counter, latency metrics.Histogram) fleet.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
		auth:    auth,
	}
}
