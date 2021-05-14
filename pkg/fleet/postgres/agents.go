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

var _ fleet.AgentRepository = (*agentRepository)(nil)

type agentRepository struct {
	db     Database
	logger *zap.Logger
}

func NewAgentRepository(db Database, logger *zap.Logger) fleet.AgentRepository {
	return &agentRepository{db: db, logger: logger}
}

func (r agentRepository) Save(ctx context.Context, agent fleet.Agent) error {

	q := `INSERT INTO agents (name, mf_thing_id, mf_owner_id, mf_channel_id, orb_tags, agent_tags, agent_metadata, state)         
			  VALUES (:name, :mf_thing_id, :mf_owner_id, :mf_channel_id, :orb_tags, :agent_tags, :agent_metadata, :state)`

	if !agent.Name.IsValid() || agent.MFOwnerID == "" {
		return fleet.ErrMalformedEntity
	}

	dba, err := toDBAgent(agent)
	if err != nil {
		return errors.Wrap(errSaveDB, err)
	}

	// enforce removed state if no ThingID
	if !dba.MFThingID.Valid {
		dba.State = fleet.Removed
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
	Name          types.Identifier `db:"name"`
	MFOwnerID     uuid.UUID        `db:"mf_owner_id"`
	MFThingID     uuid.NullUUID    `db:"mf_thing_id"`
	MFChannelID   uuid.NullUUID    `db:"mf_channel_id"`
	OrbTags       dbMetadata       `db:"orb_tags"`
	AgentTags     dbMetadata       `db:"agent_tags"`
	AgentMetadata dbMetadata       `db:"agent_metadata"`
	State         fleet.State      `db:"state"`
}

func toDBAgent(agent fleet.Agent) (dbAgent, error) {

	var tID uuid.NullUUID
	if agent.MFThingID == "" {
		tID = uuid.NullUUID{UUID: uuid.Nil, Valid: false}
	} else {
		err := tID.Scan(agent.MFThingID)
		if err != nil {
			return dbAgent{}, errors.Wrap(fleet.ErrMalformedEntity, err)
		}
	}

	var chID uuid.NullUUID
	if agent.MFChannelID == "" {
		tID = uuid.NullUUID{UUID: uuid.Nil, Valid: false}
	} else {
		err := tID.Scan(agent.MFChannelID)
		if err != nil {
			return dbAgent{}, errors.Wrap(fleet.ErrMalformedEntity, err)
		}
	}

	var oID uuid.UUID
	err := oID.Scan(agent.MFOwnerID)
	if err != nil {
		return dbAgent{}, errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	return dbAgent{
		Name:          agent.Name,
		MFThingID:     tID,
		MFChannelID:   chID,
		MFOwnerID:     oID,
		State:         agent.State,
		OrbTags:       dbMetadata(agent.OrbTags),
		AgentTags:     dbMetadata(agent.AgentTags),
		AgentMetadata: dbMetadata(agent.AgentMetadata),
	}, nil

}
