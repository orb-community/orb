// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/postgres"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"testing"
	"time"
)

const maxNameSize = 1024

var (
	invalidName = strings.Repeat("m", maxNameSize+1)
	logger, _   = zap.NewDevelopment()
	wrongValue  = "wrong-value"
)

func TestAgentSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	conflictThingID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	conflictNameID, err := types.NewIdentifier("myagent-conflict")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	// Conflict scenario
	agentCopy := agent
	agentCopy.MFThingID = conflictThingID.String()

	agentCopy.Name = conflictNameID

	err = agentRepo.Save(context.Background(), agentCopy)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	cases := map[string]struct {
		agent fleet.Agent
		err   error
	}{
		"create new agent": {
			agent: agent,
			err:   nil,
		},
		"create agent that already exist": {
			agent: agentCopy,
			err:   errors.ErrConflict,
		},
		"create new agent with empty OwnerID": {
			agent: fleet.Agent{
				Name:        nameID,
				MFOwnerID:   "",
				MFThingID:   thID.String(),
				MFChannelID: chID.String(),
			},
			err: errors.ErrMalformedEntity,
		},
		"create new agent with empty ThingID": {
			agent: fleet.Agent{
				Name:        nameID,
				MFOwnerID:   oID.String(),
				MFThingID:   "",
				MFChannelID: chID.String(),
			},
			err: errors.ErrMalformedEntity,
		},
		"create new agent with empty channelID": {
			agent: fleet.Agent{
				Name:        nameID,
				MFOwnerID:   oID.String(),
				MFThingID:   thID.String(),
				MFChannelID: "",
			},
			err: errors.ErrMalformedEntity,
		},
		"create new agent with empty invalid OwnerID": {
			agent: fleet.Agent{
				Name:        nameID,
				MFOwnerID:   "123",
				MFThingID:   thID.String(),
				MFChannelID: chID.String(),
			},
			err: errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := agentRepo.Save(context.Background(), tc.agent)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
		})
	}
}

func TestAgentRetrieve(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		thingID   string
		channelID string
		err       error
		tags      types.Tags
	}{
		"retrieve existing agent by thingID and channelID": {
			thingID:   agent.MFThingID,
			channelID: agent.MFChannelID,
			tags:      types.Tags{"testkey": "testvalue"},
			err:       nil,
		},
		"retrieve non-existent agent by thingID and channelID": {
			thingID:   agent.MFOwnerID,
			channelID: agent.MFChannelID,
			err:       errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			ag, err := agentRepo.RetrieveByIDWithChannel(context.Background(), tc.thingID, tc.channelID)
			if err == nil {
				assert.Equal(t, nameID, ag.Name, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
			}
			if len(tc.tags) > 0 {
				assert.Equal(t, tc.tags, ag.OrbTags)
				assert.Equal(t, tc.tags, ag.AgentTags)
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func TestAgentUpdateData(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		agent fleet.Agent
		err   error
	}{
		"update existing agent data by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:     thID.String(),
				MFChannelID:   chID.String(),
				AgentMetadata: types.Metadata{"newkey": "newvalue"},
			},
			err: nil,
		},
		"update non-existent agent data by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:     chID.String(),
				MFChannelID:   thID.String(),
				AgentMetadata: types.Metadata{"newkey": "newvalue"},
			},
			err: errors.ErrNotFound,
		},
		"update agent data by thingID and channelID with invalid thingID": {
			agent: fleet.Agent{
				MFThingID:     "123",
				MFChannelID:   chID.String(),
				AgentMetadata: types.Metadata{"newkey": "newvalue"},
			},
			err: errors.ErrMalformedEntity,
		},
		"update agent data by thingID with channelID with empty fields": {
			agent: fleet.Agent{
				MFThingID:     "",
				MFChannelID:   "",
				AgentMetadata: types.Metadata{"newkey": "newvalue"},
			},
			err: errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err = agentRepo.UpdateDataByIDWithChannel(context.Background(), tc.agent)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
			if err == nil {
				ag, err := agentRepo.RetrieveByIDWithChannel(context.Background(), tc.agent.MFThingID, tc.agent.MFChannelID)
				assert.Nil(t, err)
				assert.Equal(t, tc.agent.AgentMetadata, ag.AgentMetadata, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
			}
		})
	}
}

func TestAgentUpdateHeartbeat(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:        nameID,
		MFThingID:   thID.String(),
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		LastHBData:  types.Metadata{"heartbeatdata": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		agent fleet.Agent
		err   error
	}{
		"update existing agent heartbeat by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:   thID.String(),
				MFChannelID: chID.String(),
				LastHBData:  types.Metadata{"heartbeatdata2": "newvalue"},
			},
			err: nil,
		},
		"update non-existent agent heart beat by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:   chID.String(),
				MFChannelID: thID.String(),
				LastHBData:  types.Metadata{"heartbeatdata2": "newvalue"},
			},
			err: errors.ErrNotFound,
		},
		"update existing agent heartbeat with empty thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:   "",
				MFChannelID: "",
				LastHBData:  types.Metadata{"heartbeatdata2": "newvalue"},
			},
			err: errors.ErrMalformedEntity,
		},
		"update existing agent heartbeat with invalid thingID": {
			agent: fleet.Agent{
				MFThingID:   "123",
				MFChannelID: chID.String(),
				LastHBData:  types.Metadata{"heartbeatdata2": "newvalue"},
			},
			err: errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err = agentRepo.UpdateHeartbeatByIDWithChannel(context.Background(), tc.agent)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
			if err == nil {
				ag, err := agentRepo.RetrieveByIDWithChannel(context.Background(), tc.agent.MFThingID, tc.agent.MFChannelID)
				assert.Nil(t, err)
				assert.Equal(t, tc.agent.LastHBData, ag.LastHBData, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
			}
		})
	}
}

func TestMultiAgentRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	wrongoID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	name := "agent_name"
	metaStr := `{"field1":"value1","field2":{"subfield11":"value2","subfield12":{"subfield121":"value3","subfield122":"value4"}}}`
	subMetaStr := `{"field2":{"subfield12":{"subfield121":"value3"}}}`
	tagsStr := `{"node_type": "dns"}`
	subTagsStr := `{"region": "EU"}`
	mixTagsStr := `{"node_type": "dns", "region": "EU"}`

	metadata := types.Metadata{}
	json.Unmarshal([]byte(metaStr), &metadata)

	subMeta := types.Metadata{}
	json.Unmarshal([]byte(subMetaStr), &subMeta)

	tags := types.Tags{}
	json.Unmarshal([]byte(tagsStr), &tags)

	subTags := types.Tags{}
	json.Unmarshal([]byte(subTagsStr), &subTags)

	mixTags := types.Tags{}
	json.Unmarshal([]byte(mixTagsStr), &mixTags)

	wrongMeta := types.Metadata{
		"field": "value1",
	}
	wrongTags := types.Tags{
		"field": "value1",
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		thID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		chID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		th := fleet.Agent{
			MFOwnerID:   oID.String(),
			MFThingID:   thID.String(),
			MFChannelID: chID.String(),
		}

		th.Name, err = types.NewIdentifier(fmt.Sprintf("%s-%d", name, i))
		require.True(t, th.Name.IsValid(), "invalid Identifier name: %s")
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		th.AgentMetadata = metadata
		th.AgentTags = tags
		th.OrbTags = subTags

		err = agentRepo.Save(context.Background(), th)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	}

	cases := map[string]struct {
		owner        string
		pageMetadata fleet.PageMetadata
		size         uint64
	}{
		"retrieve all agents with existing owner": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
			},
			size: n,
		},
		"retrieve subset of agents with existing owner": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: n / 2,
				Limit:  n,
				Total:  n,
			},
			size: n / 2,
		},
		"retrieve agents with existing metadata": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset:   0,
				Limit:    n,
				Total:    n,
				Metadata: subMeta,
			},
			size: n,
		},
		"retrieve agents with existing tags": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Tags:   subTags,
			},
			size: n,
		},
		"retrieve agents with non-existing owner": {
			owner: wrongoID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  0,
			},
			size: 0,
		},
		"retrieve agents with non-existing name": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Name:   "wrong",
				Total:  0,
			},
			size: 0,
		},
		"retrieve agents with non-existing metadata": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset:   0,
				Limit:    n,
				Total:    0,
				Metadata: wrongMeta,
			},
			size: 0,
		},
		"retrieve agents with non-existing tags": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  0,
				Tags:   wrongTags,
			},
			size: 0,
		},
		"retrieve agents sorted by name ascendent": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Order:  "name",
				Dir:    "asc",
			},
			size: n,
		},
		"retrieve agents sorted by name descendent": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Order:  "name",
				Dir:    "desc",
			},
			size: n,
		},
		"retrieve agents with mix tags": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Tags:   mixTags,
			},
			size: n,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			page, err := agentRepo.RetrieveAll(context.Background(), tc.owner, tc.pageMetadata)
			size := uint64(len(page.Agents))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d\n", desc, tc.size, size))
			assert.Equal(t, tc.pageMetadata.Total, page.Total, fmt.Sprintf("%s: expected total %d got %d\n", desc, tc.pageMetadata.Total, page.Total))
			assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))

			// Check if Agents list have been sorted properly
			if size > 0 {
				testSortAgents(t, tc.pageMetadata, page.Agents)
			}
		})
	}
}

func TestAgentUpdate(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	duplicatedThID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-agent1")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	updatedNameID, err := types.NewIdentifier("my-agent2")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	duplicatedNameID, err := types.NewIdentifier("my-agent3")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	duplicatedAgent := fleet.Agent{
		Name:          duplicatedNameID,
		MFThingID:     duplicatedThID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	err = agentRepo.Save(context.Background(), duplicatedAgent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		agent fleet.Agent
		err   error
	}{
		"update existing agent data by thingID": {
			agent: fleet.Agent{
				MFThingID: thID.String(),
				MFOwnerID: oID.String(),
				Name:      updatedNameID,
				OrbTags:   types.Tags{"newkey": "newvalue"},
			},
			err: nil,
		},
		"update non-existent agent data by thingID": {
			agent: fleet.Agent{
				MFThingID: oID.String(),
				MFOwnerID: oID.String(),
				Name:      updatedNameID,
				OrbTags:   types.Tags{"newkey": "newvalue"},
			},
			err: errors.ErrNotFound,
		},
		"update agent data with empty thingID": {
			agent: fleet.Agent{
				MFThingID: "",
				MFOwnerID: oID.String(),
				Name:      updatedNameID,
				OrbTags:   types.Tags{"newkey": "newvalue"},
			},
			err: errors.ErrMalformedEntity,
		},
		"update agent data with empty OwnerID": {
			agent: fleet.Agent{
				MFThingID: thID.String(),
				MFOwnerID: "",
				Name:      updatedNameID,
				OrbTags:   types.Tags{"newkey": "newvalue"},
			},
			err: errors.ErrMalformedEntity,
		},
		"update agent data by thingID and channelID with invalid thingID": {
			agent: fleet.Agent{
				MFThingID:     "123",
				MFOwnerID:     oID.String(),
				Name:          updatedNameID,
				AgentMetadata: types.Metadata{"newkey": "newvalue"},
			},
			err: errors.ErrMalformedEntity,
		},
		"update agent data by thingID and channelID with duplicated nameID": {
			agent: fleet.Agent{
				MFThingID: thID.String(),
				MFOwnerID: oID.String(),
				Name:      duplicatedNameID,
			},
			err: errors.ErrConflict,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err = agentRepo.UpdateAgentByID(context.Background(), tc.agent.MFOwnerID, tc.agent)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
			if err == nil {
				ag, err := agentRepo.RetrieveByID(context.Background(), tc.agent.MFOwnerID, tc.agent.MFThingID)
				assert.Nil(t, err)
				assert.Equal(t, tc.agent.OrbTags, ag.OrbTags, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
			}
		})
	}
}

func TestDeleteAgent(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalidID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-agent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	ag := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), ag)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		ID      string
		ownerID string
		err     error
	}{
		"remove existing agent": {
			ID:      ag.MFThingID,
			ownerID: ag.MFOwnerID,
			err:     nil,
		},
		"remove a non-existing agent": {
			ID:      invalidID.String(),
			ownerID: ag.MFOwnerID,
			err:     nil,
		},
	}

	for desc, tc := range cases {
		err := agentRepo.Delete(context.Background(), tc.ownerID, tc.ID)
		require.Nil(t, err, fmt.Sprintf("%s: failed to remove agent due to: %s", desc, err))
	}
}

func TestAgentBackendTapsRetrieve(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalidID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-agent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	ag := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), ag)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		ownerID string
		total   int
		err     error
	}{
		"retrieve a list of taps by ownerID": {
			ownerID: ag.MFOwnerID,
			total:   1,
			err:     nil,
		},
		"retrieve a list of taps by a wrong ownerID": {
			ownerID: invalidID.String(),
			total:   0,
			err:     nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			ag, err := agentRepo.RetrieveAgentMetadataByOwner(context.Background(), tc.ownerID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
			assert.Equal(t, tc.total, len(ag), fmt.Sprintf("%s: expected %d got %d", desc, tc.total, len(ag)))
		})
	}
}

func TestMultiAgentRetrievalByAgentGroup(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)
	agentGroupRepo := postgres.NewAgentGroupRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	wrongoID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	groupNameID, err := types.NewIdentifier("my-group")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	metaStr := `{"field1":"value1","field2":{"subfield11":"value2","subfield12":{"subfield121":"value3","subfield122":"value4"}}}`
	metadata := types.Metadata{}
	json.Unmarshal([]byte(metaStr), &metadata)

	tagsStr := `{"region": "EU", "node_type": "dns"}`
	tags := types.Tags{}
	json.Unmarshal([]byte(tagsStr), &tags)

	subTagsStr := `{"region": "EU"}`
	subTags := types.Tags{}
	json.Unmarshal([]byte(subTagsStr), &subTags)

	group := fleet.AgentGroup{
		Name:        groupNameID,
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		Tags:        tags,
	}

	id, err := agentGroupRepo.Save(context.Background(), group)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		thID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		chID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		th := fleet.Agent{
			MFOwnerID:   oID.String(),
			MFThingID:   thID.String(),
			MFChannelID: chID.String(),
		}

		th.Name, err = types.NewIdentifier(fmt.Sprintf("agent_name-%d", i))
		require.True(t, th.Name.IsValid(), "invalid Identifier name: %s")
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		th.AgentMetadata = metadata
		th.AgentTags = tags
		th.OrbTags = subTags
		th.State = fleet.Online

		err = agentRepo.Save(context.Background(), th)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	}

	cases := map[string]struct {
		owner    string
		groupID  string
		onlinish bool
		size     uint64
		err      error
	}{
		"retrieve all agents with existing owner that are online": {
			owner:    oID.String(),
			onlinish: true,
			groupID:  id,
			size:     n,
			err:      nil,
		},
		"retrieve all agents with empty owner": {
			owner:    "",
			onlinish: true,
			groupID:  id,
			size:     0,
			err:      errors.ErrMalformedEntity,
		},
		"retrieve agents with non-existing groupID": {
			owner:    oID.String(),
			groupID:  wrongoID.String(),
			onlinish: true,
			size:     0,
			err:      nil,
		},
		"retrieve agents with non-existing owner": {
			owner:    wrongoID.String(),
			groupID:  id,
			onlinish: true,
			size:     0,
			err:      nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			agents, err := agentRepo.RetrieveAllByAgentGroupID(context.Background(), tc.owner, tc.groupID, tc.onlinish)
			size := uint64(len(agents))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d\n", desc, tc.size, size))
			assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %d\n", desc, tc.err, err))

		})
	}
}

func TestAgentRetrieveByID(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		thingID string
		ownerID string
		err     error
		tags    types.Tags
	}{
		"retrieve existing agent by thingID": {
			thingID: thID.String(),
			ownerID: oID.String(),
			tags:    types.Tags{"testkey": "testvalue"},
			err:     nil,
		},
		"retrieve non-existent agent by thingID": {
			thingID: thID.String(),
			ownerID: thID.String(),
			err:     errors.ErrNotFound,
		},
		"retrieve existing agent by thingID with invalid ownerID": {
			thingID: thID.String(),
			ownerID: "123",
			err:     errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			ag, err := agentRepo.RetrieveByID(context.Background(), tc.ownerID, tc.thingID)
			if err == nil {
				assert.Equal(t, nameID, ag.Name, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
			}
			if len(tc.tags) > 0 {
				assert.Equal(t, tc.tags, ag.OrbTags)
				assert.Equal(t, tc.tags, ag.AgentTags)
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func TestRetrieveOwnerByChannelID(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		channelID string
		ownerID   string
		name      string
		err       error
	}{
		"retrieve existing owner by channelID": {
			channelID: chID.String(),
			ownerID:   oID.String(),
			name:      nameID.String(),
			err:       nil,
		},
		"retrieve existent owner by non-existent channelID": {
			channelID: thID.String(),
			ownerID:   "",
			name:      "",
			err:       nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			ag, err := agentRepo.RetrieveOwnerByChannelID(context.Background(), tc.channelID)
			if err == nil {
				assert.Equal(t, tc.name, ag.Name.String(), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.name, ag.Name.String()))
				assert.Equal(t, tc.ownerID, ag.MFOwnerID, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.ownerID, ag.MFOwnerID))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func testSortAgents(t *testing.T, pm fleet.PageMetadata, ths []fleet.Agent) {
	switch pm.Order {
	case "name":
		current := ths[0]
		for _, res := range ths {
			if pm.Dir == "asc" {
				assert.GreaterOrEqual(t, res.Name.String(), current.Name.String())
			}
			if pm.Dir == "desc" {
				assert.GreaterOrEqual(t, current.Name.String(), res.Name.String())
			}
			current = res
		}
	default:
		break
	}
}

func TestMatchingAgentRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	name := "agent_name"
	orbTagsStr := `{"node_type": "dns"}`
	agentTagsStr := `{"region": "EU"}`
	mixTagsStr := `{"node_type": "dns", "region": "EU"}`

	orbTags := types.Tags{}
	json.Unmarshal([]byte(orbTagsStr), &orbTags)

	agentTags := types.Tags{}
	json.Unmarshal([]byte(agentTagsStr), &agentTags)

	mixTags := types.Tags{}
	json.Unmarshal([]byte(mixTagsStr), &mixTags)

	n := uint64(3)
	for i := uint64(0); i < n; i++ {
		thID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		chID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		th := fleet.Agent{
			MFOwnerID:   oID.String(),
			MFThingID:   thID.String(),
			MFChannelID: chID.String(),
		}

		th.Name, err = types.NewIdentifier(fmt.Sprintf("%s-%d", name, i))
		require.True(t, th.Name.IsValid(), "invalid Identifier name: %s")
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		th.AgentTags = agentTags
		th.OrbTags = orbTags

		err = agentRepo.Save(context.Background(), th)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	}

	cases := map[string]struct {
		owner          string
		tag            types.Tags
		matchingAgents types.Metadata
	}{
		"retrieve matching agents with mix tags": {
			owner: oID.String(),
			tag:   mixTags,
			matchingAgents: types.Metadata{
				"total":  float64(n),
				"online": float64(0),
			},
		},
		"retrieve matching agents with orb tags": {
			owner: oID.String(),
			tag:   orbTags,
			matchingAgents: types.Metadata{
				"total":  float64(n),
				"online": float64(0),
			},
		},
		"retrieve matching agents with agent tags": {
			owner: oID.String(),
			tag:   agentTags,
			matchingAgents: types.Metadata{
				"total":  float64(n),
				"online": float64(0),
			},
		},
		"retrieve unmatched agents with mix tags": {
			owner: oID.String(),
			tag: types.Tags{
				"wrong": "tag",
			},
			matchingAgents: types.Metadata{
				"total":  nil,
				"online": nil,
			},
		},
		"retrieve agents with mix tags": {
			owner: oID.String(),
			tag: types.Tags{
				"node_type": "dns",
				"region":    "EU",
				"wrong":     "tag",
			},
			matchingAgents: types.Metadata{
				"total":  nil,
				"online": nil,
			},
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			ma, err := agentRepo.RetrieveMatchingAgents(context.Background(), tc.owner, tc.tag)
			assert.True(t, reflect.DeepEqual(tc.matchingAgents, ma), fmt.Sprintf("%s: expected %v got %v\n", desc, tc.matchingAgents, ma))
			assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
		})
	}
}

func TestSetAgentStale(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
		LastHB:        time.Now().Add(fleet.DefaultTimeout),
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	agent.State = fleet.Online
	err = agentRepo.UpdateHeartbeatByIDWithChannel(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	time.Sleep(2 * time.Second)

	cases := map[string]struct {
		agent        fleet.Agent
		duration     time.Duration
		owner        string
		state        fleet.State
		affectedRows int64
	}{
		"set agent state to stale when stops do send heartbeats": {
			agent:        agent,
			duration:     1 * time.Second,
			owner:        oID.String(),
			state:        fleet.Stale,
			affectedRows: 1,
		},
		"keep agent state online when agent it's sending heartbeats": {
			agent:        agent,
			duration:     3 * time.Second,
			owner:        oID.String(),
			state:        fleet.Online,
			affectedRows: 0,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			count, err := agentRepo.SetStaleStatus(context.Background(), tc.duration)
			require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
			require.Equal(t, tc.affectedRows, count, fmt.Sprintf("%s: expected affected rows %d got %d", desc, tc.affectedRows, count))
			agent, err := agentRepo.RetrieveByID(context.Background(), tc.owner, tc.agent.MFThingID)
			require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
			require.Equal(t, tc.state, agent.State, fmt.Sprintf("%s: expected %s got %s", desc, tc.state.String(), agent.State.String()))
		})
	}

}
