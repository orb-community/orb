/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

import (
	"context"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

type Sink struct {
	ID          string
	Name        types.Identifier
	MFOwnerID   string
	Type	    string
	Description string
	Config      types.Metadata
	Created     time.Time
}

// Page contains page related metadata as well as list of sinks that
// belong to this page
type Page struct {
	PageMetadata
	Sinks []Sink
}

// SinkService Sink CRUD interface
type SinkService interface {
	// CreateSink creates new data sink
	CreateSink(ctx context.Context, token string, s Sink) (Sink, error)
	// ListSinks retrieves data about sinks
	ListSinks(ctx context.Context, token string, pm PageMetadata) (Page, error)
}

type SinkRepository interface {
	// Save persists the Sink. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, sink Sink) (string, error)
	// RetrieveAll retrieves Sinks
	RetrieveAll(ctx context.Context, owner string, pm PageMetadata) (Page, error)
}
