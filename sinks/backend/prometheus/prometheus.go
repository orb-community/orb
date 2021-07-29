/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prometheus

import (
	"fmt"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks/backend"
	"io"
	"strconv"
)

var _ backend.Backend = (*prometheusBackend)(nil)

type prometheusBackend struct {
	apiHost     string
	apiPort     uint64
	apiUser     string
	apiPassword string
}

type SinkFeature struct {
	Backend     string         `json:"backend"`
	Description string         `json:"description"`
	Config      types.Metadata `json:"config"`
}

func (p *prometheusBackend) Connect(config map[string]interface{}) error {
	for k, v := range config {
		switch k {
		case "remote_host":
			p.apiHost = fmt.Sprint(v)
		case "port":
			p.apiPort, _ = strconv.ParseUint(fmt.Sprint(v), 10, 64)
		case "username":
			p.apiUser = fmt.Sprint(v)
		case "password":
			p.apiPassword = fmt.Sprint(v)
		}
	}
	return errors.New("Error to connect to prometheus backend")
}

func (p *prometheusBackend) Metadata() interface{} {
	return SinkFeature{
		Backend:     "prometheus",
		Description: "Prometheus time series database sink",
		Config:      map[string]interface{}{"title": "Remote Host", "type": "string", "name": "remote_host"},
	}
}

func (p *prometheusBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string) error {
	return nil
}

func Register() bool {
	backend.Register("prometheus", &prometheusBackend{})

	return true
}
