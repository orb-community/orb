package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet"
	"reflect"
)

var _ fleet.AgentCommsService = (*agentCommsServiceMock)(nil)

type agentCommsServiceMock struct {
	aGroupRepoMock fleet.AgentGroupRepository
	aRepoMock      fleet.AgentRepository
	commsMock      map[string][]fleet.Agent
}

func (ac agentCommsServiceMock) NotifyGroupDatasetEdit(ctx context.Context, ag fleet.AgentGroup, datasetID, policyID, ownerID string, valid bool) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyAgentReset(agent fleet.Agent, fullReset bool, reason string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyAgentStop(agent fleet.Agent, reason string) error {
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

func NewFleetCommService(agentRepo fleet.AgentRepository, agentGroupRepo fleet.AgentGroupRepository) fleet.AgentCommsService {
	return &agentCommsServiceMock{
		aRepoMock:      agentRepo,
		aGroupRepoMock: agentGroupRepo,
		commsMock:      make(map[string][]fleet.Agent),
	}
}

func (ac agentCommsServiceMock) Start() error {
	return nil
}

func (ac agentCommsServiceMock) Stop() error {
	return nil
}

func (ac agentCommsServiceMock) NotifyAgentNewGroupMembership(a fleet.Agent, ag fleet.AgentGroup) error {
	aGroups, err := ac.aGroupRepoMock.RetrieveAllAgentGroupsByOwner(context.Background(), ag.MFOwnerID, fleet.PageMetadata{Limit: 1})
	if err != nil {
		return err
	}

	for _, group := range aGroups.AgentGroups {
		if reflect.DeepEqual(group.Tags, a.AgentTags) {
			ac.commsMock[group.ID] = append(ac.commsMock[group.ID], a)
		}
	}

	return nil
}

func (ac agentCommsServiceMock) NotifyAgentGroupMemberships(a fleet.Agent) error {
	list, err := ac.aGroupRepoMock.RetrieveAllByAgent(context.Background(), a)
	if err != nil {
		return err
	}

	for _, agentGroup := range list {
		agentList, _ := ac.commsMock[agentGroup.ID]
		for i, agent := range agentList {
			if reflect.DeepEqual(agent.AgentTags, a.AgentTags) {
				agentList[i].Name = a.Name
			} else {
				agentList[i] = agentList[len(agentList)-1]
				agentList[len(agentList)-1] = fleet.Agent{}
				agentList = agentList[:len(agentList)-1]
			}
		}
	}
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
