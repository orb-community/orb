// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package db

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/ns1labs/orb/pkg/errors"
)

const (
	ErrDuplicate  = "unique_violation"
	ErrInvalid    = "invalid_text_representation"
	ErrTruncation = "string_data_right_truncation"
)

var (
	ErrSaveDB    = errors.New("Failed to save to database")
	ErrUpdateDB  = errors.New("Failed to update database")
	ErrMarshal   = errors.New("Failed to marshal metadata")
	ErrUnmarshal = errors.New("Failed to unmarshal metadata")
	// ErrScanMetadata indicates problem with metadata in db.
	ErrScanMetadata = errors.New("failed to scan metadata")
)

// Metadata type for handling metadata properly in database/sql
type Metadata map[string]interface{}
type Tags map[string]string

// Scan - Implement the database/sql scanner interface
func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return ErrScanMetadata
	}

	if err := json.Unmarshal(b, m); err != nil {
		return err
	}

	return nil
}

// Value Implements valuer
func (m Metadata) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "{}", nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (m *Tags) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return ErrScanMetadata
	}

	if err := json.Unmarshal(b, m); err != nil {
		return err
	}

	return nil
}

// Value Implements valuer
func (m Tags) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "{}", nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}
