/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prometheus

import (
	"github.com/orb-community/orb/sinks/backend"
)

var _ backend.Backend = (*Backend)(nil)

const (
	RemoteHostURLConfigFeature = "remote_host"
	ApiTokenConfigFeature      = "api_token"
)

//type PrometheusConfigMetadata = types.Metadata

type AuthType int

const (
	BasicAuth AuthType = iota
	TokenAuth
)

type Backend struct {
	RemoteHost string `json:"remote_host"`
}

func (p *Backend) Metadata() interface{} {
	return backend.SinkFeature{
		Backend:     "prometheus",
		Description: "Prometheus time series database sink",
		Config:      p.CreateFeatureConfig(),
	}
}

func Register() bool {
	backend.Register("prometheus", &Backend{})
	return true
}
