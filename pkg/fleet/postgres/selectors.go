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
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/fleet"
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
		return fleet.ErrMalformedEntity
	}

	dba, err := toDBSelector(selector)
	if err != nil {
		return errors.Wrap(errSaveDB, err)
	}

	_, err = r.db.NamedExecContext(ctx, q, dba)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return errors.Wrap(fleet.ErrMalformedEntity, err)
			case errDuplicate:
				return errors.Wrap(fleet.ErrConflict, err)
			}
		}
		return errors.Wrap(errSaveDB, err)
	}

	return nil
}

type dbSelector struct {
	Name      types.Identifier `db:"name"`
	MFOwnerID uuid.UUID        `db:"mf_owner_id"`
	Metadata  dbMetadata       `db:"metadata"`
}

func toDBSelector(selector fleet.Selector) (dbSelector, error) {

	var oID uuid.UUID
	err := oID.Scan(selector.MFOwnerID)
	if err != nil {
		return dbSelector{}, errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	return dbSelector{
		Name:      selector.Name,
		MFOwnerID: oID,
		Metadata:  dbMetadata(selector.Metadata),
	}, nil

}
