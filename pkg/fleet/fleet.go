/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

type Agent struct {
	AgentID string
	Owner   string
	Name    string
}

type FleetRepository interface {
	// Save persists the Agent. Successful operation is indicated by non-nil
	// error response.
	Save(cfg Agent) (string, error)
}
