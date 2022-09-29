/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/backend"
	"github.com/opentracing/opentracing-go"
)

var _ backend.Backend = (*cloudprober)(nil)

type cloudprober struct {
	auth        mainflux.AuthServiceClient
	agentRepo   fleet.AgentRepository
	Backend     string
	Description string
}

func (p cloudprober) MakeHandler(tracer opentracing.Tracer, opts []kithttp.ServerOption, r *bone.Mux) {
}

func (p cloudprober) Metadata() interface{} {
	return struct {
		Backend     string `json:"backend"`
		Description string `json:"description"`
	}{
		Backend:     p.Backend,
		Description: p.Description,
	}
}

func Register(auth mainflux.AuthServiceClient, agentRepo fleet.AgentRepository) bool {
	backend.Register("cloudprober", &cloudprober{
		Backend:   "cloudprober",
		auth:      auth,
		agentRepo: agentRepo,
	})
	return true
}
