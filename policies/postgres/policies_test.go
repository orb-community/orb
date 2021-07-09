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

func TestAgentPolicyDataRetrieve(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("mypolicy")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policy := policies.Policy{
		Name:      nameID,
		MFOwnerID: oID.String(),
		Policy:    types.Metadata{"pkey1": "pvalue1"},
	}

	id, err := repo.SavePolicy(context.Background(), policy)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		policyID string
		ownerID  string
		err      error
		tags     types.Tags
	}{
		"retrieve existing policy by ID and ownerID": {
			policyID: id,
			ownerID:  policy.MFOwnerID,
			tags:     types.Tags{"testkey": "testvalue"},
			err:      nil,
		},
		"retrieve non-existent policy by ID and ownerID": {
			ownerID:  id,
			policyID: policy.MFOwnerID,
			err:      errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		name, data, err := repo.RetrievePolicyDataByID(context.Background(), tc.policyID, tc.ownerID)
		if err == nil {
			assert.Equal(t, policy.Name.String(), name, fmt.Sprintf("%s: unexpected name change expected %s got %s", desc, policy.Name.String(), name))
			var pdata types.Metadata
			if err := json.Unmarshal(data, &pdata); err != nil {
				assert.Error(t, err, "unable to unmarshal policy")
			}
			assert.Equal(t, policy.Policy, pdata, fmt.Sprintf("%s: expected %s got %s\n", desc, policy.Policy, pdata))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}
