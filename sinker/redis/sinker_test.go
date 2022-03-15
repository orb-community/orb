package redis_test

import (
	"fmt"
	r "github.com/go-redis/redis/v8"
	"github.com/mainflux/mainflux/pkg/uuid"
	"github.com/ns1labs/orb/pkg/errors"
	config2 "github.com/ns1labs/orb/sinker/config"
	"github.com/ns1labs/orb/sinker/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

var idProvider = uuid.New()

func TestSinkerConfigSave(t *testing.T) {
	sinkerCache := redis.NewSinkerCache(redisClient, logger)
	config := config2.SinkConfig{
		SinkID:          "123",
		OwnerID:         "test",
		Url:             "localhost",
		User:            "user",
		Password:        "password",
		State:           0,
		Msg:             "",
		LastRemoteWrite: time.Time{},
	}

	err := sinkerCache.Add(config)
	require.Nil(t, err, fmt.Sprintf("save sinker config to cache: expected nil got %s", err))

	cases := map[string]struct {
		config config2.SinkConfig
		err    error
	}{
		"Save sinker to cache": {
			config: config2.SinkConfig{
				SinkID:          "124",
				OwnerID:         "test",
				Url:             "localhost",
				User:            "user",
				Password:        "password",
				State:           0,
				Msg:             "",
				LastRemoteWrite: time.Time{},
			},
			err: nil,
		},
		"Save already cached sinker config to cache": {
			config: config,
			err:    nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := sinkerCache.Add(tc.config)
			assert.Nil(t, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestGetSinkerConfig(t *testing.T) {
	sinkerCache := redis.NewSinkerCache(redisClient, logger)
	config := config2.SinkConfig{
		SinkID:          "123",
		OwnerID:         "test",
		Url:             "localhost",
		User:            "user",
		Password:        "password",
		State:           0,
		Msg:             "",
		LastRemoteWrite: time.Time{},
	}

	err := sinkerCache.Add(config)
	require.Nil(t, err, fmt.Sprintf("save sinker config to cache: expected nil got %s", err))

	cases := map[string]struct {
		sinkID string
		config config2.SinkConfig
		err    error
	}{
		"Get Config by existing sinker-key": {
			sinkID: "123",
			config: config,
			err:    nil,
		},
		"Get Config by non-existing sinker-key": {
			sinkID: "000",
			config: config2.SinkConfig{},
			err:    r.Nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			sinkConfig, err := sinkerCache.Get(tc.sinkID)
			assert.True(t, reflect.DeepEqual(tc.config, sinkConfig), fmt.Sprintf("%s: expected %v got %v", desc, tc.config, sinkConfig))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestGetAllSinkerConfig(t *testing.T) {
	sinkerCache := redis.NewSinkerCache(redisClient, logger)

	sinksConfig := map[string]struct {
		config config2.SinkConfig
	}{
		"config 1": {
			config: config2.SinkConfig{
				SinkID:          "123",
				OwnerID:         "test",
				Url:             "localhost",
				User:            "user",
				Password:        "password",
				State:           0,
				Msg:             "",
				LastRemoteWrite: time.Time{},
			},
		},
		"config 2": {
			config: config2.SinkConfig{
				SinkID:          "134",
				OwnerID:         "test",
				Url:             "localhost",
				User:            "user",
				Password:        "password",
				State:           0,
				Msg:             "",
				LastRemoteWrite: time.Time{},
			},
		},
	}

	for _, val := range sinksConfig {
		err := sinkerCache.Add(val.config)
		require.Nil(t, err, fmt.Sprintf("save sinker config to cache: expected nil got %s", err))
	}

	cases := map[string]struct {
		size int
		err  error
	}{
		"Get Config by existing sinker-key": {
			size: 2,
			err:  nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			sinksConfig, err := sinkerCache.GetAll()
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			assert.GreaterOrEqual(t, len(sinksConfig), tc.size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, len(sinksConfig)))
		})
	}
}
