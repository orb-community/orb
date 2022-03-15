package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinker/config"
	"go.uber.org/zap"
	"strings"
)

const (
	keyPrefix = "sinker_key"
	idPrefix  = "sinker"
)

var _ config.ConfigRepo = (*sinkerCache)(nil)

type sinkerCache struct {
	client *redis.Client
	logger *zap.Logger
}

func NewSinkerCache(client *redis.Client, logger *zap.Logger) config.ConfigRepo {
	return &sinkerCache{client: client, logger: logger}
}

func (s *sinkerCache) Exists(sinkID string) bool {
	sinkConfig, err := s.Get(sinkID)
	if err != nil {
		return false
	}
	if sinkConfig.SinkID != "" {
		return true
	}
	return false
}

func (s *sinkerCache) Add(config config.SinkConfig) error {
	skey := fmt.Sprintf("%s-%s", keyPrefix, config.SinkID)
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err = s.client.Set(context.Background(), skey, bytes, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (s *sinkerCache) Remove(sinkID string) error {
	skey := fmt.Sprintf("%s-%s", keyPrefix, sinkID)
	if err := s.client.Del(context.Background(), skey).Err(); err != nil {
		return err
	}
	return nil
}

func (s *sinkerCache) Get(sinkID string) (config.SinkConfig, error) {
	skey := fmt.Sprintf("%s-%s", keyPrefix, sinkID)
	cachedConfig, err := s.client.Get(context.Background(), skey).Result()
	if err != nil {
		return config.SinkConfig{}, err
	}
	var cfgSinker config.SinkConfig
	if err := json.Unmarshal([]byte(cachedConfig), &cfgSinker); err != nil {
		return config.SinkConfig{}, err
	}
	return cfgSinker, nil
}

func (s *sinkerCache) Edit(config config.SinkConfig) error {
	if err := s.Remove(config.SinkID); err != nil {
		return err
	}
	if err := s.Add(config); err != nil {
		return err
	}
	return nil
}

func (s *sinkerCache) GetAll() ([]config.SinkConfig, error) {
	iter := s.client.Scan(context.Background(), 0, fmt.Sprintf("%s-*",keyPrefix), 0).Iterator()
	var configs []config.SinkConfig
	for iter.Next(context.Background()) {
		cfg, err := s.Get(strings.TrimPrefix(iter.Val(), fmt.Sprintf("%s-", keyPrefix)))
		if err != nil {
			s.logger.Error("failed to retrieve config", zap.Error(err))
			continue
		}
		configs = append(configs, cfg)
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	return configs, nil
}
