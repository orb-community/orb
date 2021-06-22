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
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/postgres"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectorSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	selectorRepo := postgres.NewSelectorRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-selector")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	selector := fleet.Selector{
		Name:      nameID,
		MFOwnerID: oID.String(),
		Metadata:  types.Metadata{"testkey": "testvalue"},
	}

	cases := []struct {
		desc     string
		selector fleet.Selector
		err      error
	}{
		{
			desc:     "create new selector",
			selector: selector,
			err:      nil,
		},
		{
			desc:     "create selector that already exist",
			selector: selector,
			err:      errors.ErrConflict,
		},
		{
			desc:     "create selector with invalid name",
			selector: fleet.Selector{MFOwnerID: oID.String()},
			err:      errors.ErrMalformedEntity,
		}, {
			desc:     "create selector with invalid owner ID",
			selector: fleet.Selector{Name: nameID, MFOwnerID: "invalid"},
			err:      errors.ErrMalformedEntity,
		},
	}

	for _, tc := range cases {
		err := selectorRepo.Save(context.Background(), tc.selector)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", tc.desc, tc.err, err))
	}

}
