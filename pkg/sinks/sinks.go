/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

type Sink struct {
	SinkID string
	Owner  string
	Name   string
}

type SinksRepository interface {
	// Save persists the Sink. Successful operation is indicated by non-nil
	// error response.
	Save(cfg Sink) (string, error)
}
