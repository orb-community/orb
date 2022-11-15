package redis_test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"
)

var redisClient *redis.Client
var logger *zap.Logger

func TestMain(m *testing.M) {
	logger, _ = zap.NewProduction()
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	container, err := pool.Run("redis", "5.0-alpine", nil)
	if err != nil {
		logger.Fatal("could not start container: %s", zap.Error(err))
	}

	if err := pool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", container.GetPort("6379/tcp")),
			Password: "",
			DB: 0,
		})
		return redisClient.Ping(context.Background()).Err()
	}); err != nil {
		logger.Fatal("could not conncet to docker: %s", zap.Error(err))
	}

	code := m.Run()

	if err := pool.Purge(container); err != nil {
		logger.Fatal("could not purge container: %s", zap.Error(err))
	}

	os.Exit(code)
}