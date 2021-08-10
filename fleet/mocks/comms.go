package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet"
)

var _ fleet.AgentCommsService = (*agentCommsServiceMock)(nil)

type agentCommsServiceMock struct{}

func NewFleetCommService() fleet.AgentCommsService {
	return &agentCommsServiceMock{}
}

func (ac agentCommsServiceMock) Start() error {
	panic("implement me")
}

func (ac agentCommsServiceMock) Stop() error {
	panic("implement me")
}

func (ac agentCommsServiceMock) NotifyNewAgentGroupMembership(a fleet.Agent, ag fleet.AgentGroup) error {
	panic("implement me")
}

func (ac agentCommsServiceMock) NotifyAgentGroupMembership(a fleet.Agent) error {
	panic("implement me")
}

func (ac agentCommsServiceMock) NotifyAgentPolicies(a fleet.Agent) error {
	panic("implement me")
}

func (ac agentCommsServiceMock) NotifyGroupNewAgentPolicy(ctx context.Context, ag fleet.AgentGroup, policyID string, ownerID string) error {
	panic("implement me")
}

func (ac agentCommsServiceMock) InactivateDatasetByAgentGroup(groupID string, ownerID string) error {
	return nil
}

func (ac agentCommsServiceMock) UnsubscribeAgentGroupMembership(a fleet.Agent) error {
	return nil
}
