package migration

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/migrate/postgres"
	"go.uber.org/zap"
	"time"
)

type M1KetoPolicies struct {
	log *zap.Logger
	dbs map[string]postgres.Database
}

type dbNetwork struct {
	ID string `db:"id"`
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

	q := `SELECT id FROM networks;`

	dbn := dbNetwork{}
	if err := m.dbs[postgres.DbKeto].QueryRowxContext(ctx, q).StructScan(&dbn); err != nil {
		return err
	}

	err := m.createUserRelations(ctx, currentTime, dbn.ID)
	if err != nil {
		return err
	}

	err = m.createChannelsOrThingsRelations(ctx, currentTime, dbn.ID, "things")
	if err != nil {
		return err
	}

	err = m.createChannelsOrThingsRelations(ctx, currentTime, dbn.ID, "channels")
	if err != nil {
		return err
	}

	return nil
}

func (m *M1KetoPolicies) Down() error {
	m.log.Info("do nothing for now")
	return nil
}

func (m *M1KetoPolicies) createUserRelations(ctx context.Context, currentTime time.Time, nid string) error {
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
		ID, _ := uuid.NewV4()
		q = `INSERT INTO keto_relation_tuples (shard_id, nid, namespace_id, object, relation, subject_id, commit_time)
			  VALUES (:shard_id, :nid, :namespace_id, :object, :relation, :subject_id, :commit_time);`
		params := map[string]interface{}{
			"shard_id":     ID,
			"nid":          nid,
			"namespace_id": 0,
			"object":       "users",
			"relation":     "member",
			"subject_id":   dbu.ID,
			"commit_time":  currentTime,
		}

		_, err = m.dbs[postgres.DbKeto].NamedExecContext(ctx, q, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *M1KetoPolicies) createChannelsOrThingsRelations(ctx context.Context, currentTime time.Time, nid string, table string) error {
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
			ID, _ := uuid.NewV4()
			q = `INSERT INTO keto_relation_tuples (shard_id, nid, namespace_id, object, relation,
				subject_set_namespace_id, subject_set_object, subject_set_relation, commit_time)
				VALUES (:shard_id, :nid, :namespace_id, :object, :relation,
				:subject_set_namespace_id, :subject_set_object, :subject_set_relation, :commit_time);`
			params := map[string]interface{}{
				"shard_id":                 ID,
				"nid":                      nid,
				"namespace_id":             0,
				"object":                   dbt.ID,
				"relation":                 relation,
				"subject_set_namespace_id": 0,
				"subject_set_object":       "authorities",
				"subject_set_relation":     "member",
				"commit_time":              currentTime,
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
			ID, _ := uuid.NewV4()
			q = `INSERT INTO keto_relation_tuples (shard_id, nid, namespace_id, object, relation, subject_id, commit_time)
				VALUES (:shard_id, :nid, :namespace_id, :object, :relation, :subject_id, :commit_time);`
			params := map[string]interface{}{
				"shard_id":     ID,
				"nid":          nid,
				"namespace_id": 0,
				"object":       dbt.ID,
				"relation":     relation,
				"subject_id":   dbu.ID,
				"commit_time":  currentTime,
			}

			_, err = m.dbs[postgres.DbKeto].NamedExecContext(ctx, q, params)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
