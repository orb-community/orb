package redis_test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"testing"
)

var redisClient *redis.Client

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	container, err := pool.Run("redis", "5.0-alpine", nil)
	if err != nil {
		log.Fatalf("could not start container: %s", err)
	}

	if err := pool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", container.GetPort("6379/tcp")),
			Password: "",
			DB: 0,
		})
		return redisClient.Ping(context.Background()).Err()
	}); err != nil {
		log.Fatalf("could not conncet to docker: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(container); err != nil {
		log.Fatalf("could not purge container: %s", err)
	}

	os.Exit(code)
}