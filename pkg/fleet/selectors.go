/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

type Selector struct {
	MFOwnerID string
	Name      types.Identifier
	Metadata  Metadata
	Created   time.Time
}

type SelectorService interface {
	// CreateSelector creates new Selector
	CreateSelector(ctx context.Context, token string, s Selector) (Selector, error)
}

type SelectorRepository interface {
	// Save persists the Selector. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, selector Selector) error
}
