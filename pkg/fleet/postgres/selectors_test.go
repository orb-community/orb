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
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/fleet"
	"github.com/ns1labs/orb/pkg/fleet/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectorSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	selectorRepo := postgres.NewSelectorRepository(dbMiddleware)

	email := "agent-owner@example.com"

	agent := fleet.Selector{
		Name:  "myselector",
		Owner: email,
	}

	cases := []struct {
		desc  string
		agent fleet.Selector
		err   error
	}{
		{
			desc:  "create new selector",
			agent: agent,
			err:   nil,
		},
		{
			desc:  "create selector that already exist",
			agent: agent,
			err:   fleet.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := selectorRepo.Save(context.Background(), tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
