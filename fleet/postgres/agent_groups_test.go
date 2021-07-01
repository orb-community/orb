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

	cases := []struct {
		desc       string
		agentGroup fleet.AgentGroup
		err        error
	}{
		{
			desc:       "create new group",
			agentGroup: group,
			err:        nil,
		},
		{
			desc:       "create group that already exist",
			agentGroup: group,
			err:        errors.ErrConflict,
		},
		{
			desc:       "create group with invalid name",
			agentGroup: fleet.AgentGroup{MFOwnerID: oID.String()},
			err:        errors.ErrMalformedEntity,
		}, {
			desc:       "create group with invalid owner ID",
			agentGroup: fleet.AgentGroup{Name: nameID, MFOwnerID: "invalid"},
			err:        errors.ErrMalformedEntity,
		},
	}

	for _, tc := range cases {
		_, err := agentGroupRepository.Save(context.Background(), tc.agentGroup)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", tc.desc, tc.err, err))
	}

}
