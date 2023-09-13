package producer

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Producer interface {
	// PublishSinkStatus to be used to publish the sink activity to the sinker
	PublishSinkStatus(ownerId string, sinkId string, status string, errorMessage string) error
}

type maestroProducer struct {
	logger      *zap.Logger
	streamRedis *redis.Client
}
