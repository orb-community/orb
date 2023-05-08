// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package errors

var (
	// ErrUnsupportedContentType indicates unacceptable or lack of Content-Type
	ErrUnsupportedContentType = New("unsupported content type")

	// ErrInvalidQueryParams indicates invalid query parameters
	ErrInvalidQueryParams = New("invalid query parameters")

	// ErrNotFoundParam indicates that the parameter was not found in the query
	ErrNotFoundParam = New("parameter not found in the query")

	// ErrMalformedEntity indicates a malformed entity specification
	ErrMalformedEntity = New("malformed entity specification")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = New("non-existent entity")

	// ErrConflict indicates that entity already exists.
	ErrConflict = New("entity already exists")

	// ErrUpdateEntity indicates error in updating entity or entities
	ErrUpdateEntity = New("update entity failed")

	// ErrViewEntity indicates error in viewing entity or entities
	ErrViewEntity = New("view entity failed")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = New("missing or invalid credentials provided")

	// ErrScanMetadata indicates problem with metadata in db.
	ErrScanMetadata = New("failed to scan metadata")

	// ErrSelectEntity indicates error while reading entity from database
	ErrSelectEntity = New("select entity from db error")
)
