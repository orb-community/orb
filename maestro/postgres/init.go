package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required for SQL access
	"github.com/orb-community/orb/pkg/config"
	migrate "github.com/rubenv/sql-migrate"
)

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg config.PostgresConfig) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.DB, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS deployments (
					    id			    UUID NOT NULL DEFAULT gen_random_uuid(),
					    owner_id                VARCHAR(255),
					    sink_id                 VARCHAR(255),
					    backend                 VARCHAR(255), 
					    config                  JSONB,
					    last_status             VARCHAR(255),
					    last_status_update      TIMESTAMP,
					    last_error_message      VARCHAR(255),
					    last_error_time         TIMESTAMP,
					    collector_name          VARCHAR(255),
					    last_collector_deploy_time TIMESTAMP,
					    last_collector_stop_time   TIMESTAMP
					);`,
				},
				Down: []string{
					"DROP TABLE deployments",
				},
			},
		},
	}
	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)

	return err
}
