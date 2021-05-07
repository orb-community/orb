// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package api

import (
	"github.com/ns1labs/orb/pkg/policies"
)

type addReq struct {
	token string
	name  string
}

func (req addReq) validate() error {
	if req.token == "" {
		return policies.ErrUnauthorizedAccess
	}

	if req.name == "" {
		return policies.ErrMalformedEntity
	}

	return nil
}
