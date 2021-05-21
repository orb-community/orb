// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/fleet"
	"github.com/ns1labs/orb/pkg/types"
)

type addSelectorReq struct {
	token    string
	Name     string                 `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (req addSelectorReq) validate() error {

	if req.token == "" {
		return fleet.ErrUnauthorizedAccess
	}
	if req.Name == "" {
		return fleet.ErrMalformedEntity
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	return nil
}

type addAgentReq struct {
	token   string
	Name    string                 `json:"name,omitempty"`
	OrbTags map[string]interface{} `json:"orb_tags,omitempty"`
}

func (req addAgentReq) validate() error {

	if req.token == "" {
		return fleet.ErrUnauthorizedAccess
	}
	if req.Name == "" {
		return fleet.ErrMalformedEntity
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	return nil
}
