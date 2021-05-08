// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/fleet"
	"go.uber.org/zap"
)

const (
	duplicateErr = "unique_violation"
	fkViolation  = "foreign_key_violation"
)

var (
	errSaveDB = errors.New("failed to save sink to database")
)

var _ fleet.FleetRepository = (*fleetRepository)(nil)

type fleetRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewFleetRepository(db *sqlx.DB, log *zap.Logger) fleet.FleetRepository {
	return &fleetRepository{db: db, logger: log}
}

func (cr fleetRepository) Save(cfg fleet.Agent) (string, error) {
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
