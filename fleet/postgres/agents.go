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
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
	"strings"
	"time"
)

var _ fleet.AgentRepository = (*agentRepository)(nil)

type agentRepository struct {
	db     Database
	logger *zap.Logger
}

func (r agentRepository) RetrieveAll(ctx context.Context, owner string, pm fleet.PageMetadata) (fleet.Page, error) {
	nq, name := getNameQuery(pm.Name)
	oq := getOrderQuery(pm.Order)
	dq := getDirQuery(pm.Dir)
	m, mq, err := getMetadataQuery(pm.Metadata)
	if err != nil {
		return fleet.Page{}, errors.Wrap(fleet.ErrSelectEntity, err)
	}

	q := fmt.Sprintf(`SELECT * FROM agents
	      WHERE mf_owner_id = :mf_owner_id %s%s ORDER BY %s %s LIMIT :limit OFFSET :offset;`, mq, nq, oq, dq)
	params := map[string]interface{}{
		"mf_owner_id": owner,
		"limit":       pm.Limit,
		"offset":      pm.Offset,
		"name":        name,
		"metadata":    m,
	}

	rows, err := r.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return fleet.Page{}, errors.Wrap(fleet.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []fleet.Agent
	for rows.Next() {
		dbth := dbAgent{MFOwnerID: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return fleet.Page{}, errors.Wrap(fleet.ErrSelectEntity, err)
		}

		th, err := toAgent(dbth)
		if err != nil {
			return fleet.Page{}, errors.Wrap(fleet.ErrViewEntity, err)
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM agents WHERE mf_owner_id = :mf_owner_id %s%s;`, nq, mq)

	total, err := total(ctx, r.db, cq, params)
	if err != nil {
		return fleet.Page{}, errors.Wrap(fleet.ErrSelectEntity, err)
	}

	page := fleet.Page{
		Agents: items,
		PageMetadata: fleet.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
			Order:  pm.Order,
			Dir:    pm.Dir,
		},
	}

	return page, nil
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
			= (:last_hb_data, now(), :state) 
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

	if !agent.Name.IsValid() || agent.MFOwnerID == "" || agent.MFThingID == "" || agent.MFChannelID == "" {
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
	Name          types.Identifier `db:"name"`
	MFOwnerID     string           `db:"mf_owner_id"`
	MFThingID     string           `db:"mf_thing_id"`
	MFChannelID   string           `db:"mf_channel_id"`
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
		MFThingID:     agent.MFThingID,
		MFChannelID:   agent.MFChannelID,
		OrbTags:       dbTags(agent.OrbTags),
		AgentTags:     dbTags(agent.AgentTags),
		AgentMetadata: dbMetadata(agent.AgentMetadata),
		State:         agent.State,
		Created:       agent.Created,
		LastHBData:    dbMetadata(agent.LastHBData),
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
		MFThingID:     dba.MFThingID,
		MFChannelID:   dba.MFChannelID,
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
func getNameQuery(name string) (string, string) {
	if name == "" {
		return "", ""
	}
	name = fmt.Sprintf(`%%%s%%`, strings.ToLower(name))
	nq := ` AND LOWER(name) LIKE :name`
	return nq, name
}

func getOrderQuery(order string) string {
	switch order {
	case "name":
		return "name"
	default:
		return "mf_thing_id"
	}
}

func getDirQuery(dir string) string {
	switch dir {
	case "asc":
		return "ASC"
	default:
		return "DESC"
	}
}

func getMetadataQuery(m fleet.Metadata) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		mq = ` AND agent_metadata @> :metadata`

		b, err := json.Marshal(m)
		if err != nil {
			return nil, "", err
		}
		mb = b
	}
	return mb, mq, nil
}

func total(ctx context.Context, db Database, query string, params interface{}) (uint64, error) {
	rows, err := db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}
	return total, nil
}

func NewAgentRepository(db Database, logger *zap.Logger) fleet.AgentRepository {
	return &agentRepository{db: db, logger: logger}
}
