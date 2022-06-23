package migration

import (
	"context"
	"github.com/ns1labs/orb/migrate/postgres"
	"go.uber.org/zap"
	"time"
)

type M1KetoPolicies struct {
	log *zap.Logger
	dbs map[string]postgres.Database
}

type dbUser struct {
	ID    string `db:"id"`
	Email string `db:"email"`
}

type dbThingOrChannel struct {
	ID    string `db:"id"`
	Owner string `db:"owner"`
}

func NewM1KetoPolicies(log *zap.Logger, dbs map[string]postgres.Database) Plan {
	return &M1KetoPolicies{log, dbs}
}

func (m *M1KetoPolicies) Up() error {
	ctx := context.Background()
	currentTime := time.Now()

	err := m.createUserRelations(ctx, currentTime)
	if err != nil {
		return err
	}

	err = m.createChannelsOrThingsRelations(ctx, currentTime, "things")
	if err != nil {
		return err
	}

	err = m.createChannelsOrThingsRelations(ctx, currentTime, "channels")
	if err != nil {
		return err
	}

	return nil
}

func (m *M1KetoPolicies) Down() error {
	m.log.Info("do nothing for now")
	return nil
}

func (m *M1KetoPolicies) createUserRelations(ctx context.Context, currentTime time.Time) error {
	dbu := dbUser{}
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
	return nil
}

func (m *M1KetoPolicies) createChannelsOrThingsRelations(ctx context.Context, currentTime time.Time, table string) error {
	dbt := dbThingOrChannel{}
	q := "SELECT id, owner FROM " + table + ";"
	rows, err := m.dbs[postgres.DbThings].NamedQueryContext(ctx, q, map[string]interface{}{})
	if err != nil {
		return err
	}
	defer rows.Close()

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

		dbu := dbUser{}
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
	return nil
}
