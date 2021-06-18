// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/postgres"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"strings"
	"testing"
)

const maxNameSize = 1024

var (
	invalidName = strings.Repeat("m", maxNameSize+1)
	logger, _   = zap.NewDevelopment()
)

func TestAgentSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		OrbTags:       fleet.Tags{"testkey": "testvalue"},
		AgentTags:     fleet.Tags{"testkey": "testvalue"},
		AgentMetadata: fleet.Metadata{"testkey": "testvalue"},
	}

	name1, _ := types.NewIdentifier("myagent2")

	cases := []struct {
		desc  string
		agent fleet.Agent
		err   error
	}{
		{
			desc:  "create new agent",
			agent: agent,
			err:   nil,
		},
		{
			desc:  "create agent that already exist",
			agent: agent,
			err:   fleet.ErrConflict,
		},
		{
			desc:  "create agent with null thing ID",
			agent: fleet.Agent{Name: name1, MFThingID: "", MFOwnerID: oID.String()},
			err:   nil,
		},
	}

	for _, tc := range cases {
		err := agentRepo.Save(context.Background(), tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", tc.desc, tc.err, err))
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
		OrbTags:       fleet.Tags{"testkey": "testvalue"},
		AgentTags:     fleet.Tags{"testkey": "testvalue"},
		AgentMetadata: fleet.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		thingID   string
		channelID string
		err       error
	}{
		"retrieve existing agent by thingID and channelID": {
			thingID:   agent.MFThingID,
			channelID: agent.MFChannelID,
			err:       nil,
		},
		"retrieve non-existent agent by thingID and channelID": {
			thingID:   agent.MFOwnerID,
			channelID: agent.MFChannelID,
			err:       fleet.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		ag, err := agentRepo.RetrieveByIDWithChannel(context.Background(), tc.thingID, tc.channelID)
		if err == nil {
			assert.Equal(t, nameID, ag.Name, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
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
		OrbTags:       fleet.Tags{"testkey": "testvalue"},
		AgentTags:     fleet.Tags{"testkey": "testvalue"},
		AgentMetadata: fleet.Metadata{"testkey": "testvalue"},
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
				AgentMetadata: fleet.Metadata{"newkey": "newvalue"},
			},
			err: nil,
		},
		"update non-existent agent data by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:     chID.String(),
				MFChannelID:   thID.String(),
				AgentMetadata: fleet.Metadata{"newkey": "newvalue"},
			},
			err: fleet.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		err = agentRepo.UpdateDataByIDWithChannel(context.Background(), tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		if err == nil {
			ag, err := agentRepo.RetrieveByIDWithChannel(context.Background(), tc.agent.MFThingID, tc.agent.MFChannelID)
			assert.Nil(t, err)
			assert.Equal(t, tc.agent.AgentMetadata, ag.AgentMetadata, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
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
		LastHBData:  fleet.Metadata{"heartbeatdata": "testvalue"},
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
				LastHBData:  fleet.Metadata{"heartbeatdata2": "newvalue"},
			},
			err: nil,
		},
		"update non-existent agent heart beat by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:   chID.String(),
				MFChannelID: thID.String(),
				LastHBData:  fleet.Metadata{"heartbeatdata2": "newvalue"},
			},
			err: fleet.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		err = agentRepo.UpdateHeartbeatByIDWithChannel(context.Background(), tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		if err == nil {
			ag, err := agentRepo.RetrieveByIDWithChannel(context.Background(), tc.agent.MFThingID, tc.agent.MFChannelID)
			assert.Nil(t, err)
			assert.Equal(t, tc.agent.LastHBData, ag.LastHBData, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
	}
}
