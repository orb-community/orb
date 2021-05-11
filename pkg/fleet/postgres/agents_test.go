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
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/fleet"
	"github.com/ns1labs/orb/pkg/fleet/postgres"
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

	agent := fleet.Agent{
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		OrbTags:       fleet.Tags{"testkey": "testvalue"},
		AgentTags:     fleet.Tags{"testkey": "testvalue"},
		AgentMetadata: fleet.Metadata{"testkey": "testvalue"},
	}

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
			agent: fleet.Agent{MFThingID: "invalid", MFOwnerID: oID.String()},
			err:   fleet.ErrMalformedEntity,
		},
		{
			desc:  "create agent with invalid owner ID",
			agent: fleet.Agent{MFThingID: thID.String(), MFOwnerID: "invalid"},
			err:   fleet.ErrMalformedEntity,
		},
	}

	for _, tc := range cases {
		err := agentRepo.Save(context.Background(), tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", tc.desc, tc.err, err))
	}
}
