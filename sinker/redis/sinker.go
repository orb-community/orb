package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinker"
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

func (s *sinkerCache) Exists(ownerID string, sinkID string) bool {
	sinkConfig, err := s.Get(ownerID, sinkID)
	if err != nil {
		return false
	}
	if sinkConfig.SinkID != "" {
		return true
	}
	return false
}

func (s *sinkerCache) Add(config config.SinkConfig) error {
	skey := fmt.Sprintf("%s-%s:%s", keyPrefix, config.OwnerID, config.SinkID)
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err = s.client.Set(context.Background(), skey, bytes, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (s *sinkerCache) Remove(ownerID string, sinkID string) error {
	skey := fmt.Sprintf("%s-%s:%s", keyPrefix, ownerID, sinkID)
	if err := s.client.Del(context.Background(), skey).Err(); err != nil {
		return err
	}
	return nil
}

func (s *sinkerCache) Get(ownerID string, sinkID string) (config.SinkConfig, error) {
	if ownerID == "" || sinkID == "" {
		return config.SinkConfig{}, sinker.ErrNotFound
	}
	skey := fmt.Sprintf("%s-%s:%s", keyPrefix, ownerID, sinkID)
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
	if err := s.Remove(config.OwnerID, config.SinkID); err != nil {
		return err
	}
	if err := s.Add(config); err != nil {
		return err
	}
	return nil
}

func (s *sinkerCache) 	GetAllOwners() ([]string, error){
	iter := s.client.Scan(context.Background(), 0, fmt.Sprintf("%s-*",keyPrefix), 0).Iterator()
	var owners []string
	for iter.Next(context.Background()) {
		keys := strings.Split(strings.TrimPrefix(iter.Val(), fmt.Sprintf("%s-", keyPrefix)), ":")
		if len(keys) > 1 {
			owners = append(owners, keys[0])
		}
	}
	if err := iter.Err(); err != nil {
		s.logger.Error("failed to retrieve config", zap.Error(err))
		return owners, err
	}
	return owners, nil
}

func (s *sinkerCache) GetAll(ownerID string) ([]config.SinkConfig, error) {
	iter := s.client.Scan(context.Background(), 0, fmt.Sprintf("%s-%s:*",keyPrefix, ownerID), 0).Iterator()
	var configs []config.SinkConfig
	for iter.Next(context.Background()) {
		keys := strings.Split(strings.TrimPrefix(iter.Val(), fmt.Sprintf("%s-", keyPrefix)), ":")
		sinkID := ""
		if len(keys) > 1 {
			sinkID = keys[1]
		}
		cfg, err := s.Get(ownerID, sinkID)
		if err != nil {
			s.logger.Error("failed to retrieve config", zap.Error(err))
			continue
		}
		configs = append(configs, cfg)
	}
	if err := iter.Err(); err != nil {
		s.logger.Error("failed to retrieve config", zap.Error(err))
		return configs, err
	}

	return configs, nil
}
