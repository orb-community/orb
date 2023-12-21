package service

import (
	"context"
	"github.com/go-kit/kit/metrics"
	maestroredis "github.com/orb-community/orb/maestro/redis"
	"go.uber.org/zap"
	"time"
)

type tracingService struct {
	logger      *zap.Logger
	counter     metrics.Counter
	latency     metrics.Histogram
	nextService EventService
}

func NewTracingService(logger *zap.Logger, service EventService, counter metrics.Counter, latency metrics.Histogram) EventService {
	return &tracingService{logger: logger, nextService: service, counter: counter, latency: latency}
}

func (t *tracingService) HandleSinkCreate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	defer func(begun time.Time) {
		labels := []string{
			"method", "HandleSinkCreate",
			"sink_id", event.SinkID,
			"owner_id", event.Owner,
		}
		t.counter.With(labels...).Add(1)
		t.latency.With(labels...).Observe(float64(time.Since(begun).Microseconds()))
	}(time.Now())
	return t.nextService.HandleSinkCreate(ctx, event)
}

func (t *tracingService) HandleSinkUpdate(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	defer func(begun time.Time) {
		labels := []string{
			"method", "HandleSinkCreate",
			"sink_id", event.SinkID,
			"owner_id", event.Owner,
		}
		t.counter.With(labels...).Add(1)
		t.latency.With(labels...).Observe(float64(time.Since(begun).Microseconds()))
	}(time.Now())
	return t.nextService.HandleSinkUpdate(ctx, event)
}

func (t *tracingService) HandleSinkDelete(ctx context.Context, event maestroredis.SinksUpdateEvent) error {
	defer func(begun time.Time) {
		labels := []string{
			"method", "HandleSinkCreate",
			"sink_id", event.SinkID,
			"owner_id", event.Owner,
		}
		t.counter.With(labels...).Add(1)
		t.latency.With(labels...).Observe(float64(time.Since(begun).Microseconds()))
	}(time.Now())
	return t.nextService.HandleSinkDelete(ctx, event)
}

func (t *tracingService) HandleSinkActivity(ctx context.Context, event maestroredis.SinkerUpdateEvent) error {
	defer func(begun time.Time) {
		labels := []string{
			"method", "HandleSinkCreate",
			"sink_id", event.SinkID,
			"owner_id", event.OwnerID,
		}
		t.counter.With(labels...).Add(1)
		t.latency.With(labels...).Observe(float64(time.Since(begun).Microseconds()))
	}(time.Now())
	return t.nextService.HandleSinkActivity(ctx, event)
}

func (t *tracingService) HandleSinkIdle(ctx context.Context, event maestroredis.SinkerUpdateEvent) error {
	defer func(begun time.Time) {
		labels := []string{
			"method", "HandleSinkCreate",
			"sink_id", event.SinkID,
			"owner_id", event.OwnerID,
		}
		t.counter.With(labels...).Add(1)
		t.latency.With(labels...).Observe(float64(time.Since(begun).Microseconds()))
	}(time.Now())
	return t.nextService.HandleSinkIdle(ctx, event)
}
