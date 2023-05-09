/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prometheus

import (
	"github.com/orb-community/orb/sinks/backend"
)

var _ backend.Backend = (*Backend)(nil)

const (
	NetboxHostURLConfigFeature  = "remote_host"
	NetboxApiTokenConfigFeature = "api_token"
)

//type PrometheusConfigMetadata = types.Metadata

type AuthType int

const (
	BasicAuth AuthType = iota
	TokenAuth
)

type Backend struct {
	apiHost     string
	apiPort     uint64
	apiPassword string
}

type configParseUtility struct {
	RemoteHost string  `yaml:"remote_host"`
	APIToken   *string `yaml:"api_token,omitempty"`
}

type SinkFeature struct {
	Backend     string                  `json:"backend"`
	Description string                  `json:"description"`
	Config      []backend.ConfigFeature `json:"config"`
}

func (p *Backend) Metadata() interface{} {
	return SinkFeature{
		Backend:     "diode-service",
		Description: "Netbox Service plugin for Network Discovery Tool",
		Config:      p.CreateFeatureConfig(),
	}
}

func Register() bool {
	backend.Register("diode-service", &Backend{})
	return true
}
