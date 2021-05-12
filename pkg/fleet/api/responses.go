/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/types"
	"net/http"
)

var (
	_ types.Response = (*selectorRes)(nil)
)

type selectorRes struct {
	Name    string `json:"name"`
	created bool
}

func (s selectorRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s selectorRes) Headers() map[string]string {
	return map[string]string{}
}

func (s selectorRes) Empty() bool {
	return false
}
