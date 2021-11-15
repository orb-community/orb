// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package mocks

import (
	"github.com/ns1labs/orb/policies"
	"sort"
)

func sortPolicies(pm policies.PageMetadata, ags []policies.Policy) []policies.Policy {
	switch pm.Order {
	case "name":
		if pm.Dir == "asc" {
			sort.SliceStable(ags, func(i, j int) bool {
				return ags[i].Name.String() < ags[j].Name.String()
			})
		}
		if pm.Dir == "desc" {
			sort.SliceStable(ags, func(i, j int) bool {
				return ags[i].Name.String() > ags[j].Name.String()
			})
		}
	case "id":
		if pm.Dir == "asc" {
			sort.SliceStable(ags, func(i, j int) bool {
				return ags[i].ID < ags[j].ID
			})
		}
		if pm.Dir == "desc" {
			sort.SliceStable(ags, func(i, j int) bool {
				return ags[i].ID > ags[j].ID
			})
		}
	default:
		sort.SliceStable(ags, func(i, j int) bool {
			return ags[i].ID < ags[j].ID
		})
	}

	return ags
}

func sortDataset(pm policies.PageMetadata, ags []policies.Dataset) []policies.Dataset {
	switch pm.Order {
	case "name":
		if pm.Dir == "asc" {
			sort.SliceStable(ags, func(i, j int) bool {
				return ags[i].Name.String() < ags[j].Name.String()
			})
		}
		if pm.Dir == "desc" {
			sort.SliceStable(ags, func(i, j int) bool {
				return ags[i].Name.String() > ags[j].Name.String()
			})
		}
	case "id":
		if pm.Dir == "asc" {
			sort.SliceStable(ags, func(i, j int) bool {
				return ags[i].ID < ags[j].ID
			})
		}
		if pm.Dir == "desc" {
			sort.SliceStable(ags, func(i, j int) bool {
				return ags[i].ID > ags[j].ID
			})
		}
	default:
		sort.SliceStable(ags, func(i, j int) bool {
			return ags[i].ID < ags[j].ID
		})
	}

	return ags
}
