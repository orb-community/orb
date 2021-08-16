package mocks

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/errors"
)

var _ fleet.AgentRepository = (*agentRepositoryMock)(nil)

type agentRepositoryMock struct {
	counter			uint64
	agentMock map[string]fleet.Agent
}

func (a agentRepositoryMock) UpdateHeartbeatByIDWithChannel(ctx context.Context, agent fleet.Agent) error {
	panic("implement me")
}

func (a agentRepositoryMock) Save(ctx context.Context, agent fleet.Agent) error {
	ID, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}
	a.counter++
	agent.MFThingID = ID.String()
	a.agentMock[ID.String()] = agent
	return nil
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
