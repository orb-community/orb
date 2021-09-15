/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package backend

import "go.uber.org/zap"

const (
	Unknown State = iota
	Running
	BackendError
	AgentError
)

type State int

var stateMap = [...]string{
	"unknown",
	"running",
	"backend_error",
	"agent_error",
}

var stateRevMap = map[string]State{
	"unknown":       Unknown,
	"running":       Running,
	"backend_error": BackendError,
	"agent_error":   AgentError,
}

func (s State) String() string {
	return stateMap[s]
}

type Backend interface {
	Configure(*zap.Logger, map[string]string) error
	Version() (string, error)
	Start() error
	Stop() error

	GetCapabilities() (map[string]interface{}, error)
	GetState() (State, string, error)

	ApplyPolicy(policyID string, policyData interface{}) error
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
