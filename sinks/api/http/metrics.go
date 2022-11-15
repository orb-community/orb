/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/sinks"
	"github.com/ns1labs/orb/sinks/backend"
	"time"
)

var _ sinks.SinkService = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	auth    mainflux.AuthServiceClient
	counter metrics.Counter
	latency metrics.Histogram
	svc     sinks.SinkService
}

// ListSinksInternal Will not count metrics since it is internal-service only rpc
func (m metricsMiddleware) ListSinksInternal(ctx context.Context, filter sinks.Filter) (sinks []sinks.Sink, err error) {
	return m.svc.ListSinksInternal(ctx, filter)
}

func (m metricsMiddleware) ChangeSinkStateInternal(ctx context.Context, sinkID string, msg string, ownerID string, state sinks.State) error {
	defer func(begin time.Time) {
		labels := []string{
			"method", "changeSinkStateInternal",
			"owner_id", ownerID,
			"sink_id", sinkID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ChangeSinkStateInternal(ctx, sinkID, msg, ownerID, state)
}

func (m metricsMiddleware) CreateSink(ctx context.Context, token string, s sinks.Sink) (sink sinks.Sink, _ error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return sinks.Sink{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "createSink",
			"owner_id", ownerID,
			"sink_id", sink.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.CreateSink(ctx, token, s)
}

func (m metricsMiddleware) UpdateSink(ctx context.Context, token string, s sinks.Sink) (sink sinks.Sink, err error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "updateSink",
			"owner_id", sink.MFOwnerID,
			"sink_id", sink.ID,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.UpdateSink(ctx, token, s)
}

func (m metricsMiddleware) ListSinks(ctx context.Context, token string, pm sinks.PageMetadata) (sink sinks.Page, err error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return sinks.Page{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "listSinks",
			"owner_id", ownerID,
			"sink_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ListSinks(ctx, token, pm)
}

func (m metricsMiddleware) ListBackends(ctx context.Context, token string) (_ []string, err error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return nil, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "listBackends",
			"owner_id", ownerID,
			"sink_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ListBackends(ctx, token)
}

func (m metricsMiddleware) ViewBackend(ctx context.Context, token string, key string) (_ backend.Backend, err error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return nil, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "viewBackend",
			"owner_id", ownerID,
			"sink_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewBackend(ctx, token, key)
}

func (m metricsMiddleware) ViewSink(ctx context.Context, token string, key string) (sinks.Sink, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return sinks.Sink{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "viewSink",
			"owner_id", ownerID,
			"sink_id", key,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewSink(ctx, token, key)
}

func (m metricsMiddleware) ViewSinkInternal(ctx context.Context, ownerID string, key string) (sinks.Sink, error) {
	defer func(begin time.Time) {
		labels := []string{
			"method", "viewSinkInternal",
			"owner_id", ownerID,
			"sink_id", key,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ViewSinkInternal(ctx, ownerID, key)
}

func (m metricsMiddleware) DeleteSink(ctx context.Context, token string, id string) (err error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "deleteSink",
			"owner_id", ownerID,
			"sink_id", id,
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.DeleteSink(ctx, token, id)
}

func (m metricsMiddleware) ValidateSink(ctx context.Context, token string, s sinks.Sink) (sinks.Sink, error) {
	ownerID, err := m.identify(token)
	if err != nil {
		return sinks.Sink{}, err
	}

	defer func(begin time.Time) {
		labels := []string{
			"method", "validateSink",
			"owner_id", ownerID,
			"sink_id", "",
		}

		m.counter.With(labels...).Add(1)
		m.latency.With(labels...).Observe(float64(time.Since(begin).Microseconds()))

	}(time.Now())

	return m.svc.ValidateSink(ctx, token, s)
}

func (m metricsMiddleware) identify(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return res.GetId(), nil
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(auth mainflux.AuthServiceClient, svc sinks.SinkService, counter metrics.Counter, latency metrics.Histogram) sinks.SinkService {
	return &metricsMiddleware{
		auth:    auth,
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
