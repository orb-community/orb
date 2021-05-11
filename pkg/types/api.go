// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package types

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	contentType = "application/json"
)

// Response contains HTTP response specific methods.
type Response interface {
	// Code returns HTTP response code.
	Code() int

	// Headers returns map of HTTP headers with their values.
	Headers() map[string]string

	// Empty indicates if HTTP response has content.
	Empty() bool
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}
