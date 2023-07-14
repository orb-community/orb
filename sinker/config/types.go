/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

import (
	"database/sql/driver"
	"github.com/orb-community/orb/pkg/types"
	"time"
)

// SinkConfigParser to be compatible with new sinks config is coming from eventbus
type SinkConfig struct {
	SinkID          string          `json:"sink_id"`
	OwnerID         string          `json:"owner_id"`
	Config          types.Metadata  `json:"config"`
	State           PrometheusState `json:"state,omitempty"`
	Msg             string          `json:"msg,omitempty"`
	LastRemoteWrite time.Time       `json:"last_remote_write,omitempty"`
}

const (
	Unknown PrometheusState = iota
	Active
	Error
	Idle
	Warning
)

type PrometheusState int

var promStateMap = [...]string{
	"unknown",
	"active",
	"error",
	"idle",
	"warning",
}

var promStateRevMap = map[string]PrometheusState{
	"unknown": Unknown,
	"active":  Active,
	"error":   Error,
	"idle":    Idle,
	"warning": Warning,
}

func (p PrometheusState) String() string {
	return promStateMap[p]
}

func (p *PrometheusState) SetFromString(value string) error {
	*p = promStateRevMap[value]
	return nil
}

func (p PrometheusState) Value() (driver.Value, error) {
	return p.String(), nil
}
