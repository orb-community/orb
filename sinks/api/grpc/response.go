// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

type sinkRes struct {
	id          string
	mfOwnerId   string
	name        string
	description string
	tags        []byte
	state       string
	error       string
	backend     string
	config      []byte
}

type sinksRes struct {
	sinks []sinkRes
}
