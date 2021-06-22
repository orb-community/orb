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
	"github.com/lib/pq"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
)

var _ fleet.SelectorRepository = (*selectorRepository)(nil)

type selectorRepository struct {
	db     Database
	logger *zap.Logger
}

func NewSelectorRepository(db Database, logger *zap.Logger) fleet.SelectorRepository {
	return &selectorRepository{db: db, logger: logger}
}

func (r selectorRepository) Save(ctx context.Context, selector fleet.Selector) error {

	q := `INSERT INTO selectors (name, mf_owner_id, metadata)         
			  VALUES (:name, :mf_owner_id, :metadata)`

	if !selector.Name.IsValid() || selector.MFOwnerID == "" {
		return errors.ErrMalformedEntity
	}

	dba, err := toDBSelector(selector)
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

type dbSelector struct {
	Name      types.Identifier `db:"name"`
	MFOwnerID uuid.UUID        `db:"mf_owner_id"`
	Metadata  db.Metadata      `db:"metadata"`
}

func toDBSelector(selector fleet.Selector) (dbSelector, error) {

	var oID uuid.UUID
	err := oID.Scan(selector.MFOwnerID)
	if err != nil {
		return dbSelector{}, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return dbSelector{
		Name:      selector.Name,
		MFOwnerID: oID,
		Metadata:  db.Metadata(selector.Metadata),
	}, nil

}
