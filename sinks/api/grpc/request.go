// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import (
	"github.com/ns1labs/orb/sinks"
)

type accessByIDReq struct {
	SinkID  string
	OwnerID string
}

func (req accessByIDReq) validate() error {
	if req.SinkID == "" || req.OwnerID == "" {
		return sinks.ErrMalformedEntity
	}

	return nil
}
