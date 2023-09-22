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
			Addr:     fmt.Sprintf("localhost:%s", container.GetPort("6379/tcp")),
			Password: "",
			DB:       0,
		})
		return redisClient.Ping(context.Background()).Err()
	}); err != nil {
		logger.Fatal("could not connect to docker: %s", zap.Error(err))
	}

	code := m.Run()

	if err := pool.Purge(container); err != nil {
		logger.Fatal("could not purge container: %s", zap.Error(err))
	}

	os.Exit(code)
}

func OnceReceiver(ctx context.Context, streamID string) error {
	go func() {
		count := 0
		err := redisClient.XGroupCreateMkStream(ctx, streamID, "unit-test", "$").Err()
		if err != nil {
			logger.Warn("error during create group", zap.Error(err))
		}
		for {
			// Redis Subscribe to stream
			if redisClient != nil {
				// create the group, or ignore if it already exists
				streams, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
					Consumer: "test_consumer",
					Group:    "unit-test",
					Streams:  []string{streamID, ">"},
					Count:    10,
				}).Result()
				if err != nil || len(streams) == 0 {
					continue
				}
				for _, stream := range streams {
					for _, msg := range stream.Messages {
						logger.Info("received message", zap.Any("message", msg.Values))
						count++
					}
				}
				if count > 0 {
					return
				}
			}
		}
	}()
	return nil
}
