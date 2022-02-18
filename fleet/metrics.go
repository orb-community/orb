package fleet

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/mainflux"
)

var _ AgentCommsService = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     AgentCommsService
	auth    mainflux.AuthServiceClient
}

func (m metricsMiddleware) Start() error {
	panic("implement me")
}

func (m metricsMiddleware) Stop() error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyAgentNewGroupMembership(a Agent, ag AgentGroup) error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyAgentGroupMemberships(a Agent) error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyAgentAllDatasets(a Agent) error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyAgentStop(MFChannelID string, reason string) error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyGroupNewDataset(ctx context.Context, ag AgentGroup, datasetID string, policyID string, ownerID string) error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyGroupRemoval(ag AgentGroup) error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyGroupPolicyRemoval(ag AgentGroup, policyID string, policyName string, backend string) error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyGroupDatasetRemoval(ag AgentGroup, dsID string, policyID string) error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyGroupPolicyUpdate(ctx context.Context, ag AgentGroup, policyID string, ownerID string) error {
	panic("implement me")
}

func (m metricsMiddleware) NotifyAgentReset(channelID string, fullReset bool, reason string) error {
	panic("implement me")
}

func MetricsMiddleware(auth mainflux.AuthServiceClient, svc AgentCommsService, counter metrics.Counter, latency metrics.Histogram) AgentCommsService {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
		auth:    auth,
	}
}
