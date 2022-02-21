package fleet

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

var _ AgentCommsService = (*commsMetricsMiddleware)(nil)

type commsMetricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     AgentCommsService
}

func (c commsMetricsMiddleware) Start() error {
	panic("implement me")
}

func (c commsMetricsMiddleware) Stop() error {
	panic("implement me")
}

func (c commsMetricsMiddleware) NotifyAgentNewGroupMembership(a Agent, ag AgentGroup) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentNewGroupMembership",
			"agent_id", a.MFThingID,
			"group_id", ag.ID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return c.svc.NotifyAgentNewGroupMembership(a, ag)
}

func (c commsMetricsMiddleware) NotifyAgentGroupMemberships(a Agent) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentGroupMemberships",
			"agent_id", a.MFThingID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyAgentGroupMemberships(a)
}

func (c commsMetricsMiddleware) NotifyAgentAllDatasets(a Agent) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentAllDatasets",
			"agent_id", a.MFThingID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyAgentAllDatasets(a)
}

func (c commsMetricsMiddleware) NotifyAgentStop(MFChannelID string, reason string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentStop",
			"channel_id", MFChannelID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyAgentStop(MFChannelID, reason)
}

func (c commsMetricsMiddleware) NotifyGroupNewDataset(ctx context.Context, ag AgentGroup, datasetID string, policyID string, ownerID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupNewDataset",
			"group_id", ag.ID,
			"policy_id", policyID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupNewDataset(ctx, ag, datasetID, policyID, ownerID)
}

func (c commsMetricsMiddleware) NotifyGroupRemoval(ag AgentGroup) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupRemoval",
			"group_id", ag.ID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupRemoval(ag)
}

func (c commsMetricsMiddleware) NotifyGroupPolicyRemoval(ag AgentGroup, policyID string, policyName string, backend string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupPolicyRemoval",
			"group_id", ag.ID,
			"policy_id", policyID,
			"policy_name", policyName,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupPolicyRemoval(ag, policyID, policyName, backend)
}

func (c commsMetricsMiddleware) NotifyGroupDatasetRemoval(ag AgentGroup, dsID string, policyID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupDatasetRemoval",
			"group_id", ag.ID,
			"policy_id", policyID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupDatasetRemoval(ag, dsID, policyID)
}

func (c commsMetricsMiddleware) NotifyGroupPolicyUpdate(ctx context.Context, ag AgentGroup, policyID string, ownerID string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupPolicyUpdate",
			"group_id", ag.ID,
			"policy_id", policyID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupPolicyUpdate(ctx, ag, policyID, ownerID)
}

func (c commsMetricsMiddleware) NotifyAgentReset(channelID string, fullReset bool, reason string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyGroupPolicyUpdate",
			"channel_id", channelID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyAgentReset(channelID, fullReset, reason)
}

func CommsMetricsMiddleware(svc AgentCommsService, counter metrics.Counter, latency metrics.Histogram) AgentCommsService {
	return &commsMetricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
