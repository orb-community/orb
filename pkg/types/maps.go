/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package types

import (
	"encoding/json"
	"github.com/orb-community/orb/pkg/errors"
)

// Tags A flat kv pair object
type Tags map[string]string

func (t *Tags) Merge(newTags map[string]string) {
	for k, v := range newTags {
		if v == "" {
			delete(*t, k)
		} else {
			(*t)[k] = v
		}
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

func (s *Metadata) SubSet(predicate func(string, interface{}) bool) (subSet Metadata) {
	for key, value := range *s {
		if predicate(key, value) {
			subSet[key] = value
		}
	}
	return subSet
}

func (s *Metadata) RestrictKeys(predicate func(string) bool) {
	for key, _ := range *s {
		if predicate(key) {
			(*s)[key] = ""
		}
	}
}

func (s *Metadata) Merge(metadataToAdd Metadata) {
	for k, v := range metadataToAdd {
		if v == "" {
			delete(*s, k)
		} else {
			(*s)[k] = v
		}
	}
}

func (s *Metadata) RemoveKeys(keys []string) {
	for _, key := range keys {
		if _, ok := (*s)[key]; ok {
			delete(*s, key)
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
