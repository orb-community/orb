// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package consumer

type removeEvent struct {
	id string
}

type updateChannelEvent struct {
	id       string
	name     string
	metadata map[string]interface{}
}

// Connection event is either connect or disconnect event.
type disconnectEvent struct {
	thingID   string
	channelID string
}
