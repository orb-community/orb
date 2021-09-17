// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"github.com/ns1labs/orb/pkg/errors"
)

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

type listAgentBackendsReq struct {
	token string
}

func (req *listAgentBackendsReq) validate() error {
	if req.token == "" {
		return errors.ErrUnauthorizedAccess
	}
	return nil
}
