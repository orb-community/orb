package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

var _ fleet.AgentRepository = (*agentRepositoryMock)(nil)

type agentRepositoryMock struct {
	counter    uint64
	agentsMock map[string]fleet.Agent
}

func (a agentRepositoryMock) SetStaleStatus(_ context.Context, _ time.Duration) (int64, error) {
	return 0, nil
}

func (a agentRepositoryMock) RetrieveAgentInfoByChannelID(_ context.Context, channelID string) (fleet.Agent, error) {
	for _, ag := range a.agentsMock {
		if ag.MFChannelID == channelID {
			return ag, nil
		}
	}
	return fleet.Agent{}, fleet.ErrNotFound
}

func (a agentRepositoryMock) RetrieveAgentMetadataByOwner(_ context.Context, _ string) ([]types.Metadata, error) {
	var taps []types.Metadata
	return taps, nil
}

func (a agentRepositoryMock) RetrieveByID(_ context.Context, ownerID string, thingID string) (fleet.Agent, error) {
	if _, ok := a.agentsMock[thingID]; ok {
		if a.agentsMock[thingID].MFOwnerID != ownerID {
			return fleet.Agent{}, fleet.ErrNotFound
		}
		return a.agentsMock[thingID], nil
	}
	return fleet.Agent{}, fleet.ErrNotFound
}

func (a agentRepositoryMock) RetrieveMatchingAgents(_ context.Context, _ string, _ types.Tags) (types.Metadata, error) {
	return nil, nil
}

func (a agentRepositoryMock) UpdateHeartbeatByIDWithChannel(_ context.Context, _ fleet.Agent) error {
	panic("implement me")
}

func (a agentRepositoryMock) Save(_ context.Context, agent fleet.Agent) error {
	for _, ag := range a.agentsMock {
		if ag.Name == agent.Name && ag.MFOwnerID == agent.MFOwnerID {
			return fleet.ErrConflict
		}
	}
	a.agentsMock[agent.MFThingID] = agent
	a.counter++
	return nil
}

func (a agentRepositoryMock) UpdateAgentByID(_ context.Context, ownerID string, agent fleet.Agent) error {
	if _, ok := a.agentsMock[agent.MFThingID]; ok {
		if a.agentsMock[agent.MFThingID].MFOwnerID != ownerID {
			return fleet.ErrUpdateEntity
		}

		currentGroup := a.agentsMock[agent.MFThingID]
		currentGroup.Name = agent.Name
		currentGroup.OrbTags = agent.OrbTags

		a.agentsMock[agent.MFThingID] = currentGroup
		return nil
	}
	return fleet.ErrNotFound
}

func (a agentRepositoryMock) UpdateDataByIDWithChannel(_ context.Context, _ fleet.Agent) error {
	panic("implement me")
}

func (a agentRepositoryMock) RetrieveByIDWithChannel(_ context.Context, thingID string, channelID string) (fleet.Agent, error) {
	if _, ok := a.agentsMock[thingID]; ok {
		if a.agentsMock[thingID].MFChannelID != channelID {
			return fleet.Agent{}, fleet.ErrNotFound
		}
		return a.agentsMock[thingID], nil
	}
	return fleet.Agent{}, fleet.ErrNotFound
}

func (a agentRepositoryMock) RetrieveAll(_ context.Context, owner string, pm fleet.PageMetadata) (fleet.Page, error) {
	first := pm.Offset
	last := first + pm.Limit

	var agents []fleet.Agent
	id := uint64(0)
	for _, v := range a.agentsMock {
		if v.MFOwnerID == owner && id >= first && id < last {
			agents = append(agents, v)
		}
		id++
	}

	agents = sortAgents(pm, agents)

	pageAgentGroup := fleet.Page{
		PageMetadata: fleet.PageMetadata{
			Total: a.counter,
		},
		Agents: agents,
	}
	return pageAgentGroup, nil
}

func (a agentRepositoryMock) RetrieveAllByAgentGroupID(_ context.Context, owner string, agentGroupID string, _ bool) ([]fleet.Agent, error) {
	if agentGroupID == "" || owner == "" {
		return nil, errors.ErrMalformedEntity
	}

	var agents []fleet.Agent
	id := uint64(0)
	for _, v := range a.agentsMock {
		if v.MFOwnerID == owner {
			agents = append(agents, v)
		}
		id++
	}

	return agents, nil
}

func (a agentRepositoryMock) Delete(_ context.Context, ownerID, thingID string) error {
	if _, ok := a.agentsMock[thingID]; ok {
		if a.agentsMock[thingID].MFOwnerID == ownerID {
			delete(a.agentsMock, thingID)
		}
	}

	return nil
}

func NewAgentRepositoryMock() fleet.AgentRepository {
	return &agentRepositoryMock{
		agentsMock: make(map[string]fleet.Agent),
	}
}
