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
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
	"go.uber.org/zap"
)

var _ policies.Repository = (*policiesRepository)(nil)

type policiesRepository struct {
	db     Database
	logger *zap.Logger
}

func (r policiesRepository) DeletePolicy(ctx context.Context, ownerID string, policyID string) error {
	if ownerID == "" || policyID == "" {
		return policies.ErrMalformedEntity
	}

	dbsk := dbPolicy{
		ID:        policyID,
		MFOwnerID: ownerID,
	}

	q := `DELETE FROM agent_policies WHERE id = :id AND mf_owner_id = :mf_owner_id;`
	if _, err := r.db.NamedExecContext(ctx, q, dbsk); err != nil {
		return errors.Wrap(policies.ErrRemoveEntity, err)
	}

	return nil
}

func (r policiesRepository) UpdatePolicy(ctx context.Context, owner string, plcy policies.Policy) error {
	q := `UPDATE agent_policies SET name = :name, description = :description, orb_tags = :orb_tags, policy = :policy, version = :version, ts_last_modified = CURRENT_TIMESTAMP, policy_data = :policy_data WHERE mf_owner_id = :mf_owner_id AND id = :id;`
	plcyDB, err := toDBPolicy(plcy)
	if err != nil {
		return errors.Wrap(policies.ErrUpdateEntity, err)
	}

	plcyDB.MFOwnerID = owner

	res, err := r.db.NamedExecContext(ctx, q, plcyDB)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(policies.ErrMalformedEntity, err)
			}
		}
		return errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	if count == 0 {
		return policies.ErrNotFound
	}

	return nil
}

func (r policiesRepository) RetrieveAll(ctx context.Context, owner string, pm policies.PageMetadata) (policies.Page, error) {
	nameQuery, name := getNameQuery(pm.Name)
	orderQuery := getOrderQuery(pm.Order)
	dirQuery := getDirQuery(pm.Dir)
	tags, tagsQuery, err := getTagsQuery(pm.Tags)
	if err != nil {
		return policies.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	q := fmt.Sprintf(`SELECT id, name, description, mf_owner_id, orb_tags, backend, version, policy, ts_created, ts_last_modified 
			FROM agent_policies
			WHERE mf_owner_id = :mf_owner_id %s%s ORDER BY %s %s LIMIT :limit OFFSET :offset;`, nameQuery, tagsQuery, orderQuery, dirQuery)

	params := map[string]interface{}{
		"mf_owner_id": owner,
		"limit":       pm.Limit,
		"offset":      pm.Offset,
		"name":        name,
		"tags":        tags,
	}
	rows, err := r.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return policies.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []policies.Policy
	for rows.Next() {
		dbPolicy := dbPolicy{MFOwnerID: owner}
		if err := rows.StructScan(&dbPolicy); err != nil {
			return policies.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
		}
		policy := toPolicy(dbPolicy)
		items = append(items, policy)
	}

	count := fmt.Sprintf(`SELECT count(*)
			FROM agent_policies
			WHERE mf_owner_id = :mf_owner_id %s%s;`, nameQuery, tagsQuery)

	total, err := total(ctx, r.db, count, params)
	if err != nil {
		return policies.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	page := policies.Page{
		Policies: items,
		PageMetadata: policies.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
			Order:  pm.Order,
			Dir:    pm.Dir,
		},
	}

	return page, nil
}

func (r policiesRepository) RetrievePoliciesByGroupID(ctx context.Context, groupIDs []string, ownerID string) ([]policies.PolicyInDataset, error) {

	q := `SELECT agent_policies.id AS id, datasets.id AS dataset_id, agent_policies.name AS name, agent_policies.mf_owner_id, orb_tags, backend, version, policy, agent_policies.ts_created 
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

	var items []policies.PolicyInDataset
	for rows.Next() {
		dbth := dbPolicy{MFOwnerID: ownerID}
		if err := rows.StructScan(&dbth); err != nil {
			return nil, errors.Wrap(errors.ErrSelectEntity, err)
		}

		th := toPolicy(dbth)
		items = append(items, policies.PolicyInDataset{Policy: th, DatasetID: dbth.DataSetID})
	}

	return items, nil
}

func (r policiesRepository) RetrievePolicyByID(ctx context.Context, policyID string, ownerID string) (policies.Policy, error) {
	q := `SELECT id, name, description, mf_owner_id, orb_tags, backend, version, policy, ts_created, schema_version, ts_last_modified, policy_data, format 
			FROM agent_policies WHERE id = $1 AND mf_owner_id = $2`

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

func (r policiesRepository) UpdateDataset(ctx context.Context, ownerID string, ds policies.Dataset) error {
	q := `UPDATE datasets SET tags = :tags, sink_ids = :sink_ids, name = :name WHERE mf_owner_id = :mf_owner_id AND id = :id;`

	params := map[string]interface{}{
		"mf_owner_id": ds.MFOwnerID,
		"tags":        db.Tags(ds.Tags),
		"sink_ids":    pq.Array(ds.SinkIDs),
		"id":          ds.ID,
		"name":        ds.Name,
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
		return errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	if count == 0 {
		return policies.ErrNotFound
	}

	return nil
}

func (r policiesRepository) DeleteDataset(ctx context.Context, ownerID string, dsID string) error {
	q := `DELETE FROM datasets WHERE mf_owner_id = :mf_owner_id AND id = :id;`

	params := map[string]interface{}{
		"mf_owner_id": ownerID,
		"id":          dsID,
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
		return errors.Wrap(fleet.ErrUpdateEntity, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(fleet.ErrRemoveEntity, err)
	}

	if count == 0 {
		return policies.ErrNotFound
	}

	return nil
}

func (r policiesRepository) SaveDataset(ctx context.Context, dataset policies.Dataset) (string, error) {

	q := `INSERT INTO datasets (name, mf_owner_id, metadata, valid, agent_group_id, agent_policy_id, sink_ids, tags)         
			  VALUES (:name, :mf_owner_id, :metadata, :valid, :agent_group_id, :agent_policy_id, :sink_ids_str, :tags) RETURNING id`

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

func (r policiesRepository) InactivateDatasetByPolicyID(ctx context.Context, policyID string, ownerID string) error {
	q := `UPDATE datasets SET valid = false WHERE mf_owner_id = :mf_owner_id and agent_policy_id = :agent_policy_id`

	params := map[string]interface{}{
		"agent_policy_id": policyID,
		"mf_owner_id":     ownerID,
	}

	_, err := r.db.NamedExecContext(ctx, q, params)
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
	return nil
}

func (r policiesRepository) SavePolicy(ctx context.Context, policy policies.Policy) (string, error) {

	q := `INSERT INTO agent_policies (name, mf_owner_id, backend, schema_version, policy, orb_tags, description, policy_data, format)         
			  VALUES (:name, :mf_owner_id, :backend, :schema_version, :policy, :orb_tags, :description, :policy_data, :format) RETURNING id`

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

func (r policiesRepository) RetrieveDatasetsByPolicyID(ctx context.Context, policyID string, ownerID string) ([]policies.Dataset, error) {

	q := `SELECT id, name, mf_owner_id, valid, agent_group_id, agent_policy_id, sink_ids, metadata, ts_created 
			FROM datasets
			WHERE agent_policy_id = ? AND mf_owner_id = ?`

	if policyID == "" || ownerID == "" {
		return nil, errors.ErrMalformedEntity
	}

	query, args, err := sqlx.In(q, policyID, ownerID)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []policies.Dataset
	for rows.Next() {
		dbth := dbDataset{MFOwnerID: ownerID}
		if err := rows.StructScan(&dbth); err != nil {
			return nil, errors.Wrap(errors.ErrSelectEntity, err)
		}

		th := toDataset(dbth)
		items = append(items, th)
	}

	return items, nil
}

func (r policiesRepository) RetrieveDatasetByID(ctx context.Context, datasetID string, ownerID string) (policies.Dataset, error) {
	q := `SELECT id, name, mf_owner_id, valid, agent_group_id, agent_policy_id, sink_ids, metadata, ts_created FROM datasets WHERE id = $1 AND mf_owner_id = $2`

	if datasetID == "" || ownerID == "" {
		return policies.Dataset{}, errors.ErrMalformedEntity
	}

	var dba dbDataset
	if err := r.db.QueryRowxContext(ctx, q, datasetID, ownerID).StructScan(&dba); err != nil {
		if err == sql.ErrNoRows {
			return policies.Dataset{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return policies.Dataset{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	return toDataset(dba), nil
}

func (r policiesRepository) RetrieveAllDatasetsByOwner(ctx context.Context, owner string, pm policies.PageMetadata) (policies.PageDataset, error) {
	nameQuery, name := getNameQuery(pm.Name)
	orderQuery := getOrderQuery(pm.Order)
	dirQuery := getDirQuery(pm.Dir)

	q := fmt.Sprintf(`SELECT id, name, mf_owner_id, valid, agent_group_id, agent_policy_id, sink_ids, metadata, ts_created 
			FROM datasets
			WHERE mf_owner_id = :mf_owner_id %s ORDER BY %s %s LIMIT :limit OFFSET :offset;`, nameQuery, orderQuery, dirQuery)

	params := map[string]interface{}{
		"mf_owner_id": owner,
		"limit":       pm.Limit,
		"offset":      pm.Offset,
		"name":        name,
	}
	rows, err := r.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return policies.PageDataset{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []policies.Dataset
	for rows.Next() {
		dbDataset := dbDataset{MFOwnerID: owner}
		if err := rows.StructScan(&dbDataset); err != nil {
			return policies.PageDataset{}, errors.Wrap(errors.ErrSelectEntity, err)
		}
		dataset := toDataset(dbDataset)
		items = append(items, dataset)
	}

	count := fmt.Sprintf(`SELECT count(*)
			FROM datasets
			WHERE mf_owner_id = :mf_owner_id %s;`, nameQuery)

	total, err := total(ctx, r.db, count, params)
	if err != nil {
		return policies.PageDataset{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	pageDataset := policies.PageDataset{
		Datasets: items,
		PageMetadata: policies.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
			Order:  pm.Order,
			Dir:    pm.Dir,
		},
	}

	return pageDataset, nil
}

func (r policiesRepository) InactivateDatasetByID(ctx context.Context, id string, ownerID string) error {
	q := `UPDATE datasets SET valid = false WHERE mf_owner_id = :mf_owner_id AND :id = id`

	params := map[string]interface{}{
		"mf_owner_id": ownerID,
		"id":          id,
	}

	_, err := r.db.NamedExecContext(ctx, q, params)
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

	return nil
}

func (r policiesRepository) DeleteSinkFromAllDatasets(ctx context.Context, sinkID string, ownerID string) ([]policies.Dataset, error) {
	q := `UPDATE datasets SET sink_ids = array_remove(sink_ids, :sink_ids) WHERE mf_owner_id = :mf_owner_id RETURNING *`

	if ownerID == "" {
		return []policies.Dataset{}, errors.ErrMalformedEntity
	}

	params := map[string]interface{}{
		"mf_owner_id": ownerID,
		"sink_ids":    sinkID,
	}

	res, err := r.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return []policies.Dataset{}, errors.Wrap(policies.ErrMalformedEntity, err)
			}
		}
		return []policies.Dataset{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	defer res.Close()

	var datasets []policies.Dataset
	for res.Next() {
		dbDataset := dbDataset{MFOwnerID: ownerID}
		if err := res.StructScan(&dbDataset); err != nil {
			return []policies.Dataset{}, errors.Wrap(errors.ErrSelectEntity, err)
		}
		dataset := toDataset(dbDataset)
		datasets = append(datasets, dataset)
	}

	return datasets, nil
}

func (r policiesRepository) ActivateDatasetByID(ctx context.Context, id string, ownerID string) error {
	q := `UPDATE datasets SET valid = true WHERE mf_owner_id = :mf_owner_id AND :id = id`

	if ownerID == "" || id == "" {
		return errors.ErrMalformedEntity
	}

	params := map[string]interface{}{
		"mf_owner_id": ownerID,
		"id":          id,
	}

	_, err := r.db.NamedExecContext(ctx, q, params)
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

	return nil
}

func (r policiesRepository) DeleteAgentGroupFromAllDatasets(ctx context.Context, groupID string, ownerID string) error {
	q := `UPDATE datasets SET agent_group_id = null WHERE mf_owner_id = :mf_owner_id AND agent_group_id = :agent_group_id`

	if ownerID == "" {
		return errors.ErrMalformedEntity
	}

	params := map[string]interface{}{
		"mf_owner_id":    ownerID,
		"agent_group_id": groupID,
	}

	res, err := r.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(policies.ErrMalformedEntity, err)
			}
		}
		return errors.Wrap(errors.ErrSelectEntity, err)
	}

	defer res.Close()

	return nil
}

type dbPolicy struct {
	ID            string           `db:"id"`
	Name          types.Identifier `db:"name"`
	MFOwnerID     string           `db:"mf_owner_id"`
	Backend       string           `db:"backend"`
	SchemaVersion string           `db:"schema_version"`
	Description   string           `db:"description"`
	OrbTags       db.Tags          `db:"orb_tags"`
	Policy        db.Metadata      `db:"policy"`
	PolicyData    string           `db:"policy_data"`
	Format        string           `db:"format"`
	Version       int32            `db:"version"`
	Created       time.Time        `db:"ts_created"`
	DataSetID     string           `db:"dataset_id"`
	LastModified  time.Time        `db:"ts_last_modified"`
}

func toDBPolicy(policy policies.Policy) (dbPolicy, error) {

	var uID uuid.UUID
	err := uID.Scan(policy.MFOwnerID)
	if err != nil {
		return dbPolicy{}, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return dbPolicy{
		ID:            policy.ID,
		Name:          policy.Name,
		Description:   policy.Description,
		Version:       policy.Version,
		MFOwnerID:     uID.String(),
		Backend:       policy.Backend,
		SchemaVersion: policy.SchemaVersion,
		OrbTags:       db.Tags(policy.OrbTags),
		Policy:        db.Metadata(policy.Policy),
		PolicyData:    policy.PolicyData,
		Format:        policy.Format,
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
	TsCreated    time.Time        `db:"ts_created"`
	Tags         db.Tags          `db:"tags"`
	SinkIDs      pq.StringArray   `db:"sink_ids"`
	SinksIDsStr  interface{}      `db:"sink_ids_str"`
}

func toDBDataset(dataset policies.Dataset) (dbDataset, error) {

	var uID uuid.UUID
	err := uID.Scan(dataset.MFOwnerID)
	if err != nil {
		return dbDataset{}, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	d := dbDataset{
		ID:          dataset.ID,
		Name:        dataset.Name,
		MFOwnerID:   uID.String(),
		Metadata:    db.Metadata(dataset.Metadata),
		Tags:        db.Tags(dataset.Tags),
		SinksIDsStr: pq.Array(dataset.SinkIDs),
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

	return d, nil

}

func NewPoliciesRepository(db Database, log *zap.Logger) policies.Repository {
	return &policiesRepository{db: db, logger: log}
}

func toPolicy(dba dbPolicy) policies.Policy {

	policy := policies.Policy{
		ID:            dba.ID,
		Name:          dba.Name,
		Description:   dba.Description,
		MFOwnerID:     dba.MFOwnerID,
		Backend:       dba.Backend,
		SchemaVersion: dba.SchemaVersion,
		Version:       dba.Version,
		OrbTags:       types.Tags(dba.OrbTags),
		Policy:        types.Metadata(dba.Policy),
		Created:       dba.Created,
		LastModified:  dba.LastModified,
		PolicyData:    dba.PolicyData,
		Format:        dba.Format,
	}

	return policy

}

func toDataset(dba dbDataset) policies.Dataset {
	dataset := policies.Dataset{
		ID:           dba.ID,
		Name:         dba.Name,
		MFOwnerID:    dba.MFOwnerID,
		Valid:        dba.Valid,
		AgentGroupID: dba.AgentGroupID.String,
		PolicyID:     dba.PolicyID.String,
		SinkIDs:      dba.SinkIDs,
		Metadata:     types.Metadata(dba.Metadata),
		Created:      dba.TsCreated,
		Tags:         types.Tags(dba.Tags),
	}

	return dataset
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
		return "id"
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

func getTagsQuery(m types.Tags) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		mq = ` AND orb_tags @> :tags`

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
