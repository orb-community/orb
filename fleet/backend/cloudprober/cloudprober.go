/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/backend"
)

var _ backend.Backend = (*cloudprober)(nil)

type cloudprober struct {
	auth        mainflux.AuthServiceClient
	agentRepo   fleet.AgentRepository
	Backend     string
	Description string
}

func Register(auth mainflux.AuthServiceClient, agentRepo fleet.AgentRepository) bool {
	backend.Register("cloudprober", &cloudprober{
		Backend:     "cloudprober",
		auth:        auth,
		agentRepo:   agentRepo,
	})
	return true
}
