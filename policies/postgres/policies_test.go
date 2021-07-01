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
	"go.uber.org/zap"
	"testing"
)

var (
	logger, _ = zap.NewDevelopment()
)

func TestPolicySave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("mypolicy")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policy := policies.Policy{
		Name:      nameID,
		MFOwnerID: oID.String(),
		OrbTags:   types.Tags{"testkey": "testvalue"},
	}

	cases := []struct {
		desc   string
		policy policies.Policy
		err    error
	}{
		{
			desc:   "create new policy",
			policy: policy,
			err:    nil,
		},
		{
			desc:   "create policy that already exist",
			policy: policy,
			err:    errors.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := repo.SavePolicy(context.Background(), tc.policy)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", tc.desc, tc.err, err))
	}
}
