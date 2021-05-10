// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"context"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/fleet"
)

const (
	duplicateErr = "unique_violation"
	fkViolation  = "foreign_key_violation"
)

var (
	errSaveDB = errors.New("failed to save agent to database")
)

var _ fleet.AgentRepository = (*agentRepository)(nil)

type agentRepository struct {
	db Database
}

func NewAgentRepository(db Database) fleet.AgentRepository {
	return &agentRepository{db: db}
}

func (cr agentRepository) Save(ctx context.Context, cfg fleet.Agent) (string, error) {
	/*
		q := `INSERT INTO fleet (sink_thing, owner, name, client_cert, client_key, ca_cert, sink_key, external_id, external_key, content, state)
			  VALUES (:sink_thing, :owner, :name, :client_cert, :client_key, :ca_cert, :sink_key, :external_id, :external_key, :content, :state)`

		tx, err := cr.db.Beginx()
		if err != nil {
			return "", errors.Wrap(errSaveDB, err)
		}

		return cfg.MFThing, nil

	*/
	return "", nil
}
