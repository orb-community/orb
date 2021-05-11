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
	"go.uber.org/zap"
)

var _ fleet.AgentRepository = (*agentRepository)(nil)

type agentRepository struct {
	db     Database
	logger *zap.Logger
}

func NewAgentRepository(db Database, logger *zap.Logger) fleet.AgentRepository {
	return &agentRepository{db: db, logger: logger}
}

func (r agentRepository) Save(ctx context.Context, agent fleet.Agent) error {

	q := `INSERT INTO agents (mf_thing_id, mf_owner_id, orb_tags, agent_tags, agent_metadata)         
			  VALUES (:mf_thing_id, :mf_owner_id, :orb_tags, :agent_tags, :agent_metadata)`

	if agent.MFThingID == "" || agent.MFOwnerID == "" {
		return fleet.ErrMalformedEntity
	}

	dba, err := toDBAgent(agent)
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

type dbAgent struct {
	MFThingID     uuid.UUID  `db:"mf_thing_id"`
	MFOwnerID     uuid.UUID  `db:"mf_owner_id"`
	OrbTags       dbMetadata `db:"orb_tags"`
	AgentTags     dbMetadata `db:"agent_tags"`
	AgentMetadata dbMetadata `db:"agent_metadata"`
}

func toDBAgent(agent fleet.Agent) (dbAgent, error) {

	var tID uuid.UUID
	err := tID.Scan(agent.MFThingID)
	if err != nil {
		return dbAgent{}, errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	var oID uuid.UUID
	err = oID.Scan(agent.MFOwnerID)
	if err != nil {
		return dbAgent{}, errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	return dbAgent{
		MFThingID:     tID,
		MFOwnerID:     oID,
		OrbTags:       dbMetadata(agent.OrbTags),
		AgentTags:     dbMetadata(agent.AgentTags),
		AgentMetadata: dbMetadata(agent.AgentMetadata),
	}, nil

}
