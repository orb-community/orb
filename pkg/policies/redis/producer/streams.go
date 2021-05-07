// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package producer

import (
	"github.com/go-redis/redis"
	"github.com/ns1labs/orb/pkg/policies"
)

const (
	streamID  = "orb.policies"
	streamLen = 1000
)

var _ policies.Service = (*eventStore)(nil)

type eventStore struct {
	svc    policies.Service
	client *redis.Client
}

// NewEventStoreMiddleware returns wrapper around policies service that sends
// events to event store.
func NewEventStoreMiddleware(svc policies.Service, client *redis.Client) policies.Service {
	return eventStore{
		svc:    svc,
		client: client,
	}
}

func (es eventStore) Add() (policies.Policy, error) {
	return policies.Policy{}, nil
}
