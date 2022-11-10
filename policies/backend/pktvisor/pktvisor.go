/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"errors"
	"github.com/ghodss/yaml"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies/backend"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
}

func (p pktvisorBackend) Validate(policy types.Metadata) error {
	// todo finish validation
	return nil
}

func (p pktvisorBackend) convertFromYAML(policy string) (types.Metadata, error) {
	j := collectionPolicy{}
	err := yaml.Unmarshal([]byte(policy), &j)
	if err != nil {
		return types.Metadata{}, err
	}

	if j.Input == nil || j.Handlers == nil || j.Kind == "" {
		return types.Metadata{}, errors.New("malformed yaml policy")
	}

	ret := types.Metadata{}

	ret["kind"] = j.Kind
	ret["input"] = j.Input
	ret["handlers"] = j.Handlers

	return ret, nil
}

func (p pktvisorBackend) ConvertFromFormat(format string, policy string) (types.Metadata, error) {
	switch format {
	case "yaml":
		return p.convertFromYAML(policy)
	default:
		return nil, errors.New("unsupported format")
	}
}

func (p pktvisorBackend) SupportsFormat(format string) bool {
	switch format {
	case "yaml":
		return true
	}
	return false
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{})
	return true
}
