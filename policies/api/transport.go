/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/ns1labs/orb"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"strings"
)

func MakeHandler(svcName string, svc policies.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r := bone.New()

	r.Post("/policies/agent", kithttp.NewServer(
		addPolicyEndpoint(svc),
		decodeAddPolicyRequest,
		types.EncodeResponse,
		opts...))

	r.Post("/policies/dataset", kithttp.NewServer(
		addDatasetEndpoint(svc),
		decodeAddDatasetRequest,
		types.EncodeResponse,
		opts...))

	r.GetFunc("/version", orb.Version(svcName))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeAddPolicyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, errors.ErrUnsupportedContentType
	}

	req := addPolicyReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeAddDatasetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, errors.ErrUnsupportedContentType
	}

	req := addDatasetReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch errorVal := err.(type) {
	case errors.Error:
		w.Header().Set("Content-Type", types.ContentType)
		switch {
		case errors.Contains(errorVal, errors.ErrUnauthorizedAccess):
			w.WriteHeader(http.StatusUnauthorized)

		case errors.Contains(errorVal, errors.ErrInvalidQueryParams):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, errors.ErrUnsupportedContentType):
			w.WriteHeader(http.StatusUnsupportedMediaType)

		case errors.Contains(errorVal, errors.ErrMalformedEntity):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Contains(errorVal, errors.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Contains(errorVal, errors.ErrConflict):
			w.WriteHeader(http.StatusConflict)

		case errors.Contains(errorVal, db.ErrScanMetadata):
			w.WriteHeader(http.StatusUnprocessableEntity)

		case errors.Contains(errorVal, io.ErrUnexpectedEOF),
			errors.Contains(errorVal, io.EOF):
			w.WriteHeader(http.StatusBadRequest)

		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		if errorVal.Msg() != "" {
			if err := json.NewEncoder(w).Encode(types.ErrorRes{Err: errorVal.Msg()}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
