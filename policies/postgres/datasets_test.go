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
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
	"github.com/ns1labs/orb/policies/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDatasetSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         nameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID.String(),
		SinkID:       sinkID.String(),
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}

	cases := []struct {
		desc    string
		dataset policies.Dataset
		err     error
	}{
		{
			desc:    "create new dataset",
			dataset: dataset,
			err:     nil,
		},
		{
			desc:    "create dataset that already exist",
			dataset: dataset,
			err:     errors.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := repo.SaveDataset(context.Background(), tc.dataset)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", tc.desc, tc.err, err))
	}
}
