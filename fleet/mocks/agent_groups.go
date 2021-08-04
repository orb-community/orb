// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package mocks

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/errors"
)

var _ fleet.AgentGroupRepository = (*agentGroupRepositoryMock)(nil)

type agentGroupRepositoryMock struct {
	counter        int64
	agentGroupMock map[string]fleet.AgentGroup
}

func NewAgentGroupRepository() agentGroupRepositoryMock {
	return agentGroupRepositoryMock{
		agentGroupMock: make(map[string]fleet.AgentGroup),
	}
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
	if c, ok := a.agentGroupMock[groupID]; ok {
		return c, nil
	}

	return fleet.AgentGroup{}, fleet.ErrNotFound
}
