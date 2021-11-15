/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prometheus

import (
	"github.com/ns1labs/orb/sinks/backend"
	"io"
)

var _ backend.Backend = (*prometheusBackend)(nil)

type prometheusBackend struct {
	apiHost     string
	apiPort     uint64
	apiUser     string
	apiPassword string
}

type SinkFeature struct {
	Backend     string                  `json:"backend"`
	Description string                  `json:"description"`
	Config      []backend.ConfigFeature `json:"config"`
}

func (p *prometheusBackend) Metadata() interface{} {
	return SinkFeature{
		Backend:     "prometheus",
		Description: "Prometheus time series database sink",
		Config:      p.CreateFeatureConfig(),
	}
}

func (p *prometheusBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string) error {
	return nil
}

func Register() bool {
	backend.Register("prometheus", &prometheusBackend{})

	return true
}

func (p *prometheusBackend) CreateFeatureConfig() []backend.ConfigFeature {
	var configs []backend.ConfigFeature

	remoteHost := backend.ConfigFeature{
		Type:     "text",
		Input:    "text",
		Title:    "Remote Host",
		Name:     "remote_host",
		Required: true,
	}
	userName := backend.ConfigFeature{
		Type:     "text",
		Input:    "text",
		Title:    "Username",
		Name:     "username",
		Required: true,
	}
	password := backend.ConfigFeature{
		Type:     "password",
		Input:    "text",
		Title:    "Password",
		Name:     "password",
		Required: true,
	}
	configs = append(configs, remoteHost, userName, password)
	return configs
}
