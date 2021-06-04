/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"github.com/ns1labs/orb/agent/backend"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{})
	return true
}
