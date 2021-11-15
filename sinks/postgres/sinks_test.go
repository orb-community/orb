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
	"github.com/ns1labs/orb/sinks"
	"github.com/ns1labs/orb/sinks/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

var (
	logger, _ = zap.NewDevelopment()
)

func TestSinkSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	skID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	conflictNameID, err := types.NewIdentifier("my-sink-conflict")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        nameID,
		Description: "An example prometheus sink",
		Backend:     "prometheus",
		ID:          skID.String(),
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		State:       sinks.Unknown,
		Error:       "",
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkCopy := sink
	sinkCopy.Name = conflictNameID
	_, err = sinkRepo.Save(context.Background(), sinkCopy)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		sink sinks.Sink
		err  error
	}{
		"create a new sink": {
			sink: sink,
			err:  nil,
		},
		"create a sink that already exist": {
			sink: sinkCopy,
			err:  errors.ErrConflict,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := sinkRepo.Save(context.Background(), tc.sink)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}

}

func TestSinkUpdate(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalideOwnerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalideID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        nameID,
		Description: "An example prometheus sink",
		Backend:     "prometheus",
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkID, err := sinkRepo.Save(context.Background(), sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	sink.ID = sinkID

	cases := map[string]struct {
		sink sinks.Sink
		err  error
	}{
		"update a existing sink": {
			sink: sink,
			err:  nil,
		},
		"update a non-existing sink with a existing user": {
			sink: sinks.Sink{
				ID:        invalideID.String(),
				MFOwnerID: oID.String(),
			},
			err: sinks.ErrNotFound,
		},
		"update a existing sink with a non-existing user": {
			sink: sinks.Sink{
				ID:        sinkID,
				MFOwnerID: invalideOwnerID.String(),
			},
			err: sinks.ErrNotFound,
		},
		"update a non-existing sink with a non-existing user": {
			sink: sinks.Sink{
				ID:        invalideID.String(),
				MFOwnerID: invalideOwnerID.String(),
			},
			err: sinks.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := sinkRepo.Update(context.Background(), tc.sink)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}

}

func TestSinkRetrieve(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        nameID,
		Description: "An example prometheus sink",
		Backend:     "prometheus",
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkID, err := sinkRepo.Save(context.Background(), sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		sinkID string
		nameID string
		err    error
	}{
		"retrive existing sink by sinkID": {
			sinkID: sinkID,
			nameID: sink.Name.String(),
			err:    nil,
		},
		"retrive non-existing sink by sinkID": {
			sinkID: "",
			nameID: sink.Name.String(),
			err:    errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := sinkRepo.RetrieveById(context.Background(), tc.sinkID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}

}

func TestMultiSinkRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	wrongoID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {

		nameID, err := types.NewIdentifier(fmt.Sprintf("my-sink-%d", i))
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		sink := sinks.Sink{
			Name:        nameID,
			Description: "An example prometheus sink",
			Backend:     "prometheus",
			Created:     time.Now(),
			MFOwnerID:   oID.String(),
			Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
			Tags:        map[string]string{"cloud": "aws"},
		}

		_, err = sinkRepo.Save(context.Background(), sink)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	}

	cases := map[string]struct {
		owner        string
		pageMetadata sinks.PageMetadata
		size         uint64
	}{
		"retrieve all sinks with existing owner": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
			},
			size: n,
		},
		"retrieve subset of sinks with existing owner": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: n / 2,
				Limit:  n,
				Total:  n,
			},
			size: n / 2,
		},
		"retrieve sinks with no-existing owner": {
			owner: wrongoID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  0,
			},
			size: 0,
		},
		"retrieve sinks with no-existing name": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Name:   "wrong",
				Total:  0,
			},
			size: 0,
		},
		"retrieve agents sorted by name ascendent": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Order:  "name",
				Dir:    "asc",
			},
			size: n,
		},
		"retrieve agents sorted by name descendent": {
			owner: oID.String(),
			pageMetadata: sinks.PageMetadata{
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
			page, err := sinkRepo.RetrieveAll(context.Background(), tc.owner, tc.pageMetadata)
			require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
			size := uint64(len(page.Sinks))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d", desc, tc.size, size))
			assert.Equal(t, tc.pageMetadata.Total, page.Total, fmt.Sprintf("%s: expected total %d got %d", desc, tc.pageMetadata.Total, page.Total))

			if size > 0 {
				testSortSinks(t, tc.pageMetadata, page.Sinks)
			}
		})
	}
}

func TestSinkRemoval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sinkName, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sink := sinks.Sink{
		Name:        sinkName,
		Description: "An example prometheus sink",
		Backend:     "prometheus",
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	sinkID, err := sinkRepo.Save(context.Background(), sink)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	sink.ID = sinkID

	cases := map[string]struct {
		sink sinks.Sink
		err  error
	}{
		"delete existing sink": {
			sink: sink,
			err:  nil,
		},
		"delete non-existent sink": {
			sink: sink,
			err:  nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := sinkRepo.Remove(context.Background(), tc.sink.MFOwnerID, tc.sink.ID)
			require.Nil(t, err, fmt.Sprintf("%s: failed to remove sink due to: %s", desc, err))

			_, err = sinkRepo.RetrieveById(context.Background(), tc.sink.ID)
			require.True(t, errors.Contains(err, sinks.ErrNotFound), fmt.Sprintf("%s: expected %s got %s", desc, sinks.ErrNotFound, err))
		})
	}
}

func testSortSinks(t *testing.T, pm sinks.PageMetadata, sks []sinks.Sink) {
	t.Helper()
	switch pm.Order {
	case "name":
		current := sks[0]
		for _, res := range sks {
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
