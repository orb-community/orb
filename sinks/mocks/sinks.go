// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package mocks

import (
	"context"
	"github.com/benbjohnson/immutable"
	"github.com/gofrs/uuid"
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks"
	"reflect"
	"sync"
)

var _ sinks.SinkRepository = (*sinkRepositoryMock)(nil)

// Mock Repository
type sinkRepositoryMock struct {
	mu        sync.Mutex
	counter   uint64
	sinksMock immutable.Map[string, sinks.Sink]
}

func (s *sinkRepositoryMock) SearchAllSinks(ctx context.Context, filter sinks.Filter) ([]sinks.Sink, error) {
	return nil, nil
}

func (s *sinkRepositoryMock) UpdateSinkState(ctx context.Context, sinkID string, msg string, ownerID string, state sinks.State) error {
	return nil
}

func (s *sinkRepositoryMock) RetrieveByOwnerAndId(ctx context.Context, ownerID string, key string) (sinks.Sink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sink, ok := s.sinksMock.Get(key); ok {
		if sink.MFOwnerID == ownerID {
			return sink, nil
		} else {
			return sinks.Sink{}, sinks.ErrNotFound
		}
	}

	return sinks.Sink{}, sinks.ErrNotFound
}

func NewSinkRepository() sinks.SinkRepository {
	mocks := immutable.NewMap[string, sinks.Sink](nil)
	return &sinkRepositoryMock{
		sinksMock: *mocks,
	}
}

func (s *sinkRepositoryMock) Save(ctx context.Context, sink sinks.Sink) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	itr := s.sinksMock.Iterator()
	for !itr.Done() {
		_, v, _ := itr.Next()
		if v.Name == sink.Name {
			return "", sinks.ErrConflictSink
		}
	}

	s.counter++
	ID, _ := uuid.NewV4()
	sink.ID = ID.String()
	// create a full copy of the Config, because somehow it changes after adding to map
	configCopy := make(types.Metadata)
	bkpConfig := sink.Config
	exporterMeta := sink.Config.GetSubMetadata("exporter")
	authMeta := sink.Config.GetSubMetadata("authentication")
	configCopy["exporter"] = types.FromMap(exporterMeta)
	configCopy["authentication"] = types.FromMap(authMeta)
	configCopy["opentelemetry"] = "enabled"
	sink.Config = configCopy
	s.sinksMock = *s.sinksMock.Set(ID.String(), sink)
	sink.Config = bkpConfig
	return sink.ID, nil
}

func (s *sinkRepositoryMock) Update(ctx context.Context, sink sinks.Sink) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.sinksMock.Get(sink.ID); ok {
		if sink.MFOwnerID != c.MFOwnerID {
			return errors.ErrUpdateEntity
		}
		// create a full copy of the Config, because somehow it changes after adding to map
		configCopy := make(types.Metadata)
		bkpConfig := sink.Config
		exporterMeta := sink.Config.GetSubMetadata("exporter")
		authMeta := sink.Config.GetSubMetadata("authentication")
		configCopy["exporter"] = types.FromMap(exporterMeta)
		configCopy["authentication"] = types.FromMap(authMeta)
		configCopy["opentelemetry"] = "enabled"
		sink.Config = configCopy
		s.sinksMock = *s.sinksMock.Set(sink.ID, sink)
		sink.Config = bkpConfig
		return nil
	}
	return sinks.ErrNotFound
}

func (s *sinkRepositoryMock) RetrieveAllByOwnerID(ctx context.Context, owner string, pm sinks.PageMetadata) (sinks.Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	first := uint64(pm.Offset) + 1
	last := first + pm.Limit

	var sks []sinks.Sink

	id := uint64(0)
	itr := s.sinksMock.Iterator()
	for !itr.Done() {
		_, v, _ := itr.Next()
		id++
		if v.MFOwnerID == owner && id >= first && id < last {
			if reflect.DeepEqual(pm.Tags, v.Tags) || pm.Tags == nil {
				sks = append(sks, v)
			}
		}
	}

	sks = sortSinks(pm, sks)

	page := sinks.Page{
		Sinks: sks,
		PageMetadata: sinks.PageMetadata{
			Total:  s.counter,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}
	return page, nil
}

func (s *sinkRepositoryMock) RetrieveById(ctx context.Context, key string) (sinks.Sink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.sinksMock.Get(key); ok {
		return c, nil
	}

	return sinks.Sink{}, sinks.ErrNotFound
}

func (s *sinkRepositoryMock) Remove(ctx context.Context, owner string, key string) error {
	if c, ok := s.sinksMock.Get(key); ok {
		if c.MFOwnerID == owner {
			s.sinksMock = *s.sinksMock.Delete(key)
			return nil
		} else {
			return sinks.ErrNotFound
		}
	}
	return nil
}
