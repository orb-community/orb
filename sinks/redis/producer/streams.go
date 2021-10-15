// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/sinks"
	"github.com/ns1labs/orb/sinks/backend"
)

const (
	streamID  = "orb.sinks"
	streamLen = 1000
)

var _ sinks.SinkService = (*eventStore)(nil)

type eventStore struct {
	svc    sinks.SinkService
	client *redis.Client
}

func (es eventStore) ChangeSinkStateInternal(ctx context.Context, sinkID string, msg string, ownerID string, state sinks.State) error {
	return es.svc.ChangeSinkStateInternal(ctx, sinkID, msg, ownerID, state)
}

func (es eventStore) ViewSinkInternal(ctx context.Context, ownerID string, key string) (sinks.Sink, error) {
	return es.svc.ViewSinkInternal(ctx, ownerID, key)
}

func (es eventStore) CreateSink(ctx context.Context, token string, s sinks.Sink) (sinks.Sink, error) {
	return es.svc.CreateSink(ctx, token, s)
}

func (es eventStore) UpdateSink(ctx context.Context, token string, s sinks.Sink) (err error) {
	if err := es.svc.UpdateSink(ctx, token, s); err != nil {
		return err
	}

	var username string
	var password string
	var remoteHost string
	for k, v := range s.Config {
		switch k {
		case "username":
			username = fmt.Sprint(v)
		case "password":
			password = fmt.Sprint(v)
		case "remote_host":
			remoteHost = fmt.Sprint(v)
		}
	}

	event := updateSinkEvent{
		sinkID:     s.ID,
		owner:      s.MFOwnerID,
		username:   username,
		password:   password,
		remoteHost: remoteHost,
	}

	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}

	es.client.XAdd(ctx, record).Err()

	return nil
}

func (es eventStore) ListSinks(ctx context.Context, token string, pm sinks.PageMetadata) (sinks.Page, error) {
	return es.svc.ListSinks(ctx, token, pm)
}

func (es eventStore) ListBackends(ctx context.Context, token string) (_ []string, err error) {
	return es.svc.ListBackends(ctx, token)
}

func (es eventStore) ViewBackend(ctx context.Context, token string, key string) (_ backend.Backend, err error) {
	return es.svc.ViewBackend(ctx, token, key)
}

func (es eventStore) ViewSink(ctx context.Context, token string, key string) (_ sinks.Sink, err error) {
	return es.svc.ViewSink(ctx, token, key)
}

func (es eventStore) DeleteSink(ctx context.Context, token, id string) error {
	if err := es.svc.DeleteSink(ctx, token, id); err != nil {
		return err
	}

	event := deleteSinkEvent{
		id: id,
	}

	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.Encode(),
	}

	es.client.XAdd(ctx, record).Err()

	return nil
}

func (es eventStore) ValidateSink(ctx context.Context, token string, sink sinks.Sink) (sinks.Sink, error) {
	return es.svc.ValidateSink(ctx, token, sink)
}

// NewEventStoreMiddleware returns wrapper around sinks service that sends
// events to event store.
func NewEventStoreMiddleware(svc sinks.SinkService, client *redis.Client) sinks.SinkService {
	return eventStore{
		svc:    svc,
		client: client,
	}
}
