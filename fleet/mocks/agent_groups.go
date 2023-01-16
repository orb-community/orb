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
	"reflect"
)

var _ fleet.AgentGroupRepository = (*agentGroupRepositoryMock)(nil)

type agentGroupRepositoryMock struct {
	counter        uint64
	agentGroupMock map[string]fleet.AgentGroup
}

func NewAgentGroupRepository() fleet.AgentGroupRepository {
	return &agentGroupRepositoryMock{
		agentGroupMock: make(map[string]fleet.AgentGroup),
	}
}

func (a *agentGroupRepositoryMock) Save(ctx context.Context, group fleet.AgentGroup) (string, error) {
	ID, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(errors.ErrMalformedEntity, err)
	}

	for _, ag := range a.agentGroupMock {
		if ag.Name == group.Name {
			return "", errors.Wrap(errors.ErrConflict, err)
		}
	}

	a.counter++
	group.ID = ID.String()
	a.agentGroupMock[ID.String()] = group
	return ID.String(), nil
}

func (a *agentGroupRepositoryMock) RetrieveAllByAgent(ctx context.Context, agent fleet.Agent) ([]fleet.AgentGroup, error) {
	var agentGroups []fleet.AgentGroup

	if agent.MFThingID == "" {
		return agentGroups, errors.ErrMalformedEntity
	}

	for _, v := range a.agentGroupMock {
		if v.MFOwnerID == agent.MFOwnerID && (reflect.DeepEqual(v.Tags, agent.AgentTags) || reflect.DeepEqual(v.Tags, agent.OrbTags)) {
			agentGroups = append(agentGroups, v)
		}
	}
	return agentGroups, nil
}

func (a *agentGroupRepositoryMock) RetrieveByID(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	if c, ok := a.agentGroupMock[groupID]; ok {
		return c, nil
	}

	return fleet.AgentGroup{}, fleet.ErrNotFound
}

func (a *agentGroupRepositoryMock) RetrieveAllAgentGroupsByOwner(ctx context.Context, ownerID string, pm fleet.PageMetadata) (fleet.PageAgentGroup, error) {
	first := uint64(pm.Offset)
	last := first + uint64(pm.Limit)

	var agentGroups []fleet.AgentGroup
	id := uint64(0)
	for _, v := range a.agentGroupMock {
		if v.MFOwnerID == ownerID && id >= first && id < last {
			agentGroups = append(agentGroups, v)
		}
		id++
	}

	agentGroups = sortAgentGroups(pm, agentGroups)

	pageAgentGroup := fleet.PageAgentGroup{
		PageMetadata: fleet.PageMetadata{
			Total: a.counter,
		},
		AgentGroups: agentGroups,
	}
	return pageAgentGroup, nil
}

func (a *agentGroupRepositoryMock) Update(ctx context.Context, ownerID string, group fleet.AgentGroup) (fleet.AgentGroup, error) {
	if _, ok := a.agentGroupMock[group.ID]; ok {
		if a.agentGroupMock[group.ID].MFOwnerID != ownerID {
			return fleet.AgentGroup{}, fleet.ErrUpdateEntity
		}
		currentGroup := a.agentGroupMock[group.ID]
		currentGroup.Name = group.Name
		currentGroup.Description = group.Description
		currentGroup.Tags = group.Tags

		a.agentGroupMock[group.ID] = currentGroup

		return a.agentGroupMock[group.ID], nil
	}
	return fleet.AgentGroup{}, fleet.ErrNotFound
}

func (a *agentGroupRepositoryMock) Delete(ctx context.Context, groupID string, ownerID string) error {
	if _, ok := a.agentGroupMock[groupID]; ok {
		if a.agentGroupMock[groupID].MFOwnerID != ownerID {
			delete(a.agentGroupMock, groupID)
		}
	}
	return nil
}

func (a *agentGroupRepositoryMock) RetrieveMatchingGroups(ctx context.Context, ownerID string, thingID string) (fleet.MatchingGroups, error) {
	//TODO find a way to correlation Agents and Agents Groups on Mock
	var groups []fleet.Group
	for _, group := range a.agentGroupMock {
		if group.MFOwnerID == ownerID {
			groups = append(groups, fleet.Group{
				GroupID:   group.ID,
				GroupName: group.Name,
			})
		}
	}
	return fleet.MatchingGroups{OwnerID: ownerID, Groups: groups}, nil
}
