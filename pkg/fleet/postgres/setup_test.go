// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres_test

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/pkg/config"
	"github.com/ns1labs/orb/pkg/fleet/postgres"

	dockertest "github.com/ory/dockertest/v3"
)

var (
	testLog, _ = zap.NewDevelopment()
	db         *sqlx.DB
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	cfg := []string{
		"POSTGRES_USER=test",
		"POSTGRES_PASSWORD=test",
		"POSTGRES_DB=test",
	}
	container, err := pool.Run("postgres", "13.2-alpine", cfg)
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	port := container.GetPort("5432/tcp")

	if err := pool.Retry(func() error {
		url := fmt.Sprintf("host=localhost port=%s user=test dbname=test password=test sslmode=disable", port)
		db, err = sqlx.Open("postgres", url)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	dbConfig := config.PostgresConfig{
		Host:        "localhost",
		Port:        port,
		User:        "test",
		Pass:        "test",
		DB:          "test",
		SSLMode:     "disable",
		SSLCert:     "",
		SSLKey:      "",
		SSLRootCert: "",
	}

	if db, err = postgres.Connect(dbConfig); err != nil {
		log.Fatalf("Could not setup test DB connection: %s", err)
	}

	testLog.Debug("connected to database")

	code := m.Run()

	// Defers will not be run when using os.Exit
	db.Close()
	/*
		if err := pool.Purge(container); err != nil {
			log.Fatalf("Could not purge container: %s", err)
		}

	*/

	testLog.Debug("purged database")

	os.Exit(code)
}
