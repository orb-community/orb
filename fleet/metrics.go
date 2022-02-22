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
	gauge   metrics.Gauge
	svc     AgentCommsService
}

func (c commsMetricsMiddleware) Start() error {
	return c.svc.Start()
}

func (c commsMetricsMiddleware) Stop() error {
	return c.svc.Stop()
}

func (c commsMetricsMiddleware) NotifyAgentNewGroupMembership(a Agent, ag AgentGroup) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentNewGroupMembership",
			"agent_id", a.MFThingID,
			"agent_name", a.Name.String(),
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", a.MFOwnerID,
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
			"agent_name", a.Name.String(),
			"group_id", "",
			"group_name", "",
			"owner_id", a.MFOwnerID,
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
			"agent_name", a.Name.String(),
			"group_id", "",
			"group_name", "",
			"owner_id", a.MFOwnerID,
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
			"agent_id", "",
			"agent_name", "",
			"group_id", "",
			"group_name", "",
			"owner_id", "",
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
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ownerID,
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
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ag.MFOwnerID,
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
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ag.MFOwnerID,
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
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ag.MFOwnerID,
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
			"agent_id", "",
			"agent_name", "",
			"group_id", ag.ID,
			"group_name", ag.Name.String(),
			"owner_id", ownerID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyGroupPolicyUpdate(ctx, ag, policyID, ownerID)
}

func (c commsMetricsMiddleware) NotifyAgentReset(channelID string, fullReset bool, reason string) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "NotifyAgentReset",
			"agent_id", "",
			"agent_name", "",
			"group_id", "",
			"group_name", "",
			"owner_id", "",
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())
	return c.svc.NotifyAgentReset(channelID, fullReset, reason)
}

func CommsMetricsMiddleware(svc AgentCommsService, counter metrics.Counter, latency metrics.Histogram, gauge metrics.Gauge) AgentCommsService {
	return &commsMetricsMiddleware{
		counter: counter,
		latency: latency,
		gauge:   gauge,
		svc:     svc,
	}
}
