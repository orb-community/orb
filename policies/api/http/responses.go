/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import "net/http"

type policyRes struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Backend string `json:"backend"`
	created bool
}

func (s policyRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s policyRes) Headers() map[string]string {
	return map[string]string{}
}

func (s policyRes) Empty() bool {
	return false
}

type policiesPageRes struct {
	pageRes
	Policies []policyRes `json:"data"`
}

func (res policiesPageRes) Code() int {
	return http.StatusOK
}

func (res policiesPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res policiesPageRes) Empty() bool {
	return false
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
	Order  string `json:"order"`
	Dir    string `json:"direction"`
}

type datasetRes struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	created bool
}
