/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package types

import (
	"encoding/json"
	"github.com/ns1labs/orb/pkg/errors"
)

// Tags A flat kv pair object
type Tags map[string]string

func (t *Tags) Append(newTags map[string]string) {
	for k, v := range newTags {
		(*t)[k] = v
	}
}

// Metadata Maybe a full object hierarchy
type Metadata map[string]interface{}

func (s *Metadata) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	}
	return errors.New("type assertion failed")
}

func (s *Metadata) RestrictKeys(predicate func(string) bool) {
	for key, _ := range *s {
		if predicate(key) {
			(*s)[key] = ""
		}
	}
}

func (s *Metadata) IsApplicable(filterFunc func(string, interface{}) bool) bool {
	for key, value := range *s {
		if filterFunc(key, value) {
			return true
		}
	}
	return false
}

func (s *Metadata) FilterMap(predicateFunc func(string) bool, mapFunc func(string, interface{}) (string, interface{})) {
	for key, value := range *s {
		if predicateFunc(key) {
			newKey, newValue := mapFunc(key, value)
			delete(*s, key)
			(*s)[newKey] = newValue
		}
	}
}
