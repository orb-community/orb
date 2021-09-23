package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet"
)

var _ fleet.AgentCommsService = (*agentCommsServiceMock)(nil)

type agentCommsServiceMock struct{}

func (ac agentCommsServiceMock) NofityDatasetRemoval(ag fleet.AgentGroup) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyPolicyRemoval(policyID string, ag fleet.AgentGroup) error {
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
