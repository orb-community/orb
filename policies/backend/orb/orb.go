/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package orb

import (
	"github.com/ns1labs/orb/policies/backend"
)

var _ backend.Backend = (*orbBackend)(nil)

type orbBackend struct {
}

func (p orbBackend) SupportsFormat(format string) bool {
	switch format {
	case "yaml":
		return true
	}
	return false
}

func Register() bool {
	backend.Register("orb", &orbBackend{})
	return true
}
