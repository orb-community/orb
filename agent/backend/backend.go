/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package backend

import (
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/orb-community/orb/agent/policies"
	"go.uber.org/zap"
)

const (
	Unknown RunningStatus = iota
	Running
	BackendError
	AgentError
	Offline
	Waiting
)

type RunningStatus int

var runningStatusMap = [...]string{
	"unknown",
	"running",
	"backend_error",
	"agent_error",
	"offline",
	"waiting",
}

var runningStatusRevMap = map[string]RunningStatus{
	"unknown":       Unknown,
	"running":       Running,
	"backend_error": BackendError,
	"agent_error":   AgentError,
	"offline":       Offline,
	"waiting":       Waiting,
}

type State struct {
	Status            RunningStatus
	RestartCount      int64
	LastError         string
	LastRestartTS     time.Time
	LastRestartReason string
}

func (s RunningStatus) String() string {
	return runningStatusMap[s]
}

type Backend interface {
	Configure(*zap.Logger, policies.PolicyRepo, map[string]string, map[string]interface{}) error
	SetCommsClient(string, *mqtt.Client, string)
	Version() (string, error)
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
	Stop(ctx context.Context) error
	FullReset(ctx context.Context) error

	GetStartTime() time.Time
	GetCapabilities() (map[string]interface{}, error)
	GetRunningStatus() (RunningStatus, string, error)
	GetInitialState() RunningStatus

	ApplyPolicy(data policies.PolicyData, updatePolicy bool) error
	RemovePolicy(data policies.PolicyData) error
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
