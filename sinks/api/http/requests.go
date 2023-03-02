// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks"
	"github.com/orb-community/orb/sinks/backend"
)

const (
	maxLimitSize = 100
	maxNameSize  = 1024
	nameOrder    = "name"
	idOrder      = "id"
	ascDir       = "asc"
	descDir      = "desc"
)

type addReq struct {
	Name        string         `json:"name,omitempty"`
	Backend     string         `json:"backend,omitempty"`
	Config      types.Metadata `json:"config,omitempty"`
	Format      string         `json:"format,omitempty"`
	ConfigData  string         `json:"config_data,omitempty"`
	Description string         `json:"description,omitempty"`
	Tags        types.Tags     `json:"tags,omitempty"`
	token       string
}

func (req addReq) validate() (err error) {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.Backend == "" || !backend.HaveBackend(req.Backend) {
		return errors.Wrap(errors.ErrMalformedEntity, errors.New("backend not found"))
	}

	reqBackend := backend.GetBackend(req.Backend)
	if req.ConfigData == "" || req.Config != nil {
		return errors.Wrap(errors.ErrMalformedEntity, errors.New("config not found"))
	}

	var config types.Metadata
	if req.Format != "" {
		config, err = reqBackend.ParseConfig(req.Format, req.ConfigData)
		if err != nil {
			return errors.Wrap(errors.ErrMalformedEntity, err)
		}
	} else {
		config = req.Config
	}

	err = reqBackend.ValidateConfiguration(config)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	if req.Name == "" {
		return errors.Wrap(errors.ErrMalformedEntity, errors.New("name not found"))
	}

	_, err = types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return nil
}

type updateSinkReq struct {
	Name        string         `json:"name,omitempty"`
	Config      types.Metadata `json:"config,omitempty"`
	Backend     string         `json:"backend,omitempty"`
	Format      string         `json:"format,omitempty"`
	ConfigData  string         `json:"config_data,omitempty"`
	Description *string        `json:"description,omitempty"`
	Tags        types.Tags     `json:"tags,omitempty"`
	id          string
	token       string
}

func (req updateSinkReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return errors.ErrMalformedEntity
	}

	if req.Backend == "" || !backend.HaveBackend(req.Backend) {
		return errors.ErrMalformedEntity
	}

	reqBackend := backend.GetBackend(req.Backend)
	if req.ConfigData == "" && req.Config == nil {
		return errors.Wrap(errors.ErrMalformedEntity, errors.New("config not found"))
	}

	var config types.Metadata
	var err error
	if req.Format != "" {
		config, err = reqBackend.ParseConfig(req.Format, req.ConfigData)
		if err != nil {
			return errors.Wrap(errors.ErrMalformedEntity, err)
		}
	} else {
		config = req.Config
	}

	err = reqBackend.ValidateConfiguration(config)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	if req.Description == nil && req.Name == "" && len(req.Config) == 0 && req.Tags == nil {
		return errors.ErrMalformedEntity
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
	pageMetadata sinks.PageMetadata
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

type listBackendsReq struct {
	token string
}

func (req *listBackendsReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	return nil
}

type deleteSinkReq struct {
	token string
	id    string
}

func (req deleteSinkReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return errors.ErrMalformedEntity
	}

	return nil
}

type validateReq struct {
	Name        string         `json:"name,omitempty"`
	Backend     string         `json:"backend,omitempty"`
	Config      types.Metadata `json:"config,omitempty"`
	Description string         `json:"description,omitempty"`
	Tags        types.Tags     `json:"tags,omitempty"`
	token       string
}

func (req validateReq) validate() error {
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
