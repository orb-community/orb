// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package consumer

import (
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

type updateSinkEvent struct {
	sinkID    string
	owner     string
	config    types.Metadata
	state     string
	timestamp time.Time
}
