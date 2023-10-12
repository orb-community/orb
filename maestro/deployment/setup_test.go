package deployment

import (
	"github.com/jmoiron/sqlx"
	"github.com/orb-community/orb/maestro/postgres"
	"github.com/orb-community/orb/pkg/config"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.uber.org/zap"
	"os"
	"testing"
)

var logger *zap.Logger
var pg *sqlx.DB

func TestMain(m *testing.M) {
	logger, _ = zap.NewProduction()
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Fatal("could not connect to docker:", zap.Error(err))
	}

	// Pull the PostgreSQL Docker image
	postgresImage := "postgres:latest"
	err = pool.Client.PullImage(docker.PullImageOptions{
		Repository: postgresImage,
		Tag:        "latest",
	}, docker.AuthConfiguration{})
	if err != nil {
		logger.Fatal("Could not pull Docker image:", zap.Error(err))
	}

	// Create a PostgreSQL container
	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=testdb",
	})
	if err != nil {
		logger.Fatal("Could not start PostgreSQL container", zap.Error(err))
	}

	retryF := func() error {
		localTest := config.PostgresConfig{
			Host:    "localhost",
			Port:    resource.GetPort("5432/tcp"),
			User:    "postgres",
			Pass:    "secret",
			DB:      "testdb",
			SSLMode: "disable",
		}
		pg, err = postgres.Connect(localTest)
		if err != nil {
			return err
		}

		return pg.Ping()
	}
	if err := pool.Retry(retryF); err != nil {
		logger.Fatal("could not connect to docker: %s", zap.Error(err))
	}
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		logger.Fatal("could not purge container: %s", zap.Error(err))
	}

	os.Exit(code)
}
