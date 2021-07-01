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
	"github.com/ns1labs/orb/fleet"
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

func (es eventStore) ListAgents(ctx context.Context, token string, pm fleet.PageMetadata) (fleet.Page, error) {
	return es.svc.ListAgents(ctx, token, pm)
}

func (es eventStore) CreateAgent(ctx context.Context, token string, a fleet.Agent) (fleet.Agent, error) {
	return es.svc.CreateAgent(ctx, token, a)
}

func (es eventStore) CreateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (fleet.AgentGroup, error) {
	return es.svc.CreateAgentGroup(ctx, token, s)
}

// NewEventStoreMiddleware returns wrapper around fleet service that sends
// events to event store.
func NewEventStoreMiddleware(svc fleet.Service, client *redis.Client) fleet.Service {
	return eventStore{
		svc:    svc,
		client: client,
	}
}
