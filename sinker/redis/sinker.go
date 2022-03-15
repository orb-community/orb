package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinker/config"
)

const (
	keyPrefix = "sinker_key"
	idPrefix  = "sinker"
)

var _ config.ConfigRepo = (*sinkerCache)(nil)

type sinkerCache struct {
	client *redis.Client
}

func NewSinkerCache(client *redis.Client) config.ConfigRepo {
	return &sinkerCache{client: client}
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
	if err := s.Add(config); err != nil {
		return err
	}
	return nil
}

func (s *sinkerCache) GetAll() ([]config.SinkConfig, error) {
	iter := s.client.Scan(context.Background(), 0, fmt.Sprintf("%s-*",keyPrefix), 0).Iterator()
	for iter.Next(context.Background()) {
		fmt.Println("keys", iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	list := s.client.ZRange(context.Background(), "sinker", 0, -1)
	for k, v := range list.Val() {
		fmt.Println("keys", k)
		fmt.Println("value", v)
	}

	return nil, nil
}
