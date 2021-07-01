/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import "github.com/ns1labs/orb/pkg/types"

type sinkRes struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Config  types.Metadata `json:"config,omitempty"`
	created bool
}
