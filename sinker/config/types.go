/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

import (
	"database/sql/driver"
)

type SinkConfig struct {
	SinkID   string
	OwnerID  string
	Url      string
	User     string
	Password string
	State    PrometheusState
}

const (
	Connected PrometheusState = iota
	FailedToConnect
)

type PrometheusState int

var promStateMap = [...]string{
	"connected",
	"failed_to_connect",
}

var promStateRevMap = map[string]PrometheusState{
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

func (p PrometheusState) Valeu() (driver.Value, error) {
	return p.String(), nil
}
