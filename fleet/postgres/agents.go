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
	"github.com/gofrs/uuid"
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

func (r agentRepository) RetrieveByIDWithOwner(ctx context.Context, thingID string, ownerID string) (fleet.Agent, error) {
	q := `SELECT name, mf_channel_id, orb_tags, agent_tags, agent_metadata, state FROM agents WHERE mf_thing_id = $1 AND mf_owner_id = $2;`

	dba, err := newDBAgentByOwner(thingID, ownerID)
	if err != nil {
		return fleet.Agent{}, errors.Wrap(errMarshal, err)
	}

	if err := r.db.QueryRowxContext(ctx, q, thingID, ownerID).StructScan(&dba); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return fleet.Agent{}, errors.Wrap(fleet.ErrNotFound, err)
		}
		return fleet.Agent{}, errors.Wrap(fleet.ErrSelectEntity, err)
	}

	return toAgent(dba)
}

func (r agentRepository) RetrieveByIDWithChannel(ctx context.Context, thingID string, channelID string) (fleet.Agent, error) {
	q := `SELECT name, mf_owner_id, orb_tags, agent_tags, agent_metadata, state FROM agents WHERE mf_thing_id = $1 AND mf_channel_id = $2;`

	dba, err := newDBAgentByChannel(thingID, channelID)
	if err != nil {
		return fleet.Agent{}, errors.Wrap(errMarshal, err)
	}

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
	MFOwnerID     uuid.UUID        `db:"mf_owner_id"`
	MFThingID     uuid.NullUUID    `db:"mf_thing_id"`
	MFChannelID   uuid.NullUUID    `db:"mf_channel_id"`
	OrbTags       dbMetadata       `db:"orb_tags"`
	AgentTags     dbMetadata       `db:"agent_tags"`
	AgentMetadata dbMetadata       `db:"agent_metadata"`
	State         fleet.State      `db:"state"`
	Created       time.Time        `db:"ts_last_hb"`
	LastHBData    dbMetadata       `db:"last_hb_data"`
	LastHB        time.Time        `db:"ts_last_hb"`
}

func getUUID(u string) (uuid.UUID, error) {
	var oID uuid.UUID
	err := oID.Scan(u)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(fleet.ErrMalformedEntity, err)
	}
	return oID, nil
}

func getNullUUID(u string) (uuid.NullUUID, error) {
	var tID uuid.NullUUID
	if u == "" {
		tID = uuid.NullUUID{UUID: uuid.Nil, Valid: false}
	} else {
		err := tID.Scan(u)
		if err != nil {
			return uuid.NullUUID{}, errors.Wrap(fleet.ErrMalformedEntity, err)
		}
	}
	return tID, nil
}

func newDBAgentByChannel(thingID string, channelID string) (dbAgent, error) {

	tID, err := getNullUUID(thingID)
	if err != nil {
		return dbAgent{}, err
	}
	chID, err := getNullUUID(channelID)
	if err != nil {
		return dbAgent{}, err
	}

	return dbAgent{
		MFThingID:   tID,
		MFChannelID: chID,
	}, nil
}

func newDBAgentByOwner(thingID string, ownerID string) (dbAgent, error) {

	tID, err := getNullUUID(thingID)
	if err != nil {
		return dbAgent{}, err
	}
	oID, err := getUUID(ownerID)
	if err != nil {
		return dbAgent{}, err
	}

	return dbAgent{
		MFThingID: tID,
		MFOwnerID: oID,
	}, nil
}

func toDBAgent(agent fleet.Agent) (dbAgent, error) {

	tID, err := getNullUUID(agent.MFThingID)
	if err != nil {
		return dbAgent{}, err
	}
	chID, err := getNullUUID(agent.MFChannelID)
	if err != nil {
		return dbAgent{}, err
	}
	oID, err := getUUID(agent.MFOwnerID)
	if err != nil {
		return dbAgent{}, err
	}

	return dbAgent{
		Name:          agent.Name,
		MFOwnerID:     oID,
		MFThingID:     tID,
		MFChannelID:   chID,
		OrbTags:       dbMetadata(agent.OrbTags),
		AgentTags:     dbMetadata(agent.AgentTags),
		AgentMetadata: dbMetadata(agent.AgentMetadata),
		State:         agent.State,
		Created:       agent.Created,
		LastHBData:    dbMetadata(agent.LastHBData),
		LastHB:        agent.LastHB,
	}, nil

}

func toAgent(dba dbAgent) (fleet.Agent, error) {

	agent := fleet.Agent{
		Name:          dba.Name,
		MFOwnerID:     dba.MFOwnerID.String(),
		MFThingID:     dba.MFThingID.UUID.String(),
		MFChannelID:   dba.MFChannelID.UUID.String(),
		Created:       dba.Created,
		OrbTags:       fleet.Tags(dba.OrbTags),
		AgentTags:     fleet.Tags(dba.AgentTags),
		AgentMetadata: fleet.Metadata(dba.AgentMetadata),
		State:         dba.State,
		LastHBData:    fleet.Metadata(dba.LastHBData),
		LastHB:        dba.LastHB,
	}

	return agent, nil

}
