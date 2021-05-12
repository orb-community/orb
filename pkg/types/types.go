/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package types

import (
	"database/sql/driver"
	"errors"
	"regexp"
)

type Identifier struct {
	string
}

const (
	idRegexp = "^[a-zA-Z_][a-zA-Z0-9_-]*$"
)

var (
	idRegexpC, _ = regexp.Compile(idRegexp)

	// ErrBadIdentifier indicates a bad identifer pattern
	ErrBadIdentifier = errors.New("invalid identifier, must match " + idRegexp)
)

const maxIdentifierLength = 64

func NewIdentifier(v string) (*Identifier, error) {
	var i Identifier
	i.string = v
	if !i.IsValid() {
		return nil, ErrBadIdentifier
	}
	return &i, nil
}

func (i *Identifier) IsValid() bool {
	if len(i.string) < 2 || len(i.string) > maxIdentifierLength {
		return false
	}
	if !idRegexpC.MatchString(i.string) {
		return false
	}
	return true
}

func (i *Identifier) String() string {
	return i.string
}

// Scan - Implement the database/sql scanner interface
func (i *Identifier) Scan(value interface{}) error {
	if value == nil {
		return ErrBadIdentifier
	}

	b, ok := value.(string)
	if !ok {
		return ErrBadIdentifier
	}
	i.string = b

	if !i.IsValid() {
		return ErrBadIdentifier
	}

	return nil
}

// Value Implements valuer
func (i Identifier) Value() (driver.Value, error) {
	return i.string, nil
}
