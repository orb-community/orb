package producer

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	streamID  = "orb.maestro.sink_status"
	streamLen = 1000
)

type SinkStatusEvent struct {
	ownerId      string
	sinkId       string
	status       string
	errorMessage string
}

func (e SinkStatusEvent) Encode() map[string]interface{} {
	return map[string]interface{}{
		"owner_id":      e.ownerId,
		"sink_id":       e.sinkId,
		"status":        e.status,
		"error_message": e.errorMessage,
		"timestamp":     time.Now().Format(time.RFC3339),
	}
}

type Producer interface {
	// PublishSinkStatus to be used to publish the sink activity to the sinker
	PublishSinkStatus(ctx context.Context, ownerId string, sinkId string, status string, errorMessage string) error
}

type maestroProducer struct {
	logger      *zap.Logger
	streamRedis *redis.Client
}

func NewMaestroProducer(logger *zap.Logger, streamRedis *redis.Client) Producer {
	return &maestroProducer{logger: logger, streamRedis: streamRedis}
}

// PublishSinkStatus to be used to publish the sink activity to the sinker
func (p *maestroProducer) PublishSinkStatus(ctx context.Context, ownerId string, sinkId string, status string, errorMessage string) error {
	event := SinkStatusEvent{
		ownerId:      ownerId,
		sinkId:       sinkId,
		status:       status,
		errorMessage: errorMessage,
	}
	streamEvent := event.Encode()
	record := &redis.XAddArgs{
		Stream: streamID,
		MaxLen: streamLen,
		Approx: true,
		Values: streamEvent,
	}
	cmd := p.streamRedis.XAdd(ctx, record)
	if cmd.Err() != nil {
		p.logger.Error("error sending event to maestro event store", zap.Error(cmd.Err()))
		return cmd.Err()
	}
	return nil
}
