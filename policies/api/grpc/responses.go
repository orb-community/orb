// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

type policyRes struct {
	id      string
	name    string
	backend string
	version int32
	data    []byte
}

type policyInDSRes struct {
	id        string
	name      string
	backend   string
	version   int32
	data      []byte
	datasetID string
}

type policyInDSListRes struct {
	policies []policyInDSRes
}

type emptyRes struct {
	err error
}
