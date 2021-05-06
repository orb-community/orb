// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/ns1labs/orb/pkg/sinks"
)

const (
	duplicateErr = "unique_violation"
	fkViolation  = "foreign_key_violation"
)

var (
	errSaveDB = errors.New("failed to save sink to database")
)

var _ sinks.SinksRepository = (*sinksRepository)(nil)

type sinksRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewSinksRepository(db *sqlx.DB, log logger.Logger) sinks.SinksRepository {
	return &sinksRepository{db: db, log: log}
}

func (cr sinksRepository) Save(cfg sinks.Sink) (string, error) {
	/*
		q := `INSERT INTO sinks (mainflux_thing, owner, name, client_cert, client_key, ca_cert, mainflux_key, external_id, external_key, content, state)
			  VALUES (:mainflux_thing, :owner, :name, :client_cert, :client_key, :ca_cert, :mainflux_key, :external_id, :external_key, :content, :state)`

		tx, err := cr.db.Beginx()
		if err != nil {
			return "", errors.Wrap(errSaveDB, err)
		}

		return cfg.MFThing, nil

	*/
	return "", nil
}
