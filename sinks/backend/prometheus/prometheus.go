/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prometheus

import (
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks/backend"
)

var _ backend.Backend = (*prometheusBackend)(nil)

type prometheusBackend struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Config      types.Metadata `json:"config"`
}

func (p prometheusBackend) Validate(config types.Metadata) error {
	return nil
}

func (p prometheusBackend) Metadata() interface{} {
	return p.Metadata()
}

func (p prometheusBackend) GetName() string {
	return p.Name
}

func (p prometheusBackend) GetDescription() string {
	return p.Description
}

func (p prometheusBackend) GetConfig() types.Metadata {
	return p.Config
}

func Register() bool {
	backend.Register("prometheus", &prometheusBackend{
		Name:        "prometheus",
		Description: "prometheus backend",
		Config:      map[string]interface{}{"title": "Remote Host", "type": "string", "name": "remote_host"},
	})
	return true
}
