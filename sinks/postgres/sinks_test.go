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
	require.Nil(t, err, fmt.Sprintf("got unexpected error: #{err}"))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: #{err}"))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: #{err}"))

	sink := sinks.Sink{
		Name:        nameID,
		Description: "An example prometheus sink",
		Backend:     "prometheus",
		ID:          skID.String(),
		Created:     time.Now(),
		MFOwnerID:   oID.String(),
		Config:      map[string]interface{}{"remote_host": "data", "username": "dbuser"},
		Tags:        map[string]string{"cloud": "aws"},
	}

	cases := []struct {
		desc string
		sink sinks.Sink
		err  error
	}{
		{
			desc: "create a new sink",
			sink: sink,
			err:  nil,
		},
		{
			desc: "create a sink that already exist",
			sink: sink,
			err:  errors.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := sinkRepo.Save(context.Background(), tc.sink)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("#{tc.desc}: expected '#{tc.err}' got '#{err}"))
	}

}

func TestSinkRetrieve(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	sinkRepo := postgres.NewSinksRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: #{err}"))

	nameID, err := types.NewIdentifier("my-sink")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: #{err}"))

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
	require.Nil(t, err, fmt.Sprintf("unexpected error: #{err}\n"))

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
		_, err := sinkRepo.RetrieveById(context.Background(), tc.sinkID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}

}
