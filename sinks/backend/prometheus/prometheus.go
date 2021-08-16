/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prometheus

import (
	"fmt"
	"github.com/ns1labs/orb/pkg/errors"
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
	Backend     string          `json:"backend"`
	Description string          `json:"description"`
	Config      []configFeature `json:"config"`
}

type configFeature struct {
	Type     string `json:"type"`
	Input    string `json:"input"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

func (p *prometheusBackend) Connect(config map[string]interface{}) error {
	if len(config) != 0 {
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
	}
	return errors.New("Error to connect to prometheus backend")
}

func (p *prometheusBackend) Metadata() interface{} {
	return SinkFeature{
		Backend:     "prometheus",
		Description: "Prometheus time series database sink",
		Config:      createFeatureConfig(),
	}
}

func (p *prometheusBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string) error {
	return nil
}

func Register() bool {
	backend.Register("prometheus", &prometheusBackend{})

	return true
}

func createFeatureConfig() []configFeature {
	var configs []configFeature

	remoteHost := configFeature{
		Type:     "text",
		Input:    "text",
		Title:    "Remote Host",
		Name:     "remote_host",
		Required: true,
	}
	userName := configFeature{
		Type:     "text",
		Input:    "text",
		Title:    "Username",
		Name:     "username",
		Required: true,
	}
	password := configFeature{
		Type:     "password",
		Input:    "text",
		Title:    "Password",
		Name:     "password",
		Required: true,
	}
	configs = append(configs, remoteHost, userName, password)
	return configs
}
