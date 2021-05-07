// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package producer

import (
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

// NewEventStoreMiddleware returns wrapper around fleet service that sends
// events to event store.
func NewEventStoreMiddleware(svc fleet.Service, client *redis.Client) fleet.Service {
	return eventStore{
		svc:    svc,
		client: client,
	}
}

func (es eventStore) Add() (fleet.Agent, error) {
	return fleet.Agent{}, nil
}
