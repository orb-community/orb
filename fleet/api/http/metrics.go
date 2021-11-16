/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/ns1labs/orb/fleet"
	"github.com/prometheus/client_golang/prometheus/push"
	"time"
)

var _ fleet.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     fleet.Service
}

func (m metricsMiddleware) ViewOwnerByChannelIDInternal(ctx context.Context, channelID string) (fleet.Agent, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "viewOwnerByChannelIDInternal").Add(1)
		m.latency.With("method", "viewOwnerByChannelIDInternal").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "view_owner_by_channel_id_internal").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ViewOwnerByChannelIDInternal(ctx, channelID)
}

func (m metricsMiddleware) ViewAgentBackend(ctx context.Context, token string, name string) (interface{}, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "viewAgentBackend").Add(1)
		m.latency.With("method", "viewAgentBackend").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "view_agent_backend").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ViewAgentBackend(ctx, token, name)
}

func (m metricsMiddleware) ListAgentBackends(ctx context.Context, token string) ([]string, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "listAgentBackends").Add(1)
		m.latency.With("method", "listAgentBackends").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "list_agent_backends").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ListAgentBackends(ctx, token)
}

func (m metricsMiddleware) ViewAgentByIDInternal(ctx context.Context, ownerID string, thingID string) (fleet.Agent, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "viewAgentByIDInternal").Add(1)
		m.latency.With("method", "viewAgentByIDInternal").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "view_agent_by_id_internal").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ViewAgentByIDInternal(ctx, ownerID, thingID)
}

func (m metricsMiddleware) ViewAgentByID(ctx context.Context, token string, thingID string) (fleet.Agent, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "viewAgentByID").Add(1)
		m.latency.With("method", "viewAgentByID").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "view_agent_by_id").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ViewAgentByID(ctx, token, thingID)
}

func (m metricsMiddleware) EditAgent(ctx context.Context, token string, agent fleet.Agent) (fleet.Agent, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "editAgent").Add(1)
		m.latency.With("method", "editAgent").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "edit_agent").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.EditAgent(ctx, token, agent)
}

func (m metricsMiddleware) ViewAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "viewAgentGroupByIDInternal").Add(1)
		m.latency.With("method", "viewAgentGroupByIDInternal").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "view_agent_group_by_id_internal").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ViewAgentGroupByIDInternal(ctx, groupID, ownerID)
}

func (m metricsMiddleware) ViewAgentGroupByID(ctx context.Context, token string, groupID string) (fleet.AgentGroup, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "viewAgentGroupByID").Add(1)
		m.latency.With("method", "viewAgentGroupByID").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "view_agent_group_by_id").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ViewAgentGroupByID(ctx, token, groupID)
}

func (m metricsMiddleware) ListAgentGroups(ctx context.Context, token string, pm fleet.PageMetadata) (fleet.PageAgentGroup, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "listAgentGroups").Add(1)
		m.latency.With("method", "listAgentGroups").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "list_agent_group").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ListAgentGroups(ctx, token, pm)
}

func (m metricsMiddleware) EditAgentGroup(ctx context.Context, token string, ag fleet.AgentGroup) (fleet.AgentGroup, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "editAgentGroup").Add(1)
		m.latency.With("method", "editAgentGroup").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "edit_agent_group").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.EditAgentGroup(ctx, token, ag)
}

func (m metricsMiddleware) ListAgents(ctx context.Context, token string, pm fleet.PageMetadata) (fleet.Page, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "listAgents").Add(1)
		m.latency.With("method", "listAgents").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "list_agent").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ListAgents(ctx, token, pm)
}

func (m metricsMiddleware) CreateAgent(ctx context.Context, token string, a fleet.Agent) (fleet.Agent, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "createAgent").Add(1)
		m.latency.With("method", "createAgent").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "create_agent").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.CreateAgent(ctx, token, a)
}

func (m metricsMiddleware) CreateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (fleet.AgentGroup, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "createAgentGroup").Add(1)
		m.latency.With("method", "createAgentGroup").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "create_agent_group").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.CreateAgentGroup(ctx, token, s)
}

func (m metricsMiddleware) RemoveAgentGroup(ctx context.Context, token string, groupID string) error {
	defer func(begin time.Time) {
		m.counter.With("method", "removeAgentGroup").Add(1)
		m.latency.With("method", "removeAgentGroup").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "remove_agent_group").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.RemoveAgentGroup(ctx, token, groupID)
}

func (m metricsMiddleware) ValidateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (fleet.AgentGroup, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "validateAgentGroup").Add(1)
		m.latency.With("method", "validateAgentGroup").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "validate_agent_group").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ValidateAgentGroup(ctx, token, s)
}

func (m metricsMiddleware) ValidateAgent(ctx context.Context, token string, a fleet.Agent) (fleet.Agent, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "validateAgent").Add(1)
		m.latency.With("method", "validateAgent").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "validate_agent").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

	}(time.Now())

	return m.svc.ValidateAgent(ctx, token, a)
}

func (m metricsMiddleware) RemoveAgent(ctx context.Context, token string, thingID string) error {
	defer func(begin time.Time) {
		m.counter.With("method", "removeAgent").Add(1)
		m.latency.With("method", "removeAgent").Observe(time.Since(begin).Seconds())

		err := push.New("http://pushgateway:9091", "remove_agent").Push()
		if err != nil{
			fmt.Println("Can not push metrics:", err)
		}

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
