/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package backend

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ns1labs/orb/agent/policies"
	"go.uber.org/zap"
)

const (
	Unknown BackendState = iota
	Running
	BackendError
	AgentError
)

type BackendState int

var backendStateMap = [...]string{
	"unknown",
	"running",
	"backend_error",
	"agent_error",
}

var backendStateRevMap = map[string]BackendState{
	"unknown":       Unknown,
	"running":       Running,
	"backend_error": BackendError,
	"agent_error":   AgentError,
}

func (s BackendState) String() string {
	return backendStateMap[s]
}

type Backend interface {
	Configure(*zap.Logger, policies.PolicyRepo, map[string]string) error
	SetCommsClient(string, mqtt.Client, string)
	Version() (string, error)
	Start() error
	Stop() error

	GetCapabilities() (map[string]interface{}, error)
	GetState() (BackendState, string, error)

	ApplyPolicy(data policies.PolicyData) error
	RemovePolicy(policyID string) error
}

var registry = make(map[string]Backend)

func Register(name string, b Backend) {
	registry[name] = b
}

func GetList() []string {
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}
	return keys
}

func HaveBackend(name string) bool {
	_, prs := registry[name]
	return prs
}

func GetBackend(name string) Backend {
	return registry[name]
}
