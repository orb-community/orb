// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package orb

import (
	"encoding/json"
	"net/http"
)

const version string = "0.9.0-develop"
const minAgentVersion string = "0.9.0-develop"

func GetVersion() string {
	return version
}

func GetMinAgentVersion() string {
	return minAgentVersion
}

// VersionInfo contains version endpoint response.
type VersionInfo struct {
	// Service contains service name.
	Service string `json:"service"`

	// Version contains service current version value.
	Version string `json:"version"`
}

// Version exposes an HTTP handler for retrieving service version.
func Version(service string) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		res := VersionInfo{service, version}

		data, _ := json.Marshal(res)

		rw.Write(data)
	})
}
