/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"database/sql/driver"
	_ "github.com/mattn/go-sqlite3"
)

type PolicyData struct {
	ID         string
	Name       string
	Backend    string
	Version    int32
	Data       interface{}
	State      PolicyState
	BackendErr string
}

const (
	Unknown PolicyState = iota
	Running
	FailedToApply
)

type PolicyState int

var policyStateMap = [...]string{
	"unknown",
	"running",
	"failed_to_apply",
}

var policyStateRevMap = map[string]PolicyState{
	"unknown":         Unknown,
	"running":         Running,
	"failed_to_apply": FailedToApply,
}

func (s PolicyState) String() string {
	return policyStateMap[s]
}

func (s *PolicyState) Scan(value interface{}) error {
	*s = policyStateRevMap[string(value.([]byte))]
	return nil
}
func (s PolicyState) Value() (driver.Value, error) { return s.String(), nil }
