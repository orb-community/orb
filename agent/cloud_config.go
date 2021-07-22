/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"strconv"
)

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "cloud_config_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS cloud_config (
						address TEXT NOT NULL,
						id TEXT	NOT NULL,
						key TEXT	NOT NULL,
						channel TEXT	NOT NULL,
						ts_created INTEGER NOT NULL
						)`,
				},
				Down: []string{
					"DROP TABLE cloud_config",
				},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "sqlite3", migrations, migrate.Up)

	return err
}

func autoProvision(c config.Config, db *sqlx.DB) (config.MQTTConfig, error) {

	err := migrateDB(db)
	if err != nil {
		return config.MQTTConfig{}, err
	}

	return config.MQTTConfig{}, errors.New("unable to auto provision agent")
}

func GetCloudConfig(c config.Config, db *sqlx.DB) (config.MQTTConfig, error) {

	// if MQTT is specified in the config file, always use that
	address, haveAddress := c.OrbAgent.Cloud["mqtt"]["address"]
	id, haveId := c.OrbAgent.Cloud["mqtt"]["id"]
	key, haveKey := c.OrbAgent.Cloud["mqtt"]["key"]
	channel, haveChannel := c.OrbAgent.Cloud["mqtt"]["channel"]

	if haveAddress && haveId && haveKey && haveChannel {
		return config.MQTTConfig{
			Address:   address,
			Id:        id,
			Key:       key,
			ChannelID: channel,
		}, nil
	}

	// if not, possibly attempt auto provision
	var ap bool
	var err error
	apConfig, haveAutoProvision := c.OrbAgent.Cloud["config"]["auto_provision"]
	if haveAutoProvision {
		ap, err = strconv.ParseBool(apConfig)
		if err != nil {
			return config.MQTTConfig{}, err
		}
	} else {
		ap = true
	}

	if ap {
		return autoProvision(c, db)
	} else {
		return config.MQTTConfig{}, errors.New("valid cloud MQTT config was not specified, and auto-provision was disabled")
	}

}
