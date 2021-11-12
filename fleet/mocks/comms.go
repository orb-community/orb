package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet"
)

var _ fleet.AgentCommsService = (*agentCommsServiceMock)(nil)

type agentCommsServiceMock struct{}

func (ac agentCommsServiceMock) NotifyAgentStop(MFChannelID string, reason string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupPolicyRemoval(ag fleet.AgentGroup, policyID string, policyName string, backend string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupPolicyUpdate(ctx context.Context, ag fleet.AgentGroup, policyID string, ownerID string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupDatasetRemoval(ag fleet.AgentGroup, dsID string, policyID string) error {
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

func (ac agentCommsServiceMock) NotifyAgentNewGroupMembership(a fleet.Agent, ag fleet.AgentGroup) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyAgentGroupMemberships(a fleet.Agent) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyAgentAllDatasets(a fleet.Agent) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupNewDataset(ctx context.Context, ag fleet.AgentGroup, datasetID string, policyID string, ownerID string) error {
	return nil
}

func (ac agentCommsServiceMock) InactivateDatasetByAgentGroup(groupID string, ownerID string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupRemoval(a fleet.AgentGroup) error {
	return nil
}
