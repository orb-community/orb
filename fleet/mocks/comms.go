package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet"
)

var _ fleet.AgentCommsService = (*agentCommsServiceMock)(nil)

type agentCommsServiceMock struct{}

func (ac agentCommsServiceMock) NotifyPolicyRemoval(ag fleet.AgentGroup) error {
	return nil
}

func NewFleetCommService() fleet.AgentCommsService {
	return &agentCommsServiceMock{}
}

func (ac agentCommsServiceMock) Start() error {
	return nil
}

func (ac agentCommsServiceMock) Stop() error {
	return nil
}

func (ac agentCommsServiceMock) NotifyNewAgentGroupMembership(a fleet.Agent, ag fleet.AgentGroup) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyAgentGroupMembership(a fleet.Agent) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyAgentPolicies(a fleet.Agent) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupNewAgentPolicy(ctx context.Context, ag fleet.AgentGroup, policyID string, ownerID string) error {
	return nil
}

func (ac agentCommsServiceMock) InactivateDatasetByAgentGroup(groupID string, ownerID string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupRemoval(a fleet.AgentGroup) error {
	return nil
}
