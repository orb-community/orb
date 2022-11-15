package fleet

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

var _ AgentCommsService = (*commsMetricsMiddleware)(nil)

type commsMetricsMiddleware struct {
	requestCounter metrics.Counter
	requestLatency metrics.Histogram
	svc            AgentCommsService
}

func (c commsMetricsMiddleware) NotifyGroupDatasetEdit(ctx context.Context, ag AgentGroup, datasetID, policyID, ownerID string, valid bool) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupDatasetEdit",
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ag.MFOwnerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return c.svc.NotifyGroupDatasetEdit(ctx, ag, datasetID, policyID, ownerID, valid)
}

func (c commsMetricsMiddleware) Start() error {
	return c.svc.Start()
}

func (c commsMetricsMiddleware) Stop() error {
	return c.svc.Stop()
}

func (c commsMetricsMiddleware) NotifyAgentNewGroupMembership(ctx context.Context, a Agent, ag AgentGroup) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentNewGroupMembership",
			"agent_id", a.MFThingID,
			"agent_name", a.Name.String(),
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", a.MFOwnerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return c.svc.NotifyAgentNewGroupMembership(ctx, a, ag)
}

func (c commsMetricsMiddleware) NotifyAgentGroupMemberships(ctx context.Context, a Agent) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentGroupMemberships",
			"agent_id", a.MFThingID,
			"agent_name", a.Name.String(),
			"group_id", "",
			"group_name", "",
			"owner_id", a.MFOwnerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyAgentGroupMemberships(ctx, a)
}

func (c commsMetricsMiddleware) NotifyAgentAllDatasets(ctx context.Context, a Agent) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentAllDatasets",
			"agent_id", a.MFThingID,
			"agent_name", a.Name.String(),
			"group_id", "",
			"group_name", "",
			"owner_id", a.MFOwnerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyAgentAllDatasets(ctx, a)
}

func (c commsMetricsMiddleware) NotifyAgentStop(ctx context.Context, agent Agent, reason string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentStop",
			"agent_id", agent.MFThingID,
			"agent_name", agent.Name.String(),
			"group_id", "",
			"group_name", "",
			"owner_id", agent.MFOwnerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyAgentStop(ctx, agent, reason)
}

func (c commsMetricsMiddleware) NotifyGroupNewDataset(ctx context.Context, ag AgentGroup, datasetID string, policyID string, ownerID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupNewDataset",
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ownerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupNewDataset(ctx, ag, datasetID, policyID, ownerID)
}

func (c commsMetricsMiddleware) NotifyGroupRemoval(ctx context.Context, ag AgentGroup) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupRemoval",
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ag.MFOwnerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupRemoval(ctx, ag)
}

func (c commsMetricsMiddleware) NotifyGroupPolicyRemoval(ctx context.Context, ag AgentGroup, policyID string, policyName string, backend string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupPolicyRemoval",
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ag.MFOwnerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupPolicyRemoval(ctx, ag, policyID, policyName, backend)
}

func (c commsMetricsMiddleware) NotifyGroupDatasetRemoval(ctx context.Context, ag AgentGroup, dsID string, policyID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupDatasetRemoval",
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ag.MFOwnerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupDatasetRemoval(ctx, ag, dsID, policyID)
}

func (c commsMetricsMiddleware) NotifyGroupPolicyUpdate(ctx context.Context, ag AgentGroup, policyID string, ownerID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupPolicyUpdate",
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ownerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupPolicyUpdate(ctx, ag, policyID, ownerID)
}

func (c commsMetricsMiddleware) NotifyAgentReset(ctx context.Context, agent Agent, fullReset bool, reason string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentReset",
			"agent_id", agent.MFThingID,
			"agent_name", agent.Name.String(),
			"group_id", "",
			"group_name", "",
			"owner_id", agent.MFOwnerID,
		}

		c.requestCounter.With(labels...).Add(1)
		c.requestLatency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyAgentReset(ctx, agent, fullReset, reason)
}

func CommsMetricsMiddleware(svc AgentCommsService, counter metrics.Counter, latency metrics.Histogram) AgentCommsService {
	return &commsMetricsMiddleware{
		requestCounter: counter,
		requestLatency: latency,
		svc:            svc,
	}
}
