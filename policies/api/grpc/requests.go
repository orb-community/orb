// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import "github.com/ns1labs/orb/policies"

type accessByIDReq struct {
	PolicyID string
	OwnerID  string
}

func (req accessByIDReq) validate() error {
	if req.PolicyID == "" || req.OwnerID == "" {
		return policies.ErrMalformedEntity
	}

	return nil
}

type accessByGroupIDReq struct {
	GroupIDs []string
	OwnerID  string
}

func (req accessByGroupIDReq) validate() error {
	if len(req.GroupIDs) == 0 || req.OwnerID == "" {
		return policies.ErrMalformedEntity
	}

	return nil
}
