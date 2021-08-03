package mocks

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/errors"
)

var _ fleet.AgentGroupRepository = (*agentGroupRepositoryMock)(nil)

type agentGroupRepositoryMock struct {
	counter int64
	agentGroupMock map[string]fleet.AgentGroup
}

func NewAgentGroupRepository() agentGroupRepositoryMock {
	return agentGroupRepositoryMock{}
}

func (a agentGroupRepositoryMock) Save(ctx context.Context, group fleet.AgentGroup) (string, error) {
	ID, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(errors.ErrMalformedEntity, err)
	}
	group.ID = ID.String()
	a.agentGroupMock[ID.String()] = group
	return ID.String(), nil
}

func (a agentGroupRepositoryMock) RetrieveAllByAgent(ctx context.Context, agent fleet.Agent) ([]fleet.AgentGroup, error) {
	panic("implement me")
}

func (a agentGroupRepositoryMock) RetrieveByID(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	panic("implement me")
}


