/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"time"
)

type Selector struct {
	Owner   string
	Name    string
	Config  Metadata
	Created time.Time
}

type SelectorRepository interface {
	// Save persists the Selector. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, cfg Selector) (string, error)
}
