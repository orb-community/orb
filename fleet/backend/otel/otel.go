/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otel

import (
	"context"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/mainflux"
	"github.com/opentracing/opentracing-go"
	"github.com/orb-community/orb/fleet"
	"github.com/orb-community/orb/fleet/backend"
	"github.com/orb-community/orb/pkg/types"
)

var _ backend.Backend = (*otelBackend)(nil)

type otelBackend struct {
	auth        mainflux.AuthServiceClient
	agentRepo   fleet.AgentRepository
	Backend     string
	Description string
}

type BackendTaps struct {
	Name             string
	InputType        string
	ConfigPredefined []string
	TotalAgents      uint64
}

func (p otelBackend) Metadata() interface{} {
	return struct {
		Backend       string `json:"backend"`
		Description   string `json:"description"`
		SchemaVersion string `json:"schema_version"`
	}{
		Backend:       p.Backend,
		Description:   p.Description,
		SchemaVersion: CurrentSchemaVersion,
	}
}

func (p otelBackend) MakeHandler(tracer opentracing.Tracer, opts []kithttp.ServerOption, r *bone.Mux) {
	MakeOtelHandler(tracer, p, opts, r)
}

func (p otelBackend) handlers() (metadata types.Metadata, err error) {
	return
}

func (p otelBackend) inputs() (metadata types.Metadata, err error) {
	return
}

func (p otelBackend) taps(ctx context.Context, ownerID string) ([]types.Metadata, error) {
	return nil, nil
}

func Register(auth mainflux.AuthServiceClient, agentRepo fleet.AgentRepository) bool {
	backend.Register("otel", &otelBackend{
		Backend:     "otel",
		Description: "OpenTelemetry configuration YAML",
		auth:        auth,
		agentRepo:   agentRepo,
	})
	return true
}
