/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

import (
	"database/sql/driver"
)

type SinkConfig struct {
	SinkID   string          `json:"sink_id"`
	OwnerID  string          `json:"owner_id"`
	Url      string          `json:"remote_host"`
	User     string          `json:"username"`
	Password string          `json:"password"`
	State    PrometheusState `json:"state"`
}

const (
	Unknown PrometheusState = iota
	Connected
	FailedToConnect
)

type PrometheusState int

var promStateMap = [...]string{
	"unknown",
	"connected",
	"failed_to_connect",
}

var promStateRevMap = map[string]PrometheusState{
	"unknown":           Unknown,
	"connected":         Connected,
	"failed_to_connect": FailedToConnect,
}

func (p PrometheusState) String() string {
	return promStateMap[p]
}

func (p *PrometheusState) Scan(value interface{}) error {
	*p = promStateRevMap[string(value.([]byte))]
	return nil
}

func (p PrometheusState) Value() (driver.Value, error) {
	return p.String(), nil
}
