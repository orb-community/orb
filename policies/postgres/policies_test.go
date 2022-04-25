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
	"math"
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

	conflictNameID, err := types.NewIdentifier("mypolicy-conflict")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policy := policies.Policy{
		Name:      nameID,
		MFOwnerID: oID.String(),
		OrbTags:   types.Tags{"testkey": "testvalue"},
	}

	policyCopy := policy
	policyCopy.Name = conflictNameID

	_, err = repo.SavePolicy(context.Background(), policyCopy)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	cases := map[string]struct {
		policy policies.Policy
		err    error
	}{
		"create new policy": {
			policy: policy,
			err:    nil,
		},
		"create policy that already exist": {
			policy: policyCopy,
			err:    errors.ErrConflict,
		},
		"create new policy with empty owner": {
			policy: policies.Policy{
				Name:      nameID,
				MFOwnerID: "",
			},
			err: errors.ErrMalformedEntity,
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
		"retrieve policies by ID with empty owner": {
			policyID: id,
			ownerID:  "",
			tags:     types.Tags{"testkey": "testvalue"},
			err:      errors.ErrMalformedEntity,
		},
		"retrieve policies by ID with empty policyID": {
			policyID: "",
			ownerID:  policy.MFOwnerID,
			tags:     types.Tags{"testkey": "testvalue"},
			err:      errors.ErrMalformedEntity,
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

func TestMultiPolicyRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		nameID, err := types.NewIdentifier(fmt.Sprintf("mypolicy-%d", i))
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		policy := policies.Policy{
			Name:      nameID,
			MFOwnerID: oID.String(),
			Policy:    types.Metadata{"pkey1": "pvalue1"},
		}
		if math.Mod(float64(i), 2) == 0 {
			policy.OrbTags = types.Tags{"node_type": "dns"}
		}

		_, err = repo.SavePolicy(context.Background(), policy)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	}

	cases := map[string]struct {
		owner        string
		pageMetadata policies.PageMetadata
		size         uint64
	}{
		"retrieve all policies with existing owner": {
			owner: oID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
			},
			size: n,
		},
		"retrieve subset of policies with existing owner": {
			owner: oID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: n / 2,
				Limit:  n,
				Total:  n,
			},
			size: n / 2,
		},
		"retrieve policies with no-existing owner": {
			owner: wrongID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  0,
			},
			size: 0,
		},
		"retrieve policies with no-existing name": {
			owner: oID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: 0,
				Limit:  n,
				Name:   "wrong",
				Total:  0,
			},
			size: 0,
		},
		"retrieve policies sorted by name ascendent": {
			owner: oID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Order:  "name",
				Dir:    "asc",
			},
			size: n,
		},
		"retrieve policies sorted by name descendent": {
			owner: oID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Order:  "name",
				Dir:    "desc",
			},
			size: n,
		},
		"retrieve policies filtered by tag": {
			owner: oID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n / 2,
				Tags:   types.Tags{"node_type": "dns"},
			},
			size: n / 2,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			page, err := repo.RetrieveAll(context.Background(), tc.owner, tc.pageMetadata)
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s\n", desc, err))
			size := uint64(len(page.Policies))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d", desc, tc.size, size))
			assert.Equal(t, tc.pageMetadata.Total, page.Total, fmt.Sprintf("%s: expected total %d got %d", desc, tc.pageMetadata.Total, page.Total))

			if size > 0 {
				testSortPolicies(t, tc.pageMetadata, page.Policies)
			}
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

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	dsnameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         dsnameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID,
		SinkIDs:      sinkIDs,
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}
	dsID, err := repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		groupID []string
		ownerID string
		dsID    string
		results int
		err     error
	}{
		"retrieve existing policies by group ID and ownerID": {
			groupID: []string{groupID.String()},
			ownerID: policy.MFOwnerID,
			dsID:    dsID,
			results: 1,
			err:     nil,
		},
		"retrieve non existing policies by group ID and ownerID": {
			groupID: []string{policy.MFOwnerID},
			ownerID: policy.MFOwnerID,
			dsID:    dsID,
			results: 0,
			err:     nil,
		},
		"retrieve policies by groupID with empty owner": {
			groupID: []string{policy.MFOwnerID},
			ownerID: "",
			dsID:    dsID,
			results: 0,
			err:     errors.ErrMalformedEntity,
		},
		"retrieve policies by groupID with empty groupID": {
			groupID: []string{},
			ownerID: policy.MFOwnerID,
			dsID:    dsID,
			results: 0,
			err:     errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			plist, err := repo.RetrievePoliciesByGroupID(context.Background(), tc.groupID, tc.ownerID)
			if err == nil {
				assert.Equal(t, tc.results, len(plist), fmt.Sprintf("%s: expected %d got %d\n", desc, tc.results, len(plist)))
				if tc.results > 0 {
					assert.Equal(t, policy.Name.String(), plist[0].Name.String(), fmt.Sprintf("%s: expected %s got %s\n", desc, policy.Name.String(), plist[0].Name.String()))
					assert.Equal(t, dsID, plist[0].DatasetID, fmt.Sprintf("%s: expected %s got %s\n", desc, policy.Name.String(), plist[0].Name.String()))
				}
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func TestPolicyUpdate(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("mypolicy")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := policies.Policy{
		Name:      nameID,
		MFOwnerID: oID.String(),
		Policy:    types.Metadata{"pkey1": "pvalue1"},
	}
	policyID, err := repo.SavePolicy(context.Background(), policy)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	policy.ID = policyID

	cases := map[string]struct {
		plcy policies.Policy
		err  error
	}{
		"update a existing policy": {
			plcy: policy,
			err:  nil,
		},
		"update a empty policy": {
			plcy: policies.Policy{},
			err:  policies.ErrUpdateEntity,
		},
		"update a non-existing policy": {
			plcy: policies.Policy{
				ID: wrongID.String(),
			},
			err: policies.ErrUpdateEntity,
		},
		"update policy with wrong owner": {
			plcy: policies.Policy{ID: wrongID.String(), MFOwnerID: wrongID.String()},
			err:  policies.ErrNotFound,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := repo.UpdatePolicy(context.Background(), tc.plcy.MFOwnerID, tc.plcy)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func TestPolicyDelete(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("mypolicy")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := policies.Policy{
		Name:      nameID,
		MFOwnerID: oID.String(),
		Policy:    types.Metadata{"pkey1": "pvalue1"},
	}
	policyID, err := repo.SavePolicy(context.Background(), policy)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	policy.ID = policyID

	cases := map[string]struct {
		id    string
		owner string
		err   error
	}{
		"delete a existing policy": {
			id:    policy.ID,
			owner: policy.MFOwnerID,
			err:   nil,
		},
		"delete a empty policy id": {
			id:    "",
			owner: policy.MFOwnerID,
			err:   policies.ErrMalformedEntity,
		},
		"delete a non-existing policy": {
			id:    wrongID.String(),
			owner: policy.MFOwnerID,
			err:   nil,
		},
		"delete policy with empty owner": {
			id:    policy.ID,
			owner: "",
			err:   policies.ErrMalformedEntity,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := repo.DeletePolicy(context.Background(), tc.owner, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func testSortPolicies(t *testing.T, pm policies.PageMetadata, ags []policies.Policy) {
	t.Helper()
	switch pm.Order {
	case "name":
		current := ags[0]
		for _, res := range ags {
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
