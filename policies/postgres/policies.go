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
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
	"go.uber.org/zap"
	"time"
)

var _ policies.Repository = (*policiesRepository)(nil)

type policiesRepository struct {
	db     Database
	logger *zap.Logger
}

func (r policiesRepository) RetrieveAll(ctx context.Context, owner string, pm policies.PageMetadata) (policies.Page, error) {
	return policies.Page{}, nil
}

func (r policiesRepository) RetrievePoliciesByGroupID(ctx context.Context, groupIDs []string, ownerID string) ([]policies.Policy, error) {

	q := `SELECT agent_policies.id AS id, agent_policies.name AS name, agent_policies.mf_owner_id, orb_tags, backend, version, policy, agent_policies.ts_created 
			FROM agent_policies, datasets
			WHERE agent_policies.id = datasets.agent_policy_id AND agent_policies.mf_owner_id = datasets.mf_owner_id AND valid = TRUE AND
				agent_group_id IN (?) AND agent_policies.mf_owner_id = ?`

	if len(groupIDs) == 0 || ownerID == "" {
		return nil, errors.ErrMalformedEntity
	}

	query, args, err := sqlx.In(q, groupIDs, ownerID)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []policies.Policy
	for rows.Next() {
		dbth := dbPolicy{MFOwnerID: ownerID}
		if err := rows.StructScan(&dbth); err != nil {
			return nil, errors.Wrap(errors.ErrSelectEntity, err)
		}

		th := toPolicy(dbth)
		items = append(items, th)
	}

	return items, nil
}

func (r policiesRepository) RetrievePolicyByID(ctx context.Context, policyID string, ownerID string) (policies.Policy, error) {
	q := `SELECT id, name, mf_owner_id, orb_tags, backend, version, policy, ts_created FROM agent_policies WHERE id = $1 AND mf_owner_id = $2`

	if policyID == "" || ownerID == "" {
		return policies.Policy{}, errors.ErrMalformedEntity
	}

	var dbp dbPolicy
	if err := r.db.QueryRowxContext(ctx, q, policyID, ownerID).StructScan(&dbp); err != nil {
		if err == sql.ErrNoRows {
			return policies.Policy{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return policies.Policy{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	return toPolicy(dbp), nil
}

func (r policiesRepository) SaveDataset(ctx context.Context, dataset policies.Dataset) (string, error) {

	q := `INSERT INTO datasets (name, mf_owner_id, metadata, valid, agent_group_id, agent_policy_id, sink_id)         
			  VALUES (:name, :mf_owner_id, :metadata, :valid, :agent_group_id, :agent_policy_id, :sink_id) RETURNING id`

	if !dataset.Name.IsValid() || dataset.MFOwnerID == "" {
		return "", errors.ErrMalformedEntity
	}

	dba, err := toDBDataset(dataset)
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

func (r policiesRepository) InactivateDatasetByGroupID(ctx context.Context, groupID string, ownerID string) error {
	q := `UPDATE datasets SET valid = false WHERE mf_owner_id = :mf_owner_id and agent_group_id = :agent_group_id`

	params := map[string]interface{}{
		"agent_group_id": groupID,
		"mf_owner_id":    ownerID,
	}

	res, err := r.db.NamedExecContext(ctx, q, params)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(policies.ErrMalformedEntity, err)
			}
		}
		return errors.Wrap(policies.ErrUpdateEntity, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(policies.ErrUpdateEntity, err)
	}

	if count == 0 {
		return policies.ErrInactivateDataset
	}
	return nil
}

func (r policiesRepository) SavePolicy(ctx context.Context, policy policies.Policy) (string, error) {

	q := `INSERT INTO agent_policies (name, mf_owner_id, backend, policy, orb_tags)         
			  VALUES (:name, :mf_owner_id, :backend, :policy, :orb_tags) RETURNING id`

	if !policy.Name.IsValid() || policy.MFOwnerID == "" {
		return "", errors.ErrMalformedEntity
	}

	dba, err := toDBPolicy(policy)
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

type dbPolicy struct {
	ID        string           `db:"id"`
	Name      types.Identifier `db:"name"`
	MFOwnerID string           `db:"mf_owner_id"`
	Backend   string           `db:"backend"`
	Format    string           `db:"format"`
	OrbTags   db.Tags          `db:"orb_tags"`
	Policy    db.Metadata      `db:"policy"`
	Version   int32            `db:"version"`
	Created   time.Time        `db:"ts_created"`
}

func toDBPolicy(policy policies.Policy) (dbPolicy, error) {

	var uID uuid.UUID
	err := uID.Scan(policy.MFOwnerID)
	if err != nil {
		return dbPolicy{}, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return dbPolicy{
		ID:        policy.ID,
		Name:      policy.Name,
		MFOwnerID: uID.String(),
		Backend:   policy.Backend,
		OrbTags:   db.Tags(policy.OrbTags),
		Policy:    db.Metadata(policy.Policy),
	}, nil

}

type dbDataset struct {
	ID           string           `db:"id"`
	Name         types.Identifier `db:"name"`
	MFOwnerID    string           `db:"mf_owner_id"`
	Metadata     db.Metadata      `db:"metadata"`
	Valid        bool             `db:"valid"`
	AgentGroupID sql.NullString   `db:"agent_group_id"`
	PolicyID     sql.NullString   `db:"agent_policy_id"`
	SinkID       sql.NullString   `db:"sink_id"`
}

func toDBDataset(dataset policies.Dataset) (dbDataset, error) {

	var uID uuid.UUID
	err := uID.Scan(dataset.MFOwnerID)
	if err != nil {
		return dbDataset{}, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	d := dbDataset{
		ID:        dataset.ID,
		Name:      dataset.Name,
		MFOwnerID: uID.String(),
		Metadata:  db.Metadata(dataset.Metadata),
	}

	d.Valid = true
	if dataset.AgentGroupID != "" {
		d.AgentGroupID = sql.NullString{String: dataset.AgentGroupID, Valid: true}
	} else {
		d.AgentGroupID = sql.NullString{Valid: false}
		d.Valid = false
	}
	if dataset.PolicyID != "" {
		d.PolicyID = sql.NullString{String: dataset.PolicyID, Valid: true}
	} else {
		d.PolicyID = sql.NullString{Valid: false}
		d.Valid = false
	}
	if dataset.SinkID != "" {
		d.SinkID = sql.NullString{String: dataset.SinkID, Valid: true}
	} else {
		d.SinkID = sql.NullString{Valid: false}
		d.Valid = false
	}

	return d, nil

}

func NewPoliciesRepository(db Database, log *zap.Logger) policies.Repository {
	return &policiesRepository{db: db, logger: log}
}

func toPolicy(dba dbPolicy) policies.Policy {

	policy := policies.Policy{
		ID:        dba.ID,
		Name:      dba.Name,
		MFOwnerID: dba.MFOwnerID,
		Backend:   dba.Backend,
		Version:   dba.Version,
		OrbTags:   types.Tags(dba.OrbTags),
		Policy:    types.Metadata(dba.Policy),
		Created:   dba.Created,
	}

	return policy

}
