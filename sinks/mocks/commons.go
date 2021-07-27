// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package mocks

import (
	"github.com/ns1labs/orb/sinks"
	"sort"
)

func sortSinks(pm sinks.PageMetadata, sks []sinks.Sink) []sinks.Sink {
	switch pm.Order {
	case "name":
		if pm.Dir == "asc" {
			sort.SliceStable(sks, func(i, j int) bool {
				return sks[i].Name.String() < sks[j].Name.String()
			})
		}
		if pm.Dir == "desc" {
			sort.SliceStable(sks, func(i, j int) bool {
				return sks[i].Name.String() > sks[j].Name.String()
			})
		}
	case "id":
		if pm.Dir == "asc" {
			sort.SliceStable(sks, func(i, j int) bool {
				return sks[i].ID < sks[j].ID
			})
		}
		if pm.Dir == "desc" {
			sort.SliceStable(sks, func(i, j int) bool {
				return sks[i].ID > sks[j].ID
			})
		}
	default:
		sort.SliceStable(sks, func(i, j int) bool {
			return sks[i].ID < sks[j].ID
		})
	}

	return sks
}
