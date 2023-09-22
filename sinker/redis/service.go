package redis

import (
	"context"
	"github.com/orb-community/orb/sinker/redis/consumer"
	"github.com/orb-community/orb/sinker/redis/producer"
	"go.uber.org/zap"
)

type StreamsHandler interface {
	Start(ctx context.Context) error
}

type pubSubCacheHandler struct {
	logger             *zap.Logger
	sinkActivity       producer.SinkActivityProducer
	expirationListener consumer.SinkerKeyExpirationListener
}

var _ StreamsHandler = (*pubSubCacheHandler)(nil)

func NewPubSubCacheHandler(l *zap.Logger, sinkActivity producer.SinkActivityProducer, expirationListener consumer.SinkerKeyExpirationListener) StreamsHandler {
	return &pubSubCacheHandler{logger: l, sinkActivity: sinkActivity, expirationListener: expirationListener}
}

func (p *pubSubCacheHandler) Start(ctx context.Context) error {
	err := p.expirationListener.SubscribeToKeyExpiration(ctx)
	if err != nil {
		p.logger.Error("error subscribing to key expiration", zap.Error(err))
		return err
	}
	return nil
}
