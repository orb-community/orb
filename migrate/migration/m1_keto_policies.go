package migration

import (
	"context"
	"github.com/ns1labs/orb/migrate/postgres"
	"github.com/ns1labs/orb/pkg/db"
	"go.uber.org/zap"
	"time"
)

type M1KetoPolicies struct {
	log *zap.Logger
	dbs map[string]postgres.Database
}

type dbUser struct {
	ID       string      `db:"id"`
	Email    string      `db:"email"`
	Metadata db.Metadata `db:"metadata"`
}

type dbThing struct {
	ID       string      `db:"id"`
	Owner    string      `db:"owner"`
	Key      string      `db:"key"`
	Name     string      `db:"name"`
	Metadata db.Metadata `db:"metadata"`
}

type dbChannel struct {
	ID       string      `db:"id"`
	Owner    string      `db:"owner"`
	Name     string      `db:"name"`
	Metadata db.Metadata `db:"metadata"`
}

func NewM1KetoPolicies(log *zap.Logger, dbs map[string]postgres.Database) Plan {
	return &M1KetoPolicies{log, dbs}
}

func (m *M1KetoPolicies) Up() error {
	ctx := context.Background()
	dbu := dbUser{}
	currentTime := time.Now()

	q := "SELECT id FROM users;"
	rows, err := m.dbs[postgres.DbUsers].NamedQueryContext(ctx, q, map[string]interface{}{})
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(&dbu); err != nil {
			return err
		}

		// create entry putting user as member
		q = `INSERT INTO keto_0000000000_relation_tuples (shard_id, object, relation, subject, commit_time)
			  VALUES (:shard_id, :object, :relation, :subject, :commit_time);`
		params := map[string]interface{}{
			"shard_id":    "default",
			"object":      "users",
			"relation":    "member",
			"subject":     dbu.ID,
			"commit_time": currentTime,
		}

		_, err = m.dbs[postgres.DbKeto].NamedExecContext(ctx, q, params)
		if err != nil {
			return err
		}
	}

	dbt := dbThing{}
	q = "SELECT id, owner FROM things;"
	rows, err = m.dbs[postgres.DbThings].NamedQueryContext(ctx, q, map[string]interface{}{})
	if err != nil {
		return err
	}

	for rows.Next() {
		if err := rows.StructScan(&dbt); err != nil {
			return err
		}

		// create read, write and create entries for authorities
		for _, relation := range []string{0: "read", 1: "write", 2: "creation"} {
			q = `INSERT INTO keto_0000000000_relation_tuples (shard_id, object, relation, subject, commit_time)
				VALUES (:shard_id, :object, :relation, :subject, :commit_time);`
			params := map[string]interface{}{
				"shard_id":    "default",
				"object":      dbt.ID,
				"relation":    relation,
				"subject":     "members:authorities#member",
				"commit_time": currentTime,
			}

			_, err = m.dbs[postgres.DbKeto].NamedExecContext(ctx, q, params)
			if err != nil {
				return err
			}
		}

		q := `SELECT id FROM users where email = $1;`

		if err := m.dbs[postgres.DbUsers].QueryRowxContext(ctx, q, dbt.Owner).StructScan(&dbu); err != nil {
			return err
		}

		// create read, write and create entries for owner
		for _, relation := range []string{0: "read", 1: "write", 2: "creation"} {
			q = `INSERT INTO keto_0000000000_relation_tuples (shard_id, object, relation, subject, commit_time)
				VALUES (:shard_id, :object, :relation, :subject, :commit_time);`
			params := map[string]interface{}{
				"shard_id":    "default",
				"object":      dbt.ID,
				"relation":    relation,
				"subject":     dbu.ID,
				"commit_time": currentTime,
			}

			_, err = m.dbs[postgres.DbKeto].NamedExecContext(ctx, q, params)
			if err != nil {
				return err
			}
		}
	}

	dbc := dbChannel{}
	q = "SELECT id, owner FROM channels;"
	rows, err = m.dbs[postgres.DbThings].NamedQueryContext(ctx, q, map[string]interface{}{})
	if err != nil {
		return err
	}

	for rows.Next() {
		if err := rows.StructScan(&dbc); err != nil {
			return err
		}

		// create read, write and create entries for authorities
		for _, relation := range []string{0: "read", 1: "write", 2: "creation"} {
			q = `INSERT INTO keto_0000000000_relation_tuples (shard_id, object, relation, subject, commit_time)
				VALUES (:shard_id, :object, :relation, :subject, :commit_time);`
			params := map[string]interface{}{
				"shard_id":    "default",
				"object":      dbc.ID,
				"relation":    relation,
				"subject":     "members:authorities#member",
				"commit_time": currentTime,
			}

			_, err = m.dbs[postgres.DbKeto].NamedExecContext(ctx, q, params)
			if err != nil {
				return err
			}
		}

		// get owner user searching by its email
		q := `SELECT id FROM users where email = $1;`

		if err := m.dbs[postgres.DbUsers].QueryRowxContext(ctx, q, dbc.Owner).StructScan(&dbu); err != nil {
			return err
		}

		// create read, write and create entries for owner
		for _, relation := range []string{0: "read", 1: "write", 2: "creation"} {
			q = `INSERT INTO keto_0000000000_relation_tuples (shard_id, object, relation, subject, commit_time)
				VALUES (:shard_id, :object, :relation, :subject, :commit_time);`
			params := map[string]interface{}{
				"shard_id":    "default",
				"object":      dbc.ID,
				"relation":    relation,
				"subject":     dbu.ID,
				"commit_time": currentTime,
			}

			_, err = m.dbs[postgres.DbKeto].NamedExecContext(ctx, q, params)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *M1KetoPolicies) Down() error {
	m.log.Info("do nothing for now")
	return nil
}
