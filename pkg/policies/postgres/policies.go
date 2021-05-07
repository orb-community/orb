// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/policies"
	"go.uber.org/zap"
)

const (
	duplicateErr = "unique_violation"
	fkViolation  = "foreign_key_violation"
)

var (
	errSaveDB = errors.New("failed to save sink to database")
)

var _ policies.PoliciesRepository = (*policiesRepository)(nil)

type policiesRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewPoliciesRepository(db *sqlx.DB, log *zap.Logger) policies.PoliciesRepository {
	return &policiesRepository{db: db, logger: log}
}

func (cr policiesRepository) Save(cfg policies.Policy) (string, error) {
	/*
		q := `INSERT INTO policies (sink_thing, owner, name, client_cert, client_key, ca_cert, sink_key, external_id, external_key, content, state)
			  VALUES (:sink_thing, :owner, :name, :client_cert, :client_key, :ca_cert, :sink_key, :external_id, :external_key, :content, :state)`

		tx, err := cr.db.Beginx()
		if err != nil {
			return "", errors.Wrap(errSaveDB, err)
		}

		return cfg.MFThing, nil

	*/
	return "", nil
}
