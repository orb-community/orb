/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"github.com/ns1labs/orb/policies/backend"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
}

func (p pktvisorBackend) SupportsFormat(format string) bool {
	switch format {
	case "yaml":
		return true
	}
	return false
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{})
	return true
}
