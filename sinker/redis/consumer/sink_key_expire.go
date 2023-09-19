package consumer

import (
	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/sinker/redis/producer"
	"go.uber.org/zap"
)

type SinkerKeyExpirationListener interface {
	// Listen to the sinker key expiration
	SubscribeToKeyExpiration() error
	ReceiveMessage(message interface{}) error
}

type sinkerKeyExpirationListener struct {
	logger           *zap.Logger
	cacheRedisClient redis.Client
	idleProducer     producer.SinkIdleProducer
}
