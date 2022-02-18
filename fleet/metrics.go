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
			"agent_id", a.MFChannelID,
			"group_id", ag.ID,
		}

		c.counter.With(labels...).Add(1)
		c.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return c.svc.NotifyAgentNewGroupMembership(a, ag)
}

func (c commsMetricsMiddleware) NotifyAgentGroupMemberships(a Agent) error {
	panic("implement me")
}

func (c commsMetricsMiddleware) NotifyAgentAllDatasets(a Agent) error {
	panic("implement me")
}

func (c commsMetricsMiddleware) NotifyAgentStop(MFChannelID string, reason string) error {
	panic("implement me")
}

func (c commsMetricsMiddleware) NotifyGroupNewDataset(ctx context.Context, ag AgentGroup, datasetID string, policyID string, ownerID string) error {
	panic("implement me")
}

func (c commsMetricsMiddleware) NotifyGroupRemoval(ag AgentGroup) error {
	panic("implement me")
}

func (c commsMetricsMiddleware) NotifyGroupPolicyRemoval(ag AgentGroup, policyID string, policyName string, backend string) error {
	panic("implement me")
}

func (c commsMetricsMiddleware) NotifyGroupDatasetRemoval(ag AgentGroup, dsID string, policyID string) error {
	panic("implement me")
}

func (c commsMetricsMiddleware) NotifyGroupPolicyUpdate(ctx context.Context, ag AgentGroup, policyID string, ownerID string) error {
	panic("implement me")
}

func (c commsMetricsMiddleware) NotifyAgentReset(channelID string, fullReset bool, reason string) error {
	panic("implement me")
}

func CommsMetricsMiddleware(svc AgentCommsService, counter metrics.Counter, latency metrics.Histogram) AgentCommsService {
	return &commsMetricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
