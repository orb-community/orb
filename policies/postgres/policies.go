// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
	"go.uber.org/zap"
)

var _ policies.Repository = (*policiesRepository)(nil)

type policiesRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewPoliciesRepository(db *sqlx.DB, log *zap.Logger) policies.Repository {
	return &policiesRepository{db: db, logger: log}
}

func (r policiesRepository) Save(ctx context.Context, policy policies.Policy) error {

	q := `INSERT INTO policies (name, mf_owner_id, policy_yaml)         
			  VALUES (:name, :mf_owner_id, :policy_yaml)`

	if !policy.Name.IsValid() || policy.MFOwnerID == "" {
		return errors.ErrMalformedEntity
	}

	dba, err := toDBPolicy(policy)
	if err != nil {
		return errors.Wrap(db.ErrSaveDB, err)
	}

	_, err = r.db.NamedExecContext(ctx, q, dba)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(errors.ErrMalformedEntity, err)
			case db.ErrDuplicate:
				return errors.Wrap(errors.ErrConflict, err)
			}
		}
		return errors.Wrap(db.ErrSaveDB, err)
	}

	return nil

}

type dbPolicy struct {
	Name       types.Identifier `db:"name"`
	MFOwnerID  string           `db:"mf_owner_id"`
	PolicyYAML string           `db:"policy_yaml"`
}

func toDBPolicy(policy policies.Policy) (dbPolicy, error) {

	var uID uuid.UUID
	err := uID.Scan(policy.MFOwnerID)
	if err != nil {
		return dbPolicy{}, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return dbPolicy{
		Name:       policy.Name,
		MFOwnerID:  uID.String(),
		PolicyYAML: policy.PolicyYAML,
	}, nil

}
