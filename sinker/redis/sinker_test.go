package redis_test

import (
	"context"
	"fmt"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinker/redis/consumer"
	"github.com/orb-community/orb/sinker/redis/producer"
	"testing"
	"time"

	"github.com/mainflux/mainflux/pkg/uuid"
	"github.com/orb-community/orb/pkg/errors"
	config2 "github.com/orb-community/orb/sinker/config"
	"github.com/orb-community/orb/sinker/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var idProvider = uuid.New()

func TestSinkActivityStoreAndMessage(t *testing.T) {
	// Create SinkActivityService
	sinkTTLSvc := producer.NewSinkerKeyService(logger, redisClient)
	sinkActivitySvc := producer.NewSinkActivityProducer(logger, redisClient, sinkTTLSvc)
	args := []struct {
		testCase string
		event    producer.SinkActivityEvent
	}{
		{
			testCase: "sink activity for new sink",
			event: producer.SinkActivityEvent{
				OwnerID:   "1",
				SinkID:    "1",
				State:     "active",
				Size:      "40",
				Timestamp: time.Now(),
			},
		},
		{
			testCase: "sink activity for existing sink",
			event: producer.SinkActivityEvent{
				OwnerID:   "1",
				SinkID:    "1",
				State:     "active",
				Size:      "55",
				Timestamp: time.Now(),
			},
		},
		{
			testCase: "sink activity for another new sink",
			event: producer.SinkActivityEvent{
				OwnerID:   "2",
				SinkID:    "1",
				State:     "active",
				Size:      "37",
				Timestamp: time.Now(),
			},
		},
	}
	for _, tt := range args {
		ctx := context.WithValue(context.Background(), "test_case", tt.testCase)
		err := sinkActivitySvc.PublishSinkActivity(ctx, tt.event)
		require.NoError(t, err, fmt.Sprintf("%s: unexpected error: %s", tt.testCase, err))
	}
	logger.Debug("debugging breakpoint")
}

func TestSinkIdle(t *testing.T) {
	sinkTTLSvc := producer.NewSinkerKeyService(logger, redisClient)
	sinkActivitySvc := producer.NewSinkActivityProducer(logger, redisClient, sinkTTLSvc)
	sinkIdleSvc := producer.NewSinkIdleProducer(logger, redisClient)
	sinkExpire := consumer.NewSinkerKeyExpirationListener(logger, redisClient, sinkIdleSvc)
	event := producer.SinkActivityEvent{
		OwnerID:   "1",
		SinkID:    "1",
		State:     "active",
		Size:      "40",
		Timestamp: time.Now(),
	}
	ctx := context.WithValue(context.Background(), "test", "TestSinkIdle")
	err := sinkExpire.SubscribeToKeyExpiration(ctx)
	require.NoError(t, err, fmt.Sprintf("unexpected error: %s", err))
	err = sinkActivitySvc.PublishSinkActivity(ctx, event)
	require.NoError(t, err, fmt.Sprintf("unexpected error: %s", err))
	err = sinkTTLSvc.RenewSinkerKeyInternal(ctx, producer.SinkerKey{
		OwnerID: "1",
		SinkID:  "1",
	}, 30*time.Second)
	require.NoError(t, err, fmt.Sprintf("unexpected error: %s", err))
}

func TestSinkerConfigSave(t *testing.T) {
	sinkerCache := redis.NewSinkerCache(redisClient, logger)
	var config config2.SinkConfig
	config.SinkID = "123"
	config.OwnerID = "test"
	config.Config = types.Metadata{
		"authentication": types.Metadata{
			"password": "password",
			"type":     "basicauth",
			"username": "user",
		},
		"exporter": types.Metadata{
			"headers": map[string]string{
				"X-Tenant": "MY_TENANT_1",
			},
			"remote_host": "localhost",
		},
		"opentelemetry": "enabled",
	}

	config.State = 0
	config.Msg = ""
	config.LastRemoteWrite = time.Time{}

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
				Config:          config.Config,
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
	var config config2.SinkConfig
	config.SinkID = "123"
	config.OwnerID = "test"
	config.Config = types.Metadata{
		"authentication": types.Metadata{
			"password": "password",
			"type":     "basicauth",
			"username": "user",
		},
		"exporter": types.Metadata{
			"headers": map[string]string{
				"X-Tenant": "MY_TENANT_1",
			},
			"remote_host": "localhost",
		},
		"opentelemetry": "enabled",
	}
	config.State = 0
	config.Msg = ""
	config.LastRemoteWrite = time.Time{}

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
			err:    errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			sinkConfig, err := sinkerCache.Get(tc.config.OwnerID, tc.sinkID)
			assert.Equal(t, tc.config.SinkID, sinkConfig.SinkID, fmt.Sprintf("%s: expected %s got %s", desc, tc.config.SinkID, sinkConfig.SinkID))
			assert.Equal(t, tc.config.State, sinkConfig.State, fmt.Sprintf("%s: expected %s got %s", desc, tc.config.State, sinkConfig.State))
			assert.Equal(t, tc.config.OwnerID, sinkConfig.OwnerID, fmt.Sprintf("%s: expected %s got %s", desc, tc.config.OwnerID, sinkConfig.OwnerID))
			assert.Equal(t, tc.config.Msg, sinkConfig.Msg, fmt.Sprintf("%s: expected %s got %s", desc, tc.config.Msg, sinkConfig.Msg))
			assert.Equal(t, tc.config.LastRemoteWrite, sinkConfig.LastRemoteWrite, fmt.Sprintf("%s: expected %s got %s", desc, tc.config.LastRemoteWrite, sinkConfig.LastRemoteWrite))
			if tc.config.Config != nil {
				_, ok := sinkConfig.Config["authentication"]
				assert.True(t, ok, fmt.Sprintf("%s: should contain authentication metadata", desc))
				_, ok = sinkConfig.Config["exporter"]
				assert.True(t, ok, fmt.Sprintf("%s: should contain exporter metadata", desc))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestGetAllSinkerConfig(t *testing.T) {
	sinkerCache := redis.NewSinkerCache(redisClient, logger)
	var config config2.SinkConfig
	config.SinkID = "123"
	config.OwnerID = "test"
	config.State = 0
	config.Msg = ""
	config.Config = types.Metadata{
		"authentication": types.Metadata{
			"password": "password",
			"type":     "basicauth",
			"username": "user",
		},
		"exporter": types.Metadata{
			"headers": map[string]string{
				"X-Tenant": "MY_TENANT_1",
			},
			"remote_host": "localhost",
		},
		"opentelemetry": "enabled",
	}
	config.LastRemoteWrite = time.Time{}
	sinksConfig := map[string]struct {
		config config2.SinkConfig
	}{
		"config 1": {
			config: config2.SinkConfig{
				SinkID:          "123",
				OwnerID:         "test",
				Config:          config.Config,
				State:           0,
				Msg:             "",
				LastRemoteWrite: time.Time{},
			},
		},
		"config 2": {
			config: config2.SinkConfig{
				SinkID:          "134",
				OwnerID:         "test",
				Config:          config.Config,
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
		size    int
		ownerID string
		err     error
	}{
		"Get Config by existing sinker-key": {
			size:    2,
			ownerID: "test",
			err:     nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			sinksConfig, err := sinkerCache.GetAll(tc.ownerID)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			assert.GreaterOrEqual(t, len(sinksConfig), tc.size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, len(sinksConfig)))
		})
	}
}
