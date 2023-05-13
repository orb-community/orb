/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otelinf

import (
	"errors"

	"github.com/ghodss/yaml"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/policies/backend"
)

var _ backend.Backend = (*otelinfBackend)(nil)

type otelinfBackend struct {
}

func (p otelinfBackend) Validate(policy types.Metadata) error {
	// todo finish validation
	return nil
}

func (p otelinfBackend) convertFromYAML(policy string) (types.Metadata, error) {
	j := collectionPolicy{}
	err := yaml.Unmarshal([]byte(policy), &j)
	if err != nil {
		return types.Metadata{}, err
	}

	if j.Config == nil || j.Kind == "" {
		return types.Metadata{}, errors.New("malformed yaml policy")
	}

	ret := types.Metadata{}

	ret["kind"] = j.Kind
	ret["config"] = j.Config

	return ret, nil
}

func (p otelinfBackend) ConvertFromFormat(format string, policy string) (types.Metadata, error) {
	switch format {
	case "yaml":
		return p.convertFromYAML(policy)
	default:
		return nil, errors.New("unsupported format")
	}
}

func (p otelinfBackend) SupportsFormat(format string) bool {
	switch format {
	case "yaml":
		return true
	}
	return false
}

func Register() bool {
	backend.Register("otelinf", &otelinfBackend{})
	return true
}
