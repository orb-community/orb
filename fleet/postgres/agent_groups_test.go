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
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentGroupSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentGroupRepository := postgres.NewAgentGroupRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-group")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	group := fleet.AgentGroup{
		Name:        nameID,
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		Tags:        types.Tags{"testkey": "testvalue"},
	}

	cases := map[string]struct {
		agentGroup fleet.AgentGroup
		err        error
	}{
		"create new group": {
			agentGroup: group,
			err:        nil,
		},
		"create group that already exist": {
			agentGroup: group,
			err:        errors.ErrConflict,
		},
		"create group with invalid name": {
			agentGroup: fleet.AgentGroup{MFOwnerID: oID.String()},
			err:        errors.ErrMalformedEntity,
		},
		"create group with invalid owner ID": {
			agentGroup: fleet.AgentGroup{Name: nameID, MFOwnerID: "invalid"},
			err:        errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		_, err := agentGroupRepository.Save(context.Background(), tc.agentGroup)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
	}

}

func TestAgentGroupRetrieve(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentGroupRepo := postgres.NewAgentGroupRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-group")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	group := fleet.AgentGroup{
		Name:        nameID,
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		Tags:        types.Tags{"testkey": "testvalue"},
	}

	id, err := agentGroupRepo.Save(context.Background(), group)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		groupID string
		ownerID string
		err     error
		tags    types.Tags
	}{
		"retrieve existing agent group by groupID and ownerID": {
			groupID: id,
			ownerID: group.MFOwnerID,
			tags:    types.Tags{"testkey": "testvalue"},
			err:     nil,
		},
		"retrieve non-existent agent group by groupID and ownerID": {
			ownerID: id,
			groupID: group.MFOwnerID,
			err:     errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		ag, err := agentGroupRepo.RetrieveByID(context.Background(), tc.groupID, tc.ownerID)
		if err == nil {
			assert.Equal(t, nameID, ag.Name, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
		if len(tc.tags) > 0 {
			assert.Equal(t, tc.tags, ag.Tags)
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestAgentGroupValidate(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentGroupRepo := postgres.NewAgentGroupRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	existingNameID, err := types.NewIdentifier("existingAG")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	noNameID, err := types.NewIdentifier("nonexistentAG")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalidOID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	group := fleet.AgentGroup{
		Name:        existingNameID,
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		Tags:        types.Tags{"testkey": "testvalue"},
	}

	id, err := agentGroupRepo.Save(context.Background(), group)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	group.ID = id

	cases := map[string]struct {
		name    types.Identifier
		ownerID string
		tags    types.Tags
		err     error
	}{
		"validate a new agent group": {
			name:    noNameID,
			ownerID: group.MFOwnerID,
			tags:    types.Tags{"testkey": "testvalue"},
			err:     nil,
		},
		"validate an agent group that already exists": {
			name:    group.Name,
			ownerID: group.MFOwnerID,
			tags:    types.Tags{"testkey": "testvalue"},
			err:     fleet.ErrConflict,
		},
		"validate an agent group with an invalid owner id": {
			name:    group.Name,
			ownerID: invalidOID.String(),
			err:     nil,
		},
	}

	for desc, tc := range cases {
		err := agentGroupRepo.RetrieveToValidate(context.Background(), tc.name.String(), tc.ownerID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}
