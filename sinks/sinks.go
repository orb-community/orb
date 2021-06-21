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
	Name      types.Identifier
	MFOwnerID string
	Config    types.Metadata
	Created   time.Time
}

type SinkRepository interface {
	// Save persists the Sink. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, cfg Sink) error
}
