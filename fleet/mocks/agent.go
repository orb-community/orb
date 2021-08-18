package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/types"
)

var _ fleet.AgentRepository = (*agentRepositoryMock)(nil)

type agentRepositoryMock struct {
}

func (a agentRepositoryMock) RetrieveMatchingAgents(ctx context.Context, ownerID string, tags types.Tags) (types.Metadata, error) {
	return nil, nil
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
