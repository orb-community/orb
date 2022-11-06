/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import "github.com/etaques/orb/pkg/types"

const CurrentSchemaVersion = "1.0"

type collectionPolicy struct {
	Handlers types.Metadata `json:"handlers"`
	Input    types.Metadata `json:"input"`
	Kind     string         `json:"kind"`
}
