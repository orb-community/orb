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
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
	"time"
)

var _ fleet.AgentGroupRepository = (*agentGroupRepository)(nil)

type agentGroupRepository struct {
	db     Database
	logger *zap.Logger
}

func (a agentGroupRepository) Delete(ctx context.Context, groupID string, ownerID string) error {
	params := map[string]interface{}{
		"id":    groupID,
		"owner": ownerID,
	}

	q := `DELETE FROM agent_groups WHERE id = :id AND mf_owner_id = :owner;`
	if _, err := a.db.NamedQueryContext(ctx, q, params); err != nil {
		return errors.Wrap(fleet.ErrRemoveEntity, err)
	}
	return nil
}

func (a agentGroupRepository) RetrieveAllAgentGroupsByOwner(ctx context.Context, ownerID string, pm fleet.PageMetadata) (fleet.PageAgentGroup, error) {
	nameQuery, name := getNameQuery(pm.Name)
	orderQuery := getAgentGroupOrderQuery(pm.Order)
	dirQuery := getDirQuery(pm.Dir)
	metadata, metadataQuery, err := getMetadataQuery(pm.Metadata)
	if err != nil {
		return fleet.PageAgentGroup{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	tags, tagsQuery, err := getTagsQuery(pm.Tags)
	if err != nil {
		return fleet.PageAgentGroup{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	q := fmt.Sprintf(
		`select
			id,
			name,
			description,
			mf_owner_id,
			mf_channel_id,
			tags,
			ts_created,
			json_build_object('total', total, 'online', online) AS matching_agents
		from
			(select
				ag.id,
				ag.name,
				ag.description,
				ag.mf_owner_id,
				ag.mf_channel_id,
				ag.tags,
				ag.ts_created,
				sum(case when agm.agent_groups_id is not null then 1 else 0 end) as total,
				sum(case when agm.agent_state = 'online' then 1 else 0 end) as online
			from agent_groups ag
				left join agent_group_membership agm
					on ag.id = agm.agent_groups_id
					and ag.mf_owner_id = agm.mf_owner_id
			WHERE ag.mf_owner_id = :mf_owner_id %s%s%s
				group by ag.id,
					ag.name,
					ag.description,
					ag.mf_owner_id,
					ag.mf_channel_id,
					ag.tags,
					ag.ts_created)
			as agent_groups ORDER BY %s %s LIMIT :limit OFFSET :offset;`, nameQuery, tagsQuery, metadataQuery, orderQuery, dirQuery)

	params := map[string]interface{}{
		"mf_owner_id": ownerID,
		"limit":       pm.Limit,
		"offset":      pm.Offset,
		"name":        name,
		"metadata":    metadata,
		"tags":        tags,
	}
	rows, err := a.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return fleet.PageAgentGroup{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []fleet.AgentGroup
	for rows.Next() {
		dbAgentGroup := dbAgentGroup{MFOwnerID: ownerID}
		if err := rows.StructScan(&dbAgentGroup); err != nil {
			return fleet.PageAgentGroup{}, errors.Wrap(errors.ErrSelectEntity, err)
		}

		agentGroup, err := toAgentGroup(dbAgentGroup)
		if err != nil {
			return fleet.PageAgentGroup{}, errors.Wrap(errors.ErrSelectEntity, err)
		}

		items = append(items, agentGroup)
	}

	count := fmt.Sprintf(
		`select
			COUNT(*)
		from
		(select
			ag.id,
			ag.name,
			ag.description,
			ag.mf_owner_id,
			ag.mf_channel_id,
			ag.tags,
			ag.ts_created,
			sum(case when agm.agent_groups_id is not null then 1 else 0 end) as total,
			sum(case when agm.agent_state = 'online' then 1 else 0 end) as online
		from agent_groups ag
			left join agent_group_membership agm
				on ag.id = agm.agent_groups_id
					and ag.mf_owner_id = agm.mf_owner_id
		WHERE ag.mf_owner_id = :mf_owner_id %s%s%s
		group by ag.id,
			ag.name,
			ag.description,
			ag.mf_owner_id,
			ag.mf_channel_id,
			ag.tags,
			ag.ts_created) 
		as agent_groups;`, nameQuery, tagsQuery, metadataQuery)

	total, err := total(ctx, a.db, count, params)
	if err != nil {
		return fleet.PageAgentGroup{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	page := fleet.PageAgentGroup{
		AgentGroups: items,
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

func (a agentGroupRepository) RetrieveByID(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	q :=
		`select
		id,
		name,
		description,
		mf_owner_id,
		mf_channel_id,
		tags,
		ts_created,
		json_build_object('total', total, 'online', online) AS matching_agents
	from
	(select
		ag.id,
		ag.name,
		ag.description,
		ag.mf_owner_id,
		ag.mf_channel_id,
		ag.tags,
		ag.ts_created,
		sum(case when agm.agent_groups_id is not null then 1 else 0 end) as total,
		sum(case when agm.agent_state = 'online' then 1 else 0 end) as online
	from agent_groups ag
	left join agent_group_membership agm
		on ag.id = agm.agent_groups_id
		and ag.mf_owner_id = agm.mf_owner_id
	WHERE ag.id = $1 AND ag.mf_owner_id = $2
	group by ag.id,
		ag.name,
		ag.description,
		ag.mf_owner_id,
		ag.mf_channel_id,
		ag.tags,
		ag.ts_created) as agent_groups`

	if groupID == "" || ownerID == "" {
		return fleet.AgentGroup{}, errors.ErrMalformedEntity
	}

	var group dbAgentGroup
	if err := a.db.QueryRowxContext(ctx, q, groupID, ownerID).StructScan(&group); err != nil {
		if err == sql.ErrNoRows {
			return fleet.AgentGroup{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return fleet.AgentGroup{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	return toAgentGroup(group)
}

func (a agentGroupRepository) Update(ctx context.Context, ownerID string, group fleet.AgentGroup) (fleet.AgentGroup, error) {
	q := `UPDATE agent_groups SET name = :name, description = :description, tags = :tags WHERE mf_owner_id = :mf_owner_id AND id = :id;`
	groupDB, err := toDBAgentGroup(group)
	if err != nil {
		return fleet.AgentGroup{}, errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	groupDB.MFOwnerID = ownerID

	res, err := a.db.NamedExecContext(ctx, q, groupDB)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return fleet.AgentGroup{}, errors.Wrap(fleet.ErrMalformedEntity, err)
			}
		}
		return fleet.AgentGroup{}, errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fleet.AgentGroup{}, errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	if count == 0 {
		return fleet.AgentGroup{}, fleet.ErrNotFound
	}

	agMatchs, err := a.RetrieveByID(ctx, group.ID, ownerID)
	if err != nil {
		return fleet.AgentGroup{}, errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	return agMatchs, nil
}

func (a agentGroupRepository) RetrieveAllByAgent(ctx context.Context, ag fleet.Agent) ([]fleet.AgentGroup, error) {

	q := `SELECT agent_groups_id AS id, agent_groups_name AS name, group_mf_channel_id AS mf_channel_id, mf_owner_id FROM agent_group_membership WHERE agent_mf_thing_id = :agent_id`

	if ag.MFThingID == "" {
		return nil, errors.ErrMalformedEntity
	}

	params := map[string]interface{}{
		"agent_id": ag.MFThingID,
	}

	rows, err := a.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return nil, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []fleet.AgentGroup
	for rows.Next() {
		dbth := dbAgentGroup{}
		if err := rows.StructScan(&dbth); err != nil {
			return nil, errors.Wrap(errors.ErrSelectEntity, err)
		}

		th, err := toAgentGroup(dbth)
		if err != nil {
			return nil, errors.Wrap(errors.ErrViewEntity, err)
		}

		items = append(items, th)
	}

	return items, nil
}

func (a agentGroupRepository) Save(ctx context.Context, group fleet.AgentGroup) (string, error) {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	q := `INSERT INTO agent_groups (name, description, mf_owner_id, mf_channel_id, tags)         
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	if !group.Name.IsValid() || group.MFOwnerID == "" || group.MFChannelID == "" {
		return "", errors.ErrMalformedEntity
	}

	dba, err := toDBAgentGroup(group)
	if err != nil {
		return "", errors.Wrap(db.ErrSaveDB, err)
	}

	row, err := tx.QueryContext(ctx, q, dba.Name, dba.Description, dba.MFOwnerID, dba.MFChannelID, dba.Tags)
	if err != nil {
		tx.Rollback()
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

	if err = tx.Commit(); err != nil {
		return "", errors.Wrap(db.ErrSaveDB, err)
	}

	return id, nil
}

func (a agentGroupRepository) RetrieveMatchingGroups(ctx context.Context, ownerID string, thingID string) (fleet.MatchingGroups, error) {
	q := `select agent_groups_id as group_id, agent_groups_name as group_name from agent_group_membership where agent_mf_thing_id = :mf_thing_id and mf_owner_id = :mf_owner_id`

	params := map[string]interface{}{
		"mf_thing_id": thingID,
		"mf_owner_id": ownerID,
	}

	rows, err := a.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return fleet.MatchingGroups{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var groups []fleet.Group
	for rows.Next() {
		db := dbMatchingGroups{}
		if err := rows.StructScan(&db); err != nil {
			return fleet.MatchingGroups{}, errors.Wrap(errors.ErrSelectEntity, err)
		}
		groups = append(groups, fleet.Group{
			GroupID:   db.GroupID,
			GroupName: db.GroupName,
		})
	}

	return fleet.MatchingGroups{OwnerID: ownerID, Groups: groups}, nil
}

type dbAgentGroup struct {
	ID             string           `db:"id"`
	Name           types.Identifier `db:"name"`
	Description    string           `db:"description"`
	MFOwnerID      string           `db:"mf_owner_id"`
	MFChannelID    string           `db:"mf_channel_id"`
	Tags           db.Tags          `db:"tags"`
	Created        time.Time        `db:"ts_created"`
	MatchingAgents db.Metadata      `db:"matching_agents"`
}

type dbMatchingGroups struct {
	GroupID   string           `db:"group_id"`
	GroupName types.Identifier `db:"group_name"`
}

func toDBAgentGroup(group fleet.AgentGroup) (dbAgentGroup, error) {

	return dbAgentGroup{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		MFOwnerID:   group.MFOwnerID,
		MFChannelID: group.MFChannelID,
		Tags:        db.Tags(group.Tags),
	}, nil

}
func toAgentGroup(dba dbAgentGroup) (fleet.AgentGroup, error) {

	return fleet.AgentGroup{
		ID:             dba.ID,
		Name:           dba.Name,
		Description:    dba.Description,
		MFOwnerID:      dba.MFOwnerID,
		MFChannelID:    dba.MFChannelID,
		Tags:           types.Tags(dba.Tags),
		MatchingAgents: types.Metadata(dba.MatchingAgents),
	}, nil

}

func getAgentGroupOrderQuery(order string) string {
	switch order {
	case "name":
		return "name"
	default:
		return "id"
	}
}

func getTagsQuery(m types.Tags) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		mq = ` AND tags @> :tags`

		b, err := json.Marshal(m)
		if err != nil {
			return nil, "", err
		}
		mb = b
	}
	return mb, mq, nil
}

func NewAgentGroupRepository(db Database, logger *zap.Logger) fleet.AgentGroupRepository {
	return &agentGroupRepository{db: db, logger: logger}
}
