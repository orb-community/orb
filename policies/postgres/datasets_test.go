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

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	nameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	conflictNameID, err := types.NewIdentifier("mydataset-conflict")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         nameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID.String(),
		SinkIDs:      sinkIDs,
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}

	// Conflict scenario
	datasetCopy := dataset
	datasetCopy.Name = conflictNameID

	_, err = repo.SaveDataset(context.Background(), datasetCopy)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	cases := map[string]struct {
		dataset policies.Dataset
		err     error
	}{
		"create new dataset": {
			dataset: dataset,
			err:     nil,
		},
		"create dataset that already exist": {
			dataset: datasetCopy,
			err:     errors.ErrConflict,
		},
		"create new dataset with empty ownerID": {
			dataset: policies.Dataset{
				Name:         nameID,
				MFOwnerID:    "",
				Valid:        true,
				AgentGroupID: groupID.String(),
				PolicyID:     policyID.String(),
				SinkIDs:      sinkIDs,
			},
			err: errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := repo.SaveDataset(context.Background(), tc.dataset)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
		})
	}
}

func TestDatasetUpdate(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	nameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         nameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID.String(),
		SinkIDs:      sinkIDs,
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}

	dsID, err := repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset.ID = dsID

	cases := map[string]struct {
		dataset policies.Dataset
		err     error
	}{
		"update a existing dataset": {
			dataset: dataset,
			err:     nil,
		},
		"update a non-existing dataset": {
			dataset: policies.Dataset{
				Name:         nameID,
				MFOwnerID:    oID.String(),
				Valid:        true,
				AgentGroupID: groupID.String(),
				PolicyID:     policyID.String(),
				SinkIDs:      sinkIDs,
				Metadata:     types.Metadata{"testkey": "testvalue"},
				Created:      time.Time{},
				ID:           wrongID.String(),
			},
			err: policies.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := repo.UpdateDataset(context.Background(), tc.dataset.MFOwnerID, tc.dataset)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
		})
	}
}

func TestDatasetDelete(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	nameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         nameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID.String(),
		SinkIDs:      sinkIDs,
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}

	dsID, err := repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset.ID = dsID

	cases := map[string]struct {
		owner string
		id    string
		err   error
	}{
		"delete a existing dataset": {
			owner: dataset.MFOwnerID,
			id:    dataset.ID,
			err:   nil,
		},
		"delete a non-existing dataset": {
			owner: dataset.MFOwnerID,
			id:    wrongID.String(),
			err:   errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := repo.DeleteDataset(context.Background(), tc.owner, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
		})
	}
}

func TestDatasetRetrieveByID(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("mypolicy")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	dataset := policies.Dataset{
		Name:         nameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID.String(),
		SinkIDs:      sinkIDs,
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}

	id, err := repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		datasetID string
		ownerID   string
		err       error
	}{
		"retrieve existing dataset by ID and ownerID": {
			datasetID: id,
			ownerID:   dataset.MFOwnerID,
			err:       nil,
		},
		"retrieve non-existent dataset by ID and ownerID": {
			datasetID: dataset.MFOwnerID,
			ownerID:   id,
			err:       errors.ErrNotFound,
		},
		"retrieve dataset by ID and ownerID with emmpty owner field": {
			datasetID: id,
			ownerID:   "",
			err:       errors.ErrMalformedEntity,
		},
		"retrieve dataset by ID and ownerID with emmpty datasetID field": {
			datasetID: "",
			ownerID:   dataset.MFOwnerID,
			err:       errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			tcd, err := repo.RetrieveDatasetByID(context.Background(), tc.datasetID, tc.ownerID)
			if err == nil {
				assert.Equal(t, dataset.Name, tcd.Name, fmt.Sprintf("%s: unexpected name change expected %s got %s", desc, dataset.Name, tcd.Name))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func TestMultiDatasetRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		nameID, err := types.NewIdentifier(fmt.Sprintf("mydataset-%d", i))
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		dataset := policies.Dataset{
			Name:      nameID,
			MFOwnerID: oID.String(),
		}

		_, err = repo.SaveDataset(context.Background(), dataset)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	}

	cases := map[string]struct {
		owner        string
		pageMetadata policies.PageMetadata
		size         uint64
	}{
		"retrieve all datasets with existing owner": {
			owner: oID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
			},
			size: n,
		},
		"retrieve subset of datasets with existing owner": {
			owner: oID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: n / 2,
				Limit:  n,
				Total:  n,
			},
			size: n / 2,
		},
		"retrieve datasets with no-existing owner": {
			owner: wrongID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  0,
			},
			size: 0,
		},
		"retrieve datasets with no-existing name": {
			owner: oID.String(),
			pageMetadata: policies.PageMetadata{
				Offset: 0,
				Limit:  n,
				Name:   "wrong",
				Total:  0,
			},
			size: 0,
		},
		"retrieve datasets sorted by name ascendent": {
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
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			pageDataset, err := repo.RetrieveAllDatasetsByOwner(context.Background(), tc.owner, tc.pageMetadata)
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s\n", desc, err))
			size := uint64(len(pageDataset.Datasets))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d", desc, tc.size, size))
			assert.Equal(t, tc.pageMetadata.Total, pageDataset.Total, fmt.Sprintf("%s: expected total %d got %d", desc, tc.pageMetadata.Total, pageDataset.Total))

			if size > 0 {
				testSortDataset(t, tc.pageMetadata, pageDataset.Datasets)
			}
		})
	}
}

func TestInactivateDatasetByGroupID(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongOID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongGroupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	nameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         nameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID.String(),
		SinkIDs:      sinkIDs,
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}

	dsID, err := repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset.ID = dsID

	cases := map[string]struct {
		ownerID string
		groupID string
		err     error
	}{
		"inactivate a existing dataset by group ID": {
			ownerID: dataset.MFOwnerID,
			groupID: dataset.AgentGroupID,
			err:     nil,
		},
		"inactivate a dataset with non-existent owner": {
			groupID: dataset.AgentGroupID,
			ownerID: wrongOID.String(),
			err:     policies.ErrInactivateDataset,
		},
		"inactivate a non-existing dataset with existent owner": {
			groupID: wrongGroupID.String(),
			ownerID: dataset.MFOwnerID,
			err:     policies.ErrInactivateDataset,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := repo.InactivateDatasetByGroupID(context.Background(), tc.groupID, tc.ownerID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
		})
	}
}

func TestInactivateDatasetByPolicyID(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	nameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         nameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID.String(),
		SinkIDs:      sinkIDs,
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}

	dsID, err := repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset.ID = dsID

	cases := map[string]struct {
		ownerID  string
		policyID string
		err      error
	}{
		"inactivate a existing dataset by policy ID": {
			ownerID:  dataset.MFOwnerID,
			policyID: dataset.PolicyID,
			err:      nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := repo.InactivateDatasetByPolicyID(context.Background(), tc.policyID, tc.ownerID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
		})
	}
}

func TestMultiDatasetRetrievalPolicyID(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		nameID, err := types.NewIdentifier(fmt.Sprintf("mydataset-%d", i))
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		dataset := policies.Dataset{
			Name:      nameID,
			MFOwnerID: oID.String(),
			PolicyID:  policyID.String(),
		}

		_, err = repo.SaveDataset(context.Background(), dataset)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	}

	cases := map[string]struct {
		owner    string
		policyID string
		size     uint64
		err      error
	}{
		"retrieve all datasets by policyID": {
			owner:    oID.String(),
			policyID: policyID.String(),
			size:     n,
			err:      nil,
		},
		"retrieve datasets with no-existing owner": {
			owner:    wrongID.String(),
			policyID: policyID.String(),
			size:     0,
			err:      nil,
		},
		"retrieve all datasets by policyID with empty policyID": {
			owner:    oID.String(),
			policyID: "",
			size:     0,
			err:      errors.ErrMalformedEntity,
		},
		"retrieve datasets with no-existing policyID": {
			owner:    "",
			policyID: policyID.String(),
			size:     0,
			err:      errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			datasets, err := repo.RetrieveDatasetsByPolicyID(context.Background(), tc.policyID, tc.owner)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
			size := uint64(len(datasets))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d", desc, tc.size, size))
		})
	}
}

func TestInactivateDatasetBySinkID(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID2, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	nameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID2, err := types.NewIdentifier("mydataset-2")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         nameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID.String(),
		SinkIDs:      sinkIDs,
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}

	dataset2 := dataset
	dataset2.Name = nameID2
	dataset2.MFOwnerID = oID2.String()

	dsID, err := repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset.ID = dsID

	dsID2, err := repo.SaveDataset(context.Background(), dataset2)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset2.ID = dsID2

	cases := map[string]struct {
		ownerID string
		id      string
		dataset policies.Dataset
		valid   bool
		err     error
	}{
		"inactivate dataset with a non-existing sink": {
			id:      wrongID.String(),
			ownerID: dataset.MFOwnerID,
			dataset: dataset2,
			valid:   true,
			err:     nil,
		},
		"inactivate a existing dataset by sink ID": {
			ownerID: dataset.MFOwnerID,
			id:      dataset.ID,
			dataset: dataset,
			valid:   false,
			err:     nil,
		},
		"inactivate dataset with an invalid ownerID": {
			id:      dataset.ID,
			ownerID: "",
			dataset: dataset2,
			valid:   true,
			err:     policies.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := repo.InactivateDatasetByID(context.Background(), tc.id, tc.ownerID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))

			assertate, _ := repo.RetrieveDatasetByID(context.Background(), tc.dataset.ID, tc.dataset.MFOwnerID)
			assert.Equal(t, tc.valid, assertate.Valid, fmt.Sprintf("%s: expected '%t' got '%t'", desc, tc.valid, assertate.Valid))
		})
	}
}

func TestDeleteSinkFromDataset(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	repo := postgres.NewPoliciesRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID2, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongSinkID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	nameID, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID2, err := types.NewIdentifier("mydataset")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         nameID,
		MFOwnerID:    oID.String(),
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policyID.String(),
		SinkIDs:      sinkIDs,
		Metadata:     types.Metadata{"testkey": "testvalue"},
		Created:      time.Time{},
	}

	dataset2 := dataset
	dataset2.Name = nameID2
	dataset2.MFOwnerID = oID2.String()

	dsID, err := repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset.ID = dsID

	dsID2, err := repo.SaveDataset(context.Background(), dataset2)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset2.ID = dsID2

	cases := map[string]struct {
		owner    string
		sinkID   string
		contains bool
		dataset  policies.Dataset
		err      error
	}{
		"delete a sink from existing dataset": {
			owner:    dataset.MFOwnerID,
			sinkID:   dataset.SinkIDs[0],
			contains: false,
			dataset:  dataset,
			err:      nil,
		},
		"delete a non-existing sink from a dataset": {
			owner:    dataset.MFOwnerID,
			sinkID:   wrongSinkID.String(),
			contains: false,
			dataset:  dataset,
			err:      nil,
		},
		"delete a sink from a dataset with an invalid ownerID": {
			sinkID:   dataset2.SinkIDs[0],
			owner:    "",
			contains: true,
			dataset:  dataset2,
			err:      errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			dataset, err := repo.DeleteSinkFromDataset(context.Background(), tc.sinkID, tc.owner)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))

			for _, d := range dataset {
				switch tc.contains {
				case false:
					assert.NotContains(t, d.SinkIDs, tc.sinkID, fmt.Sprintf("%s: expected '%s' to not contains '%s'", desc, d.SinkIDs, tc.sinkID))
				case true:
					assert.Contains(t, d.SinkIDs, tc.sinkID, fmt.Sprintf("%s: expected '%s' to contains '%s'", desc, d.SinkIDs, tc.sinkID))
				}
			}
		})
	}
}

func testSortDataset(t *testing.T, pm policies.PageMetadata, ags []policies.Dataset) {
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
