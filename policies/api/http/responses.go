/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"github.com/ns1labs/orb/pkg/types"
	"net/http"
)

type policyRes struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Tags        types.Tags     `json:"tags"`
	Backend     string         `json:"backend"`
	Policy      types.Metadata `json:"policy,omitempty"`
	Format      string         `json:"format,omitempty"`
	PolicyData  string         `json:"policy_data,omitempty"`
	created     bool
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

type policyUpdateRes struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Tags        types.Tags     `json:"tags,omitempty"`
	Policy      types.Metadata `json:"policy,omitempty"`
}

func (s policyUpdateRes) Code() int {
	return http.StatusOK
}

func (s policyUpdateRes) Headers() map[string]string {
	return map[string]string{}
}

func (s policyUpdateRes) Empty() bool {
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

type removeRes struct{}

func (res removeRes) Code() int {
	return http.StatusNoContent
}

func (res removeRes) Headers() map[string]string {
	return map[string]string{}
}

func (res removeRes) Empty() bool {
	return true
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

func (s datasetRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s datasetRes) Headers() map[string]string {
	return map[string]string{}
}

func (s datasetRes) Empty() bool {
	return false
}

type policyValidateRes struct {
	Name        string         `json:"name"`
	Backend     string         `json:"backend"`
	Description string         `json:"description"`
	Tags        types.Tags     `json:"tags"`
	Policy      types.Metadata `json:"policy"`
}

func (s policyValidateRes) Code() int {
	return http.StatusOK
}

func (s policyValidateRes) Headers() map[string]string {
	return map[string]string{}
}

func (s policyValidateRes) Empty() bool {
	return false
}
