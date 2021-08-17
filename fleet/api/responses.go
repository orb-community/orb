/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"github.com/ns1labs/orb/pkg/types"
	"net/http"
	"time"
)

var (
	_ types.Response = (*agentGroupRes)(nil)
	_ types.Response = (*agentRes)(nil)
)

type agentGroupRes struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description,omitempty"`
	Tags           types.Tags     `json:"tags"`
	TsCreated      time.Time      `json:"ts_created,omitempty"`
	MatchingAgents types.Metadata `json:"matching_agents,omitempty"`
	created        bool
}

func (s agentGroupRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s agentGroupRes) Headers() map[string]string {
	return map[string]string{}
}

func (s agentGroupRes) Empty() bool {
	return false
}

type agentGroupsPageRes struct {
	pageRes
	AgentGroups []agentGroupRes `json:"agentGroups"`
}

func (res agentGroupsPageRes) Code() int {
	return http.StatusOK
}

func (res agentGroupsPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res agentGroupsPageRes) Empty() bool {
	return false
}

type agentRes struct {
	ID        string `json:"id"`
	Key       string `json:"key,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	Name      string `json:"name"`
	State     string `json:"state"`
	created   bool
}

func (s agentRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s agentRes) Headers() map[string]string {
	return map[string]string{}
}

func (s agentRes) Empty() bool {
	return false
}

type viewAgentRes struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	ChannelID     string         `json:"channel_id,omitempty"`
	AgentTags     types.Tags     `json:"agent_tags"`
	OrbTags       types.Tags     `json:"orb_tags"`
	TsCreated     time.Time      `json:"ts_created"`
	AgentMetadata types.Metadata `json:"agent_metadata"`
	State         string         `json:"state"`
	LastHBData    types.Metadata `json:"last_hb_data"`
	LastHB        time.Time      `json:"ts_last_hb"`
}

func (res viewAgentRes) Code() int {
	return http.StatusOK
}

func (res viewAgentRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewAgentRes) Empty() bool {
	return false
}

type agentsPageRes struct {
	pageRes
	Agents []viewAgentRes `json:"agents"`
}

func (res agentsPageRes) Code() int {
	return http.StatusOK
}

func (res agentsPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res agentsPageRes) Empty() bool {
	return false
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
	Order  string `json:"order"`
	Dir    string `json:"direction"`
}

type removeRes struct{}

func (r removeRes) Code() int {
	return http.StatusNoContent
}

func (r removeRes) Headers() map[string]string {
	return map[string]string{}
}

func (r removeRes) Empty() bool {
	return true
}

type validateAgentRes struct {
	ID        string `json:"id"`
	Key       string `json:"key,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	Name      string `json:"name"`
	State     string `json:"state"`
}

func (s validateAgentRes) Code() int {
	return http.StatusOK
}

func (s validateAgentRes) Headers() map[string]string {
	return map[string]string{}
}

func (s validateAgentRes) Empty() bool {
	return false
}
