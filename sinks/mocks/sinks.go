package mocks

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/sinks"
	"strconv"
	"strings"
	"sync"
)

var _ sinks.SinkRepository = (*sinkRepositoryMock)(nil)

type sinkRepositoryMock struct {
	mu      sync.Mutex
	counter uint64
	sinksDb map[string]sinks.Sink
}

func NewSinkRepository() sinks.SinkRepository {
	return &sinkRepositoryMock{
		sinksDb: make(map[string]sinks.Sink),
	}
}

func (s *sinkRepositoryMock) Save(ctx context.Context, sink sinks.Sink) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ID, _ := uuid.NewV4()
	sink.ID = ID.String()
	s.sinksDb[sink.ID] = sink
	return ID.String(), nil
}

func (s* sinkRepositoryMock) RetrieveAll(ctx context.Context, owner string, pm sinks.PageMetadata) (sinks.Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pm.Limit < 0 {
		return sinks.Page{}, nil
	}

	first := uint64(pm.Offset) + 1
	last := first + uint64(pm.Limit)

	var sks[]sinks.Sink

	prefix := fmt.Sprintf("%s", owner)
	for k, v := range s.sinksDb {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if strings.HasPrefix(k, prefix) && id >= first && id < last {
			sks = append(sks, v)
		}
	}

	page := sinks.Page{
		Sinks: sks,
		PageMetadata: sinks.PageMetadata{
			Total: s.counter,
			Offset: pm.Offset,
			Limit: pm.Limit,
		},
	}
	return page, nil
}
