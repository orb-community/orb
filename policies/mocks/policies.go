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
	policyCounter  uint64
	pdb            map[string]policies.Policy
	dataSetCounter uint64
	ddb            map[string]policies.Dataset
	gdb            map[string][]policies.Policy
}

func (m *mockPoliciesRepository) InactivateDatasetByPolicyID(ctx context.Context, policyID string, ownerID string) error {
	//todo implement when create unit tests to dataset
	return nil
}

func (m *mockPoliciesRepository) DeletePolicy(ctx context.Context, ownerID string, policyID string) error {
	if _, ok := m.pdb[policyID]; ok {
		if m.pdb[policyID].MFOwnerID != ownerID {
			delete(m.gdb, policyID)
		}
	}
	return nil
}

func (m *mockPoliciesRepository) RetrieveDatasetsByPolicyID(ctx context.Context, policyID string, ownerID string) ([]policies.Dataset, error) {
	var datasetList []policies.Dataset
	for _, d := range m.ddb{
		if d.PolicyID == policyID && d.MFOwnerID == ownerID{
			datasetList = append(datasetList, d)
		}
	}

	return datasetList, nil
}

func (m *mockPoliciesRepository) UpdatePolicy(ctx context.Context, owner string, pol policies.Policy) error {
	if _, ok := m.pdb[pol.ID]; ok {
		if m.pdb[pol.ID].MFOwnerID != owner {
			return policies.ErrUpdateEntity
		}
		pol.MFOwnerID = owner
		m.pdb[pol.ID] = pol
		return nil
	}
	return policies.ErrNotFound
}

func NewPoliciesRepository() policies.Repository {
	return &mockPoliciesRepository{
		pdb: make(map[string]policies.Policy),
		ddb: make(map[string]policies.Dataset),
		gdb: make(map[string][]policies.Policy),
	}
}

func (m *mockPoliciesRepository) RetrieveAll(ctx context.Context, owner string, pm policies.PageMetadata) (policies.Page, error) {
	first := uint64(pm.Offset)
	last := first + uint64(pm.Limit)

	var policyList []policies.Policy
	id := uint64(0)
	for _, p := range m.pdb {
		if p.MFOwnerID == owner && id >= first && id < last {
			policyList = append(policyList, p)
		}
		id++
	}

	policyList = sortPolicies(pm, policyList)

	pagePolicies := policies.Page{
		PageMetadata: policies.PageMetadata{
			Total: m.policyCounter,
		},
		Policies: policyList,
	}
	return pagePolicies, nil
}

func (m *mockPoliciesRepository) SavePolicy(ctx context.Context, policy policies.Policy) (string, error) {
	ID, _ := uuid.NewV4()
	policy.ID = ID.String()
	m.pdb[policy.ID] = policy
	m.policyCounter++
	return ID.String(), nil
}

func (m *mockPoliciesRepository) RetrievePolicyByID(ctx context.Context, policyID string, ownerID string) (policies.Policy, error) {
	if _, ok := m.pdb[policyID]; ok {
		if m.pdb[policyID].MFOwnerID != ownerID {
			return policies.Policy{}, policies.ErrNotFound
		}
		return m.pdb[policyID], nil
	}
	return policies.Policy{}, policies.ErrNotFound
}

func (m *mockPoliciesRepository) RetrievePoliciesByGroupID(ctx context.Context, groupIDs []string, ownerID string) (ret []policies.Policy, err error) {
	for _, d := range groupIDs {
		ret = append(ret, m.pdb[d])
	}
	return ret, nil
}

func (m *mockPoliciesRepository) SaveDataset(ctx context.Context, dataset policies.Dataset) (string, error) {
	ID, _ := uuid.NewV4()
	dataset.ID = ID.String()
	m.ddb[dataset.ID] = dataset
	m.gdb[dataset.AgentGroupID] = make([]policies.Policy, 1)
	m.gdb[dataset.AgentGroupID][0] = m.pdb[dataset.PolicyID]
	m.dataSetCounter++
	return ID.String(), nil
}

func (m *mockPoliciesRepository) RetrieveDatasetByID(ctx context.Context, datasetID string, ownerID string) (policies.Dataset, error) {
	if _, ok := m.ddb[datasetID]; ok {
		if m.ddb[datasetID].MFOwnerID != ownerID {
			return policies.Dataset{}, policies.ErrNotFound
		}
		return m.ddb[datasetID], nil
	}
	return policies.Dataset{}, policies.ErrNotFound
}

func (m *mockPoliciesRepository) InactivateDatasetByGroupID(ctx context.Context, groupID string, ownerID string) error {
	panic("implement me")
}

func (m *mockPoliciesRepository) RetrieveAllDatasetByOwner(ctx context.Context, owner string, pm policies.PageMetadata) (policies.PageDataset, error) {
	first := uint64(pm.Offset)
	last := first + uint64(pm.Limit)

	var datasetList []policies.Dataset
	id := uint64(0)
	for _, p := range m.ddb {
		if p.MFOwnerID == owner && id >= first && id < last {
			datasetList = append(datasetList, p)
		}
		id++
	}

	datasetList = sortDataset(pm, datasetList)

	pageDataset := policies.PageDataset{
		PageMetadata: policies.PageMetadata{
			Total:  m.dataSetCounter,
			Offset: pm.Offset,
			Limit:  pm.Limit,
			Dir:    pm.Dir,
		},
		Datasets: datasetList,
	}
	return pageDataset, nil
}