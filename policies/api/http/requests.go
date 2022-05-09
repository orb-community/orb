// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
)

const (
	maxLimitSize = 100
	maxNameSize  = 1024
	nameOrder    = "name"
	idOrder      = "id"
	ascDir       = "asc"
	descDir      = "desc"
)

type addPolicyReq struct {
	Name          string         `json:"name"`
	Backend       string         `json:"backend"`
	SchemaVersion string         `json:"schema_version"`
	Policy        types.Metadata `json:"policy,omitempty"`
	Tags          types.Tags     `json:"tags"`
	Format        string         `json:"format,omitempty"`
	PolicyData    string         `json:"policy_data,omitempty"`
	Description   string         `json:"description"`
	token         string
}

func (req *addPolicyReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.Name == "" {
		return errors.ErrMalformedEntity
	}
	if req.Backend == "" {
		return errors.ErrMalformedEntity
	}

	if req.Policy == nil {
		// passing policy data blob in the specified format
		if req.Format == "" || req.PolicyData == "" {
			return errors.ErrMalformedEntity
		}
	} else {
		// policy is in json, verified by the back ends later
		if req.Format != "" || req.PolicyData != "" {
			return errors.ErrMalformedEntity
		}
	}

	if req.SchemaVersion == "" {
		req.SchemaVersion = "1.0"
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

type updatePolicyReq struct {
	id          string
	token       string
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Tags        types.Tags     `json:"tags,omitempty"`
	Format      string         `json:"format,omitempty"`
	Policy      types.Metadata `json:"policy,omitempty"`
	PolicyData  string         `json:"policy_data,omitempty"`
}

func (req updatePolicyReq) validate() error {

	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	if req.Name == "" {
		return errors.ErrMalformedEntity
	}

	if req.Policy == nil {
		// passing policy data blob in the specified format
		if req.Format == "" || req.PolicyData == "" {
			return errors.ErrMalformedEntity
		}
	} else {
		// policy is in json, verified by the back ends later
		if req.Format != "" || req.PolicyData != "" {
			return errors.ErrMalformedEntity
		}
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return nil
}

type addDatasetReq struct {
	Name         string     `json:"name"`
	AgentGroupID string     `json:"agent_group_id"`
	PolicyID     string     `json:"agent_policy_id"`
	SinkIDs      []string   `json:"sink_ids"`
	Tags         types.Tags `json:"tags"`
	token        string
}

func (req addDatasetReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	if req.Name == "" || req.AgentGroupID == "" || req.PolicyID == "" || len(req.SinkIDs) == 0 {
		return errors.ErrMalformedEntity
	}

	_, err := types.NewIdentifier(req.Name)
	if err != nil {
		return errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return nil
}

type listResourcesReq struct {
	token        string
	pageMetadata policies.PageMetadata
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

type updateDatasetReq struct {
	Name    string `json:"name,omitempty"`
	id      string
	token   string
	Tags    types.Tags `json:"tags,omitempty"`
	SinkIDs []string   `json:"sink_ids,omitempty"`
}

func (req updateDatasetReq) validate() error {
	if req.id == "" {
		return errors.ErrMalformedEntity
	}
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}

	return nil
}

type duplicatePolicyReq struct {
	id          string
	token       string
	Name        string         `json:"name,omitempty"`
}

func (req duplicatePolicyReq) validate() error {

	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	if req.id == "" {
		return errors.ErrMalformedEntity
	}
	if req.Name != "" {
		_, err := types.NewIdentifier(req.Name)
		if err != nil {
			return errors.Wrap(errors.ErrMalformedEntity, err)
		}
	}

	return nil
}