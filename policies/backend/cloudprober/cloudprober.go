/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cloudprober

import (
	"errors"

	"github.com/ghodss/yaml"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies/backend"
)

var _ backend.Backend = (*cloudproberBackend)(nil)

type cloudproberBackend struct {
}

func (p cloudproberBackend) Validate(policy types.Metadata) error {
	// todo finish validation
	return nil
}

func (p cloudproberBackend) convertFromYAML(policy string) (types.Metadata, error) {
	j := collectionPolicy{}
	err := yaml.Unmarshal([]byte(policy), &j)
	if err != nil {
		return types.Metadata{}, err
	}

	if j.Probes == nil {
		return types.Metadata{}, errors.New("malformed yaml policy")
	}

	ret := types.Metadata{}

	ret["probes"] = j.Probes

	return ret, nil
}

func (p cloudproberBackend) ConvertFromFormat(format string, policy string) (types.Metadata, error) {
	switch format {
	case "yaml":
		return p.convertFromYAML(policy)
	default:
		return nil, errors.New("unsupported format")
	}
}

func (p cloudproberBackend) SupportsFormat(format string) bool {
	switch format {
	case "yaml":
		return true
	}
	return false
}

func Register() bool {
	backend.Register("cloudprober", &cloudproberBackend{})
	return true
}
