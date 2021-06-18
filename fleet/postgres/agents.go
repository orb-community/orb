// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
	"time"
)

var _ fleet.AgentRepository = (*agentRepository)(nil)

type agentRepository struct {
	db     Database
	logger *zap.Logger
}

func NewAgentRepository(db Database, logger *zap.Logger) fleet.AgentRepository {
	return &agentRepository{db: db, logger: logger}
}

func (r agentRepository) UpdateDataByIDWithChannel(ctx context.Context, agent fleet.Agent) error {

	q := `UPDATE agents SET (orb_tags, agent_tags, agent_metadata)         
			= (:orb_tags, :agent_tags, :agent_metadata) 
			WHERE mf_thing_id = :mf_thing_id AND mf_channel_id = :mf_channel_id;`

	if agent.MFThingID == "" || agent.MFChannelID == "" {
		return fleet.ErrMalformedEntity
	}

	dba, err := toDBAgent(agent)
	if err != nil {
		return errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	res, err := r.db.NamedExecContext(ctx, q, dba)
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
		return errors.Wrap(errUpdateDB, err)
	}

	cnt, errdb := res.RowsAffected()
	if errdb != nil {
		return errors.Wrap(fleet.ErrUpdateEntity, errdb)
	}

	if cnt == 0 {
		return fleet.ErrNotFound
	}

	return nil
}

func (r agentRepository) UpdateHeartbeatByIDWithChannel(ctx context.Context, agent fleet.Agent) error {

	q := `UPDATE agents SET (last_hb_data, ts_last_hb, state)         
			= (:last_hb_data, now(), 'online') 
			WHERE mf_thing_id = :mf_thing_id AND mf_channel_id = :mf_channel_id;`

	if agent.MFThingID == "" || agent.MFChannelID == "" {
		return fleet.ErrMalformedEntity
	}

	dba, err := toDBAgent(agent)
	if err != nil {
		return errors.Wrap(fleet.ErrUpdateEntity, err)
	}
	res, err := r.db.NamedExecContext(ctx, q, dba)
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
		return errors.Wrap(errUpdateDB, err)
	}

	cnt, errdb := res.RowsAffected()
	if errdb != nil {
		return errors.Wrap(fleet.ErrUpdateEntity, errdb)
	}

	if cnt == 0 {
		return fleet.ErrNotFound
	}

	return nil

}

func (r agentRepository) RetrieveByIDWithChannel(ctx context.Context, thingID string, channelID string) (fleet.Agent, error) {

	q := `SELECT * FROM agents WHERE mf_thing_id = $1 AND mf_channel_id = $2;`

	dba := dbAgent{}

	if err := r.db.QueryRowxContext(ctx, q, thingID, channelID).StructScan(&dba); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return fleet.Agent{}, errors.Wrap(fleet.ErrNotFound, err)
		}
		return fleet.Agent{}, errors.Wrap(fleet.ErrSelectEntity, err)
	}

	return toAgent(dba)
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
	MFOwnerID     string           `db:"mf_owner_id"`
	MFThingID     sql.NullString   `db:"mf_thing_id"`
	MFChannelID   sql.NullString   `db:"mf_channel_id"`
	OrbTags       dbTags           `db:"orb_tags"`
	AgentTags     dbTags           `db:"agent_tags"`
	AgentMetadata dbMetadata       `db:"agent_metadata"`
	State         fleet.State      `db:"state"`
	Created       time.Time        `db:"ts_created"`
	LastHBData    dbMetadata       `db:"last_hb_data"`
	LastHB        sql.NullTime     `db:"ts_last_hb"`
}

func toDBAgent(agent fleet.Agent) (dbAgent, error) {

	a := dbAgent{
		Name:          agent.Name,
		MFOwnerID:     agent.MFOwnerID,
		OrbTags:       dbTags(agent.OrbTags),
		AgentTags:     dbTags(agent.AgentTags),
		AgentMetadata: dbMetadata(agent.AgentMetadata),
		State:         agent.State,
		Created:       agent.Created,
		LastHBData:    dbMetadata(agent.LastHBData),
	}

	if agent.MFThingID != "" {
		a.MFThingID = sql.NullString{
			String: agent.MFThingID,
			Valid:  true,
		}
	}
	if agent.MFChannelID != "" {
		a.MFChannelID = sql.NullString{
			String: agent.MFChannelID,
			Valid:  true,
		}
	}
	if !agent.LastHB.IsZero() {
		a.LastHB = sql.NullTime{
			Time:  agent.LastHB,
			Valid: true,
		}
	}
	return a, nil
}

func toAgent(dba dbAgent) (fleet.Agent, error) {

	agent := fleet.Agent{
		Name:          dba.Name,
		MFOwnerID:     dba.MFOwnerID,
		MFThingID:     dba.MFThingID.String,
		MFChannelID:   dba.MFChannelID.String,
		Created:       dba.Created,
		OrbTags:       fleet.Tags(dba.OrbTags),
		AgentTags:     fleet.Tags(dba.AgentTags),
		AgentMetadata: fleet.Metadata(dba.AgentMetadata),
		State:         dba.State,
		LastHBData:    fleet.Metadata(dba.LastHBData),
		LastHB:        dba.LastHB.Time,
	}

	return agent, nil

}
