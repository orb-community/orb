// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package consumer

import (
	"time"
)

type sinksUpdateEvent struct {
	ownerID   string
	sinkID    string
	config    string
	timestamp time.Time
}

type sinkerUpdateEvent struct {
	ownerID   string
	sinkID    string
	state     string
	timestamp time.Time
}
