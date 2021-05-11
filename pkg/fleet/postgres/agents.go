// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"github.com/mainflux/mainflux/users"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/fleet"
	"go.uber.org/zap"
)

const (
	errDuplicate  = "unique_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

var (
	errSaveAgentDB = errors.New("failed to save agent to database")
	errMarshal     = errors.New("Failed to marshal metadata")
	errUnmarshal   = errors.New("Failed to unmarshal metadata")
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
		return errors.Wrap(errSaveAgentDB, err)
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
		return errors.Wrap(errSaveAgentDB, err)
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

// dbMetadata type for handling metadata properly in database/sql
type dbMetadata map[string]interface{}

// Scan - Implement the database/sql scanner interface
func (m *dbMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return users.ErrScanMetadata
	}

	if err := json.Unmarshal(b, m); err != nil {
		return err
	}

	return nil
}

// Value Implements valuer
func (m dbMetadata) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "{}", nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
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
