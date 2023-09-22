package consumer

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/sinker/redis/producer"
	"go.uber.org/zap"
)

type SinkerKeyExpirationListener interface {
	// SubscribeToKeyExpiration Listen to the sinker key expiration
	SubscribeToKeyExpiration(ctx context.Context) error
	// ReceiveMessage to be used to receive the message from the sinker key expiration, async
	ReceiveMessage(ctx context.Context, message interface{}) error
}

type sinkerKeyExpirationListener struct {
	logger           *zap.Logger
	cacheRedisClient *redis.Client
	idleProducer     producer.SinkIdleProducer
}

func NewSinkerKeyExpirationListener(l *zap.Logger, cacheRedisClient *redis.Client, idleProducer producer.SinkIdleProducer) SinkerKeyExpirationListener {
	logger := l.Named("sinker_key_expiration_listener")
	return &sinkerKeyExpirationListener{logger: logger, cacheRedisClient: cacheRedisClient, idleProducer: idleProducer}
}

// SubscribeToKeyExpiration to be used to subscribe to the sinker key expiration
func (s *sinkerKeyExpirationListener) SubscribeToKeyExpiration(ctx context.Context) error {
	go func() {
		pubsub := s.cacheRedisClient.Subscribe(ctx, "__keyevent@0__:expired")
		defer func(pubsub *redis.PubSub) {
			_ = pubsub.Close()
		}(pubsub)
		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				s.logger.Info(fmt.Sprintf("key %s expired", msg.Payload))
				subCtx := context.WithValue(ctx, "msg", msg.Payload)
				err := s.ReceiveMessage(subCtx, msg.Payload)
				if err != nil {
					s.logger.Error("error receiving message", zap.Error(err))
					return
				}
			}
		}
	}()
	return nil
}

// ReceiveMessage to be used to receive the message from the sinker key expiration
func (s *sinkerKeyExpirationListener) ReceiveMessage(ctx context.Context, message interface{}) error {
	// goroutine
	//sinkID := msg.Payload
	//event := producer.SinkIdleEvent{
	//	OwnerID:  "owner_id",
	//	SinkID: "sink_id",
	//	State: "idle",
	//}
	//s.idleProducer.PublishSinkIdle(ctx, event)
	return nil
}
