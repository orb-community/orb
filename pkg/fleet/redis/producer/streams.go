// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/ns1labs/orb/pkg/fleet"
)

const (
	streamID  = "orb.fleet"
	streamLen = 1000
)

var _ fleet.Service = (*eventStore)(nil)

type eventStore struct {
	svc    fleet.Service
	client *redis.Client
}

func (es eventStore) CreateAgent(ctx context.Context, token string, a fleet.Agent) (fleet.Agent, error) {
	return es.svc.CreateAgent(ctx, token, a)
}

func (es eventStore) CreateSelector(ctx context.Context, token string, s fleet.Selector) (fleet.Selector, error) {
	return es.svc.CreateSelector(ctx, token, s)
}

// NewEventStoreMiddleware returns wrapper around fleet service that sends
// events to event store.
func NewEventStoreMiddleware(svc fleet.Service, client *redis.Client) fleet.Service {
	return eventStore{
		svc:    svc,
		client: client,
	}
}
