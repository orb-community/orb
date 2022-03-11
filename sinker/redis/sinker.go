package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinker/config"
)

const (
	keyPrefix = "sinker_key"
	idPrefix = "sinker"
)

var _ config.ConfigRepo = (*sinkerCache)(nil)

type sinkerCache struct {
	client *redis.Client
}

func NewSinkerCache(client *redis.Client) config.ConfigRepo {
	return &sinkerCache{client: client}
}

func (s sinkerCache) Exists(sinkID string) bool {
	panic("implement me")
}

func (s sinkerCache) Add(config config.SinkConfig) error {
	panic("implement me")
}

func (s sinkerCache) Get(sinkID string) (config.SinkConfig, error) {
	panic("implement me")
}

func (s sinkerCache) Edit(config config.SinkConfig) error {
	panic("implement me")
}

func (s sinkerCache) GetAll() ([]config.SinkConfig, error) {
	panic("implement me")
}
