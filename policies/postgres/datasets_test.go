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

	conflictNameID, err := types.NewIdentifier("mydataset-conflict")
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
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := repo.SaveDataset(context.Background(), tc.dataset)
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

	dataset := policies.Dataset{
		Name:      nameID,
		MFOwnerID: oID.String(),
	}

	id, err := repo.SaveDataset(context.Background(), dataset)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		datasetID string
		ownerID   string
		err       error
	}{
		"retrieve existing policy by ID and ownerID": {
			datasetID: id,
			ownerID:   dataset.MFOwnerID,
			err:       nil,
		},
		"retrieve non-existent policy by ID and ownerID": {
			datasetID: dataset.MFOwnerID,
			ownerID:   id,
			err:       errors.ErrNotFound,
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
			pageDataset, err := repo.RetrieveAllDatasetByOwner(context.Background(), tc.owner, tc.pageMetadata)
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
