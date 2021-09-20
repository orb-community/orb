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
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	"go.uber.org/zap"
	"strings"
	"time"
)

var _ sinks.SinkRepository = (*sinksRepository)(nil)

type sinksRepository struct {
	db     Database
	logger *zap.Logger
}

func (s sinksRepository) Save(ctx context.Context, sink sinks.Sink) (string, error) {
	q := `INSERT INTO sinks (name, mf_owner_id, metadata, description, backend, tags, state, error)         
			  VALUES (:name, :mf_owner_id, :metadata, :description, :backend, :tags, :state, :error) RETURNING id`

	if !sink.Name.IsValid() || sink.MFOwnerID == "" {
		return "", errors.ErrMalformedEntity
	}

	dba, err := toDBSink(sink)
	if err != nil {
		return "", errors.Wrap(db.ErrSaveDB, err)
	}

	row, err := s.db.NamedQueryContext(ctx, q, dba)
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

func (s sinksRepository) Update(ctx context.Context, sink sinks.Sink) error {
	q := `UPDATE sinks SET description = :description, tags = :tags, metadata = :metadata WHERE mf_owner_id = :mf_owner_id AND id = :id;`

	sinkDB, err := toDBSink(sink)
	if err != nil {
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	res, err := s.db.NamedExecContext(ctx, q, sinkDB)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case db.ErrInvalid, db.ErrTruncation:
				return errors.Wrap(sinks.ErrMalformedEntity, err)
			}
		}
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(sinks.ErrUpdateEntity, err)
	}

	if count == 0 {
		return sinks.ErrNotFound
	}
	return nil
}

func (s sinksRepository) RetrieveAll(ctx context.Context, owner string, pm sinks.PageMetadata) (sinks.Page, error) {
	name, nameQuery := getNameQuery(pm.Name)
	orderQuery := getOrderQuery(pm.Order)
	dirQuery := getDirQuery(pm.Dir)
	metadata, metadataQuery, err := getMetadataQuery(pm.Metadata)
	if err != nil {
		return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	tags, tagsQuery, err := getTagsQuery(pm.Tags)
	if err != nil {
		return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	q := fmt.Sprintf(`SELECT id, name, mf_owner_id, description, tags, coalesce(state, '') as state, coalesce(error, '') as error, backend, metadata, ts_created
								FROM sinks WHERE mf_owner_id = :mf_owner_id %s%s%s ORDER BY %s %s LIMIT :limit OFFSET :offset;`, tagsQuery, metadataQuery, nameQuery, orderQuery, dirQuery)
	params := map[string]interface{}{
		"mf_owner_id": owner,
		"limit":       pm.Limit,
		"offset":      pm.Offset,
		"name":        name,
		"metadata":    metadata,
		"tags":        tags,
	}
	rows, err := s.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []sinks.Sink
	for rows.Next() {
		dbSink := dbSink{MFOwnerID: owner}
		if err := rows.StructScan(&dbSink); err != nil {
			return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
		}

		sink, err := toSink(dbSink)
		if err != nil {
			return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
		}

		items = append(items, sink)
	}

	count := fmt.Sprintf(`SELECT COUNT(*) FROM sinks WHERE mf_owner_id = :mf_owner_id %s%s%s`, tagsQuery, metadataQuery, nameQuery)

	total, err := total(ctx, s.db, count, params)
	if err != nil {
		return sinks.Page{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	page := sinks.Page{
		Sinks: items,
		PageMetadata: sinks.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
			Order:  pm.Order,
			Dir:    pm.Dir,
		},
	}

	return page, nil
}

func (s sinksRepository) RetrieveById(ctx context.Context, id string) (sinks.Sink, error) {

	q := `SELECT id, name, mf_owner_id, description, tags, backend, metadata, ts_created, coalesce(state, '') as state, coalesce(error, '') as error
			FROM sinks where id = $1`

	dba := dbSink{}

	if err := s.db.QueryRowxContext(ctx, q, id).StructScan(&dba); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && db.ErrInvalid == pqErr.Code.Name() {
			return sinks.Sink{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return sinks.Sink{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	return toSink(dba)
}

func (s sinksRepository) RetrieveByOwnerAndId(ctx context.Context, ownerID string, id string) (sinks.Sink, error) {

	q := `SELECT id, name, mf_owner_id, description, tags, backend, metadata, ts_created, coalesce(state, '') as state, coalesce(error, '') as error
			FROM sinks where id = $1 and mf_owner_id = $2`

	if ownerID == "" || id == "" {
		return sinks.Sink{}, errors.ErrSelectEntity
	}

	dba := dbSink{}

	if err := s.db.QueryRowxContext(ctx, q, id, ownerID).StructScan(&dba); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && db.ErrInvalid == pqErr.Code.Name() {
			return sinks.Sink{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return sinks.Sink{}, errors.Wrap(errors.ErrSelectEntity, err)
	}

	return toSink(dba)
}

func (s sinksRepository) Remove(ctx context.Context, owner, id string) error {
	dbsk := dbSink{
		ID:        id,
		MFOwnerID: owner,
	}

	q := `DELETE FROM sinks WHERE id = :id AND mf_owner_id = :mf_owner_id;`
	if _, err := s.db.NamedExecContext(ctx, q, dbsk); err != nil {
		return errors.Wrap(sinks.ErrRemoveEntity, err)
	}

	return nil
}

type dbSink struct {
	ID          string           `db:"id"`
	Name        types.Identifier `db:"name"`
	MFOwnerID   string           `db:"mf_owner_id"`
	Metadata    db.Metadata      `db:"metadata"`
	Backend     string           `db:"backend"`
	Description string           `db:"description"`
	Created     time.Time        `db:"ts_created"`
	Tags        db.Tags          `db:"tags"`
	State       string           `db:"state"`
	Error       string           `db:"error"`
}

func toDBSink(sink sinks.Sink) (dbSink, error) {

	var uID uuid.UUID
	err := uID.Scan(sink.MFOwnerID)
	if err != nil {
		return dbSink{}, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return dbSink{
		ID:          sink.ID,
		Name:        sink.Name,
		MFOwnerID:   uID.String(),
		Metadata:    db.Metadata(sink.Config),
		Backend:     sink.Backend,
		Description: sink.Description,
		State:       sink.State,
		Error:       sink.Error,
		Tags:        db.Tags(sink.Tags),
	}, nil

}

func toSink(dba dbSink) (sinks.Sink, error) {
	sink := sinks.Sink{
		ID:          dba.ID,
		Name:        dba.Name,
		MFOwnerID:   dba.MFOwnerID,
		Backend:     dba.Backend,
		Description: dba.Description,
		State:       dba.State,
		Error:       dba.Error,
		Config:      types.Metadata(dba.Metadata),
		Created:     dba.Created,
		Tags:        types.Tags(dba.Tags),
	}
	return sink, nil
}

func getNameQuery(name string) (string, string) {
	if name == "" {
		return "", ""
	}
	name = fmt.Sprintf(`%%%s%%`, strings.ToLower(name))
	nameQuey := ` AND LOWER(name) LIKE :name`
	return name, nameQuey
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

func getMetadataQuery(m types.Metadata) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		mq = ` AND metadata @> :metadata`

		b, err := json.Marshal(m)
		if err != nil {
			return nil, "", err
		}
		mb = b
	}
	return mb, mq, nil
}

func getTagsQuery(m types.Tags) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		// todo add in orb tags
		mq = ` AND tags @> :tags`

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

func NewSinksRepository(db Database, logger *zap.Logger) sinks.SinkRepository {
	return &sinksRepository{db: db, logger: logger}
}
