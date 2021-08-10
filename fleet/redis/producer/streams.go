// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"context"
	"github.com/go-redis/redis/v8"
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

func (es eventStore) ViewAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	return es.svc.ViewAgentGroupByIDInternal(ctx, groupID, ownerID)
}

func (es eventStore) ViewAgentGroupByID(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	return es.svc.ViewAgentGroupByID(ctx, groupID, ownerID)
}

func (es eventStore) ListAgentGroups(ctx context.Context, token string, pm fleet.PageMetadata) (fleet.PageAgentGroup, error) {
	return es.svc.ListAgentGroups(ctx, token, pm)
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

func (es eventStore) ValidateAgent(ctx context.Context, token string, a fleet.Agent) (fleet.Agent, error) {
	return es.svc.ValidateAgent(ctx, token, a)
}

func (es eventStore) ValidateAgentGroup(ctx context.Context, token string, ag fleet.AgentGroup) (fleet.AgentGroup, error) {
	return es.svc.ValidateAgentGroup(ctx, token, ag)
}

// NewEventStoreMiddleware returns wrapper around fleet service that sends
// events to event store.
func NewEventStoreMiddleware(svc fleet.Service, client *redis.Client) fleet.Service {
	return eventStore{
		svc:    svc,
		client: client,
	}
}
