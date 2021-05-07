/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package writer

// Service represents a fleet service that can consume agent output and maintain fleet information
type Service interface {
	// Begin consuming agent output
	Run() error
}
