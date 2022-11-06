/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"github.com/etaques/orb/pkg/types"
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
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	State         string         `json:"state"`
	Key           string         `json:"key,omitempty"`
	ChannelID     string         `json:"channel_id,omitempty"`
	AgentTags     types.Tags     `json:"agent_tags"`
	OrbTags       types.Tags     `json:"orb_tags"`
	AgentMetadata types.Metadata `json:"agent_metadata"`
	LastHBData    types.Metadata `json:"last_hb_data"`
	TsCreated     time.Time      `json:"ts_created"`
	TsLastHB      time.Time      `json:"ts_last_hb"`
	PolicyState   types.Metadata `json:"policy_state,omitempty"`
	created       bool
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

type agentsPageRes struct {
	pageRes
	Agents []agentRes `json:"agents"`
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

type validateAgentGroupRes struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name"`
	Description    string         `json:"description,omitempty"`
	Tags           types.Tags     `json:"tags"`
	MatchingAgents types.Metadata `json:"matching_agents,omitempty"`
}

func (s validateAgentGroupRes) Code() int {
	return http.StatusOK
}

func (s validateAgentGroupRes) Headers() map[string]string {
	return map[string]string{}
}

func (s validateAgentGroupRes) Empty() bool {
	return false
}

type validateAgentRes struct {
	Name    string     `json:"name"`
	OrbTags types.Tags `json:"orb_tags"`
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

type agentBackendsRes struct {
	Backends []interface{} `json:"backends,omitempty"`
}

func (s agentBackendsRes) Code() int {
	return http.StatusOK
}

func (s agentBackendsRes) Headers() map[string]string {
	return map[string]string{}
}

func (s agentBackendsRes) Empty() bool {
	return false
}

type matchingGroupsRes struct {
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

func (s matchingGroupsRes) Code() int {
	return http.StatusOK
}

func (s matchingGroupsRes) Headers() map[string]string {
	return map[string]string{}
}

func (s matchingGroupsRes) Empty() bool {
	return false
}
