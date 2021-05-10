// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres

import (
	"fmt"
	"github.com/ns1labs/orb/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required for SQL access
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
				Id: "fleet_1",
				Up: []string{
					`CREATE TYPE agent_state AS ENUM ('new', 'online', 'offline', 'stale', 'removed');`,
					`CREATE TABLE IF NOT EXISTS agents (
						mf_thing_id        UUID UNIQUE,
						owner              VARCHAR(254),
                        ts_created         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

						orb_tags           JSONB NOT NULL DEFAULT '{}',

						agent_tags         JSONB NOT NULL DEFAULT '{}',
						agent_metadata     JSONB NOT NULL DEFAULT '{}',

						state              agent_state NOT NULL DEFAULT 'new',

						last_hb_data       JSONB NOT NULL DEFAULT '{}',
                        ts_last_hb         TIMESTAMPTZ DEFAULT NULL,

						PRIMARY KEY (mf_thing_id, owner)
					)`,
					`CREATE TABLE IF NOT EXISTS selectors (
						name        	   TEXT UNIQUE,
						owner              VARCHAR(254),
	
						config             JSONB NOT NULL DEFAULT '{}',
                        ts_created         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
						PRIMARY KEY (name, owner)
					)`,
				},
				Down: []string{
					"DROP TABLE agents",
					"DROP TABLE selectors",
				},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)

	return err
}
