// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
)

const (
	maxLimitSize = 100
	maxNameSize  = 1024
	nameOrder    = "name"
	idOrder      = "id"
	ascDir       = "asc"
	descDir      = "desc"
)

type addAgentGroupReq struct {
	token       string
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Tags        types.Tags `json:"tags"`
}

func (req addAgentGroupReq) validate() error {

	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	if req.Name == "" {
		return errors.ErrMalformedEntity
	}
	if len(req.Tags) == 0 {
		return errors.ErrMalformedEntity
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return nil
}

type updateAgentGroupReq struct {
	id          string
	token       string
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Tags        types.Tags `json:"tags"`
}

func (req updateAgentGroupReq) validate() error {

	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	if req.Name == "" {
		return errors.ErrMalformedEntity
	}
	if len(req.Tags) == 0 {
		return errors.ErrMalformedEntity
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return nil
}

type addAgentReq struct {
	token   string
	Name    string     `json:"name,omitempty"`
	OrbTags types.Tags `json:"orb_tags,omitempty"`
}

func (req addAgentReq) validate() error {

	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	if req.Name == "" {
		return errors.ErrMalformedEntity
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return nil
}

type updateAgentReq struct {
	id    string
	token string
	Name  string     `json:"name,omitempty"`
	Tags  types.Tags `json:"orb_tags,omitempty"`
}

func (req updateAgentReq) validate() error {

	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	if req.Name == "" {
		return errors.ErrMalformedEntity
	}
	if len(req.Tags) == 0 {
		return errors.ErrMalformedEntity
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return nil
}

type viewResourceReq struct {
	token string
	id    string
}

func (req viewResourceReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	if req.id == "" {
		return errors.ErrMalformedEntity
	}
	return nil
}

type listResourcesReq struct {
	token        string
	pageMetadata fleet.PageMetadata
}

func (req *listResourcesReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.pageMetadata.Limit == 0 {
		req.pageMetadata.Limit = defLimit
	}

	if req.pageMetadata.Limit > maxLimitSize {
		return errors.ErrMalformedEntity
	}

	if len(req.pageMetadata.Name) > maxNameSize {
		return errors.ErrMalformedEntity
	}

	if req.pageMetadata.Order != "" &&
		req.pageMetadata.Order != nameOrder && req.pageMetadata.Order != idOrder {
		return errors.ErrMalformedEntity
	}

	if req.pageMetadata.Dir != "" &&
		req.pageMetadata.Dir != ascDir && req.pageMetadata.Dir != descDir {
		return errors.ErrMalformedEntity
	}

	return nil
}

type listAgentBackendsReq struct {
	token string
}

func (req *listAgentBackendsReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	return nil
}

type agentsStatisticsReq struct {
	token string
}

func (req *agentsStatisticsReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	return nil
}