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
			desc:  "create agent with invalid thing ID",
			agent: fleet.Agent{Name: name1, MFThingID: "invalid", MFOwnerID: oID.String()},
			err:   fleet.ErrMalformedEntity,
		},
		{
			desc:  "create agent with invalid owner ID",
			agent: fleet.Agent{Name: name1, MFThingID: thID.String(), MFOwnerID: "invalid"},
			err:   fleet.ErrMalformedEntity,
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
		ownerID   string
		err       error
	}{
		"retrieve existing agent by thingID and channelID": {
			thingID:   agent.MFThingID,
			channelID: agent.MFChannelID,
			ownerID:   "",
			err:       nil,
		},
		"retrieve existing agent by thingID and ownerID": {
			thingID:   agent.MFThingID,
			channelID: "",
			ownerID:   agent.MFOwnerID,
			err:       nil,
		},
		"retrieve non-existent agent by thingID and ownerID": {
			thingID:   agent.MFOwnerID,
			channelID: agent.MFChannelID,
			ownerID:   "",
			err:       fleet.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		var ag fleet.Agent
		if tc.channelID != "" {
			ag, err = agentRepo.RetrieveByIDWithChannel(context.Background(), tc.thingID, tc.channelID)
		} else {
			ag, err = agentRepo.RetrieveByIDWithOwner(context.Background(), tc.thingID, tc.ownerID)
		}
		if err == nil {
			assert.Equal(t, nameID, ag.Name, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}
