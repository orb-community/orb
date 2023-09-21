package producer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type SinkIdleEvent struct {
	OwnerID   string
	SinkID    string
	State     string
	Size      string
	Timestamp time.Time
}

func (s *SinkIdleEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"owner_id":  s.OwnerID,
		"sink_id":   s.SinkID,
		"state":     s.State,
		"size":      s.Size,
		"timestamp": s.Timestamp.Format(time.RFC3339),
	}
}

type SinkIdleProducer interface {
	// PublishSinkIdle to be used to publish the sink activity to the sinker, mainly used by Otel Bridge Service
	PublishSinkIdle(ctx context.Context, event SinkIdleEvent) error
}

var _ SinkIdleProducer = (*sinkIdleProducer)(nil)

type sinkIdleProducer struct {
	logger            *zap.Logger
	redisStreamClient *redis.Client
}

func NewSinkIdleProducer(logger *zap.Logger, redisStreamClient *redis.Client) SinkIdleProducer {
	return &sinkIdleProducer{logger: logger, redisStreamClient: redisStreamClient}
}

func (s *sinkIdleProducer) PublishSinkIdle(ctx context.Context, event SinkIdleEvent) error {
	const maxLen = 1000
	record := &redis.XAddArgs{
		Stream: "orb.sinker.sink_idle",
		Values: event.Encode(),
		MaxLen: maxLen,
		Approx: true,
	}
	err := s.redisStreamClient.XAdd(ctx, record).Err()
	if err != nil {
		s.logger.Error("error sending event to sinker event store", zap.Error(err))
	}
	return err
}
