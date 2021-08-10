/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package mocks

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/policies"
)

var _ policies.Repository = (*mockPoliciesRepository)(nil)

type mockPoliciesRepository struct {
	pdb map[string]policies.Policy
	ddb map[string]policies.Dataset
	gdb map[string][]policies.Policy
}

func NewPoliciesRepository() policies.Repository {
	return &mockPoliciesRepository{
		pdb: make(map[string]policies.Policy),
		ddb: make(map[string]policies.Dataset),
		gdb: make(map[string][]policies.Policy),
	}
}

func (m mockPoliciesRepository) SavePolicy(ctx context.Context, policy policies.Policy) (string, error) {
	ID, _ := uuid.NewV4()
	policy.ID = ID.String()
	m.pdb[policy.ID] = policy
	return ID.String(), nil
}

func (m mockPoliciesRepository) RetrievePolicyByID(ctx context.Context, policyID string, ownerID string) (policies.Policy, error) {
	return m.pdb[policyID], nil
}

func (m mockPoliciesRepository) RetrievePoliciesByGroupID(ctx context.Context, groupIDs []string, ownerID string) (ret []policies.Policy, err error) {
	for _, d := range groupIDs {
		ret = append(ret, m.pdb[d])
	}
	return ret, nil
}

func (m mockPoliciesRepository) SaveDataset(ctx context.Context, dataset policies.Dataset) (string, error) {
	ID, _ := uuid.NewV4()
	dataset.ID = ID.String()
	m.ddb[dataset.ID] = dataset
	m.gdb[dataset.AgentGroupID] = make([]policies.Policy, 1)
	m.gdb[dataset.AgentGroupID][0] = m.pdb[dataset.PolicyID]
	return ID.String(), nil
}

func (m mockPoliciesRepository) UpdateDatasetToInactivate(ctx context.Context, groupID string, ownerID string) error {
	panic("implement me")
}
