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
	"time"
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

	cases := map[string]struct {
		policy policies.Policy
		err    error
	}{
		"create new policy": {
			policy: policy,
			err:    nil,
		},
		"create policy that already exist": {
			policy: policy,
			err:    errors.ErrConflict,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := repo.SavePolicy(context.Background(), tc.policy)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
		})
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
		t.Run(desc, func(t *testing.T) {
			tcp, err := repo.RetrievePolicyByID(context.Background(), tc.policyID, tc.ownerID)
			if err == nil {
				assert.Equal(t, policy.Name, tcp.Name, fmt.Sprintf("%s: unexpected name change expected %s got %s", desc, policy.Name, tcp.Name))
				assert.Equal(t, policy.Policy, tcp.Policy, fmt.Sprintf("%s: expected %s got %s\n", desc, policy.Policy, tcp.Policy))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func TestAgentPoliciesRetrieveByGroup(t *testing.T) {
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
	policyID, err := repo.SavePolicy(context.Background(), policy)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dsnameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         dsnameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID,
		SinkID:       sinkID.String(),
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}
	_, err = repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		groupID []string
		ownerID string
		results int
		err     error
	}{
		"retrieve existing policies by group ID and ownerID": {
			groupID: []string{groupID.String()},
			ownerID: policy.MFOwnerID,
			results: 1,
			err:     nil,
		},
		"retrieve non existing policies by group ID and ownerID": {
			groupID: []string{policy.MFOwnerID},
			ownerID: policy.MFOwnerID,
			results: 0,
			err:     nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			plist, err := repo.RetrievePoliciesByGroupID(context.Background(), tc.groupID, tc.ownerID)
			if err == nil {
				assert.Equal(t, tc.results, len(plist), fmt.Sprintf("%s: expected %d got %d\n", desc, tc.results, len(plist)))
				if tc.results > 0 {
					assert.Equal(t, policy.Name.String(), plist[0].Name.String(), fmt.Sprintf("%s: expected %s got %s\n", desc, policy.Name.String(), plist[0].Name.String()))
				}
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}
