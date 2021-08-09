package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet"
)

var _ fleet.AgentRepository = (*agentRepositoryMock)(nil)

type agentRepositoryMock struct {
}

func (a agentRepositoryMock) UpdateHeartbeatByIDWithChannel(ctx context.Context, agent fleet.Agent) error {
	panic("implement me")
}

func (a agentRepositoryMock) Save(ctx context.Context, agent fleet.Agent) error {
	panic("implement me")
}

func (a agentRepositoryMock) UpdateDataByIDWithChannel(ctx context.Context, agent fleet.Agent) error {
	panic("implement me")
}

func (a agentRepositoryMock) RetrieveByIDWithChannel(ctx context.Context, thingID string, channelID string) (fleet.Agent, error) {
	panic("implement me")
}

func (a agentRepositoryMock) RetrieveAll(ctx context.Context, owner string, pm fleet.PageMetadata) (fleet.Page, error) {
	panic("implement me")
}

func (a agentRepositoryMock) RetrieveAllByAgentGroupID(ctx context.Context, owner string, agentGroupID string, onlinishOnly bool) ([]fleet.Agent, error) {
	return []fleet.Agent{}, nil
}

func NewAgentRepositoryMock() fleet.AgentRepository {
	return agentRepositoryMock{}
}
