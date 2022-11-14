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

func (ac agentCommsServiceMock) NotifyAgentConfig(ctx context.Context, a fleet.Agent) error {
	//TODO implement me
	panic("implement me")
}

func (ac agentCommsServiceMock) NotifyGroupDatasetEdit(_ context.Context, _ fleet.AgentGroup, _, _, _ string, _ bool) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyAgentReset(_ context.Context, _ fleet.Agent, _ bool, _ string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyAgentStop(_ context.Context, _ fleet.Agent, _ string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupPolicyRemoval(_ context.Context, _ fleet.AgentGroup, _ string, _ string, _ string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupPolicyUpdate(_ context.Context, _ fleet.AgentGroup, _ string, _ string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupDatasetRemoval(_ context.Context, _ fleet.AgentGroup, _ string, _ string) error {
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

func (ac agentCommsServiceMock) NotifyAgentNewGroupMembership(ctx context.Context, a fleet.Agent, ag fleet.AgentGroup) error {
	aGroups, err := ac.aGroupRepoMock.RetrieveAllAgentGroupsByOwner(ctx, ag.MFOwnerID, fleet.PageMetadata{Limit: 1})
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

func (ac agentCommsServiceMock) NotifyAgentGroupMemberships(ctx context.Context, a fleet.Agent) error {
	list, err := ac.aGroupRepoMock.RetrieveAllByAgent(ctx, a)
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

func (ac agentCommsServiceMock) NotifyAgentAllDatasets(_ context.Context, _ fleet.Agent) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupNewDataset(_ context.Context, _ fleet.AgentGroup, _ string, _ string, _ string) error {
	return nil
}

func (ac agentCommsServiceMock) InactivateDatasetByAgentGroup(_ context.Context, _, _ string) error {
	return nil
}

func (ac agentCommsServiceMock) NotifyGroupRemoval(_ context.Context, _ fleet.AgentGroup) error {
	return nil
}
