/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"database/sql/driver"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type PolicyData struct {
	ID                 string
	Datasets           map[string]bool
	GroupIds           map[string]bool
	Name               string
	Backend            string
	Version            int32
	Data               interface{}
	State              PolicyState
	BackendErr         string
	LastScrapeBytes    int64
	LastScrapeTS       time.Time
	PreviousPolicyData *PolicyData
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
	Offline
	NoTapMatch
)

type PolicyState int

var policyStateMap = [...]string{
	"unknown",
	"running",
	"failed_to_apply",
	"offline",
	"no_tap_match",
}

var policyStateRevMap = map[string]PolicyState{
	"unknown":         Unknown,
	"running":         Running,
	"failed_to_apply": FailedToApply,
	"offline":         Offline,
	"no_tap_match":    NoTapMatch,
}

func (s PolicyState) String() string {
	return policyStateMap[s]
}

func (s *PolicyState) Scan(value interface{}) error {
	*s = policyStateRevMap[string(value.([]byte))]
	return nil
}
func (s PolicyState) Value() (driver.Value, error) { return s.String(), nil }
