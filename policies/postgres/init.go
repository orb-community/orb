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
				Id: "policies_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS agent_policies (
						id			   UUID NOT NULL DEFAULT gen_random_uuid(),

						name           TEXT NOT NULL,
						mf_owner_id    UUID NOT NULL,
						description    TEXT NOT NULL DEFAULT '',

						orb_tags       JSONB NOT NULL DEFAULT '{}',

						backend        TEXT NOT NULL,
						version        INTEGER NOT NULL DEFAULT 0,
						policy		   JSONB NOT NULL DEFAULT '{}',

                        ts_created     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
						PRIMARY KEY (name, mf_owner_id),
					    UNIQUE(id)
					)`,
					`CREATE TABLE IF NOT EXISTS datasets (
						id			   UUID NOT NULL DEFAULT gen_random_uuid(),

						name             TEXT NOT NULL,
						mf_owner_id      UUID NOT NULL,

						valid			 BOOLEAN,
						agent_group_id	 UUID,
						agent_policy_id  UUID,
						sink_ids	     UUID[],
						
						metadata       JSONB NOT NULL DEFAULT '{}',
						tags		   JSONB NOT NULL DEFAULT '{}',

                        ts_created     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
						PRIMARY KEY (name, mf_owner_id),
					    UNIQUE(id)
					)`,
				},
				Down: []string{
					"DROP TABLE agent_policies",
					"DROP TABLE datasets",
				},
			},
			{
				Id: "policies_2",
				Up: []string{
					`ALTER TABLE IF EXISTS agent_policies ADD COLUMN IF NOT EXISTS
					 schema_version TEXT NOT NULL`,
				},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)

	return err
}
