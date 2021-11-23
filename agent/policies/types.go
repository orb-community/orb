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
	Datasets   map[string]bool
	Name       string
	Backend    string
	Version    int32
	Data       interface{}
	State      PolicyState
	BackendErr string
}

func (d *PolicyData) GetDatasetIDs() []string {
	keys := make([]string, len(d.Datasets))

	i := 0
	for k := range d.Datasets {
		keys[i] = k
		i++
	}
	return keys
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
