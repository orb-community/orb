package producer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type SinkActivityProducer interface {
	// PublishSinkActivity to be used to publish the sink activity to the sinker, mainly used by Otel Bridge Service
	PublishSinkActivity(ctx context.Context, event SinkActivityEvent) error
}

type SinkActivityEvent struct {
	OwnerID   string
	SinkID    string
	State     string
	Size      string
	Timestamp time.Time
}

func (s *SinkActivityEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"owner_id":  s.OwnerID,
		"sink_id":   s.SinkID,
		"state":     s.State,
		"size":      s.Size,
		"timestamp": s.Timestamp.Format(time.RFC3339),
	}
}

var _ SinkActivityProducer = (*sinkActivityProducer)(nil)

type sinkActivityProducer struct {
	logger            *zap.Logger
	redisStreamClient *redis.Client
	sinkTTL           SinkerKeyService
}

func NewSinkActivityProducer(l *zap.Logger, redisStreamClient *redis.Client, sinkTTL SinkerKeyService) SinkActivityProducer {
	logger := l.Named("sink_activity_producer")
	return &sinkActivityProducer{logger: logger, redisStreamClient: redisStreamClient, sinkTTL: sinkTTL}
}

// PublishSinkActivity BridgeService will notify stream of sink activity
func (sp *sinkActivityProducer) PublishSinkActivity(ctx context.Context, event SinkActivityEvent) error {
	const maxLen = 1000
	record := &redis.XAddArgs{
		Stream: "orb.sink_activity",
		Values: event.Encode(),
		MaxLen: maxLen,
		Approx: true,
	}
	err := sp.redisStreamClient.XAdd(ctx, record).Err()
	if err != nil {
		sp.logger.Error("error sending event to sinker event store", zap.Error(err))
	}
	err = sp.sinkTTL.AddNewSinkerKey(ctx, SinkerKey{
		OwnerID:      event.OwnerID,
		SinkID:       event.SinkID,
		Size:         event.Size,
		LastActivity: event.Timestamp,
	})
	return err
}
