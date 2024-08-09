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

	// 	ErrEntityNameNotFound indicates that the entity name was not found
	ErrEntityNameNotFound = New("malformed entity specificiation. name not found")

	// ErrBackendNotFound indicates that the backend field was not found
	ErrBackendNotFound = New("malformed entity specification. backend field is expected")

	// ErrInvalidBackend indicates a malformed entity specification on backend field
	ErrInvalidBackend = New("malformed entity specification. backend field is invalid")

	// ErrConfigFieldNotFound indicates that configuration field was not found
	ErrConfigFieldNotFound = New("malformed entity specification. configuration field is expected")

	// ErrExporterFieldNotFound indicates that exporter field was not found
	ErrExporterFieldNotFound = New("malformed entity specification. exporter field is expected on configuration field")

	// ErrEndPointNotFound indicates that endpoint field was not found on exporter field for otlp backend
	ErrEndpointNotFound = New("malformed entity specification. endpoint field is expected on exporter field")

	// ErrInvalidEndpoint indicates that endpoint field is not valid
	ErrInvalidEndpoint = New("malformed entity specification. endpoint field is invalid")

	// ErrAuthFieldNotFound indicates that authentication field was not found on configuration field
	ErrAuthFieldNotFound = New("malformed entity specification. authentication fields are expected on configuration field")

	// ErrAuthTypeNotFound indicates that authentication type field was not found on the authentication field
	ErrAuthTypeNotFound = New("malformed entity specification: authentication type field is expected on configuration field")

	// ErrAuthInvalidType indicates invalid authentication type
	ErrAuthInvalidType = New("malformed entity specification. type key on authentication field is invalid")

	// ErrAuthUsernameNotFound indicates that username key was not found
	ErrAuthUsernameNotFound = New("malformed entity specification. username key is expected on authentication field")

	// ErrAuthPasswordNotFound indicates that password key was not found
	ErrAuthPasswordNotFound = New("malformed entity specification. password key is expected on authentication field")

	// ErrAuthSchemeNotFound indicates that scheme key was not found
	ErrAuthSchemeNotFound = New("malformed entity specification. scheme key is expected on authentication field")

	// ErrAuthTokenNotFound indicates that token key was not found
	ErrAuthTokenNotFound = New("malformed entity specification. token key is expected on authentication field")

	// ErrAuthInvalidPasswordType indicates invalid password key on authentication field
	ErrAuthInvalidPasswordType = New("malformed entity specification. password key on authentication field is invalid")

	// ErrAuthInvalidSchemeType indicates invalid scheme key on authentication field
	ErrAuthInvalidSchemeType = New("malformed entity specification. scheme key on authentication field is invalid")

	// ErrAuthInvalidTokenType indicates invalid token key on authentication field
	ErrAuthInvalidTokenType = New("malformed entity specification. token key on authentication field is invalid")

	// ErrAuthInvalidUsernameType indicates invalid username key on authentication field
	ErrAuthInvalidUsernameType = New("malformed entity specification. username key on authentication field is invalid")

	// ErrRemoteHostNotFound indicates that remote host field was not found
	ErrRemoteHostNotFound = New("malformed entity specification. remote host is expected on exporter field")

	// ErrInvalidRemoteHost indicates that remote host field is invalid
	ErrInvalidRemoteHost = New("malformed entity specification. remote host type is invalid")

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
