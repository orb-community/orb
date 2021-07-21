/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/types"
	"net/http"
	"time"
)

type sinkRes struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Tags        types.Tags     `json:"tags"`
	Backend     string         `json:"backend"`
	Config      types.Metadata `json:"config,omitempty"`
	TsCreated   time.Time      `json:"ts_created"`
	created     bool
}

func (s sinkRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s sinkRes) Headers() map[string]string {
	return map[string]string{}
}

func (s sinkRes) Empty() bool {
	return false
}

type sinksPagesRes struct {
	pageRes
	Sinks []sinkRes `json:"sinks"`
}

func (res sinksPagesRes) Code() int {
	return http.StatusOK
}

func (res sinksPagesRes) Headers() map[string]string {
	return map[string]string{}
}

func (res sinksPagesRes) Empty() bool {
	return false
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
	Order  string `json:"order"`
	Dir    string `json:"direction"`
}
