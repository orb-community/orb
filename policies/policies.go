/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"context"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

type Policy struct {
	Name      types.Identifier
	MFOwnerID string
	Backend   string
	Policy    types.Metadata
	Created   time.Time
}

type Repository interface {
	// Save persists the Policy. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, policy Policy) error
}
