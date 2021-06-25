// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"context"
	"github.com/lib/pq"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
)

var _ fleet.AgentGroupRepository = (*agentGroupRepository)(nil)

type agentGroupRepository struct {
	db     Database
	logger *zap.Logger
}

func NewAgentGroupRepository(db Database, logger *zap.Logger) fleet.AgentGroupRepository {
	return &agentGroupRepository{db: db, logger: logger}
}

func (r agentGroupRepository) Save(ctx context.Context, group fleet.AgentGroup) (string, error) {

	q := `INSERT INTO agent_groups (name, mf_owner_id, mf_channel_id, tags)         
			  VALUES (:name, :mf_owner_id, :mf_channel_id, :tags) RETURNING id`

	if !group.Name.IsValid() || group.MFOwnerID == "" || group.MFChannelID == "" {
		return "", errors.ErrMalformedEntity
	}

	dba, err := toDBAgentGroup(group)
	if err != nil {
		return "", errors.Wrap(db.ErrSaveDB, err)
	}

	row, err := r.db.NamedQueryContext(ctx, q, dba)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return "", errors.Wrap(errors.ErrMalformedEntity, err)
			case db.ErrDuplicate:
				return "", errors.Wrap(errors.ErrConflict, err)
			}
		}
		return "", errors.Wrap(db.ErrSaveDB, err)
	}

	defer row.Close()
	row.Next()
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

type dbAgentGroup struct {
	ID          string           `db:"id"`
	Name        types.Identifier `db:"name"`
	MFOwnerID   string           `db:"mf_owner_id"`
	MFChannelID string           `db:"mf_channel_id"`
	Tags        db.Tags          `db:"tags"`
}

func toDBAgentGroup(group fleet.AgentGroup) (dbAgentGroup, error) {

	return dbAgentGroup{
		ID:          group.ID,
		Name:        group.Name,
		MFOwnerID:   group.MFOwnerID,
		MFChannelID: group.MFChannelID,
		Tags:        db.Tags(group.Tags),
	}, nil

}
