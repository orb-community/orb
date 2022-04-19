/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package http

import (
	"context"
	"encoding/json"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/ns1labs/orb/buildinfo"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/internal/httputil"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"strings"
)

const (
	contentType = "application/json"
	offsetKey   = "offset"
	limitKey    = "limit"
	nameKey     = "name"
	orderKey    = "order"
	dirKey      = "dir"
	metadataKey = "metadata"
	tagsKey     = "tags"
	defOffset   = 0
	defLimit    = 10
)

func MakeHandler(tracer opentracing.Tracer, svcName string, svc policies.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r := bone.New()

	r.Post("/policies/agent", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_policy")(addPolicyEndpoint(svc)),
		decodeAddPolicyRequest,
		types.EncodeResponse,
		opts...))
	r.Get("/policies/agent/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_policy")(viewPolicyEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))
	r.Get("/policies/agent", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_policies")(listPoliciesEndpoint(svc)),
		decodeList,
		types.EncodeResponse,
		opts...))
	r.Put("/policies/agent/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "edit_policy")(editPoliciyEndpoint(svc)),
		decodePolicyUpdate,
		types.EncodeResponse,
		opts...))
	r.Post("/policies/agent/:id/duplicate", kithttp.NewServer(
		kitot.TraceServer(tracer, "duplicate_policy")(duplicatePolicyEndpoint(svc)),
		decodePolicyDuplicate,
		types.EncodeResponse,
		opts...))
	r.Delete("/policies/agent/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "remove_policy")(removePolicyEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))
	r.Post("/policies/agent/validate", kithttp.NewServer(
		kitot.TraceServer(tracer, "validate_policy")(validatePolicyEndpoint(svc)),
		decodeAddPolicyRequest,
		types.EncodeResponse,
		opts...))

	r.Post("/policies/dataset", kithttp.NewServer(
		kitot.TraceServer(tracer, "add_dataset")(addDatasetEndpoint(svc)),
		decodeAddDatasetRequest,
		types.EncodeResponse,
		opts...))
	r.Put("/policies/dataset/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "edit_dataset")(editDatasetEndpoint(svc)),
		decodeDatasetUpdate,
		types.EncodeResponse,
		opts...))
	r.Delete("/policies/dataset/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "remove_dataset")(removeDatasetEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))
	r.Get("/policies/dataset/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_dataset")(viewDatasetEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...))
	r.Get("/policies/dataset", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_datasets")(listDatasetEndpoint(svc)),
		decodeList,
		types.EncodeResponse,
		opts...))

	r.Post("/policies/dataset/validate", kithttp.NewServer(
		kitot.TraceServer(tracer, "validate_dataset")(validateDatasetEndpoint(svc)),
		decodeAddDatasetRequest,
		types.EncodeResponse,
		opts...))

	r.GetFunc("/version", buildinfo.Version(svcName))
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

func decodeView(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewResourceReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodePolicyUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errors.ErrUnsupportedContentType
	}

	req := updatePolicyReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeDatasetUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errors.ErrUnsupportedContentType
	}

	req := updateDatasetReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(fleet.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeList(_ context.Context, r *http.Request) (interface{}, error) {
	o, err := httputil.ReadUintQuery(r, offsetKey, defOffset)
	if err != nil {
		return nil, err
	}

	l, err := httputil.ReadUintQuery(r, limitKey, defLimit)
	if err != nil {
		return nil, err
	}

	n, err := httputil.ReadStringQuery(r, nameKey, "")
	if err != nil {
		return nil, err
	}

	or, err := httputil.ReadStringQuery(r, orderKey, "")
	if err != nil {
		return nil, err
	}

	d, err := httputil.ReadStringQuery(r, dirKey, "")
	if err != nil {
		return nil, err
	}

	m, err := httputil.ReadMetadataQuery(r, metadataKey, nil)
	if err != nil {
		return nil, err
	}

	t, err := httputil.ReadTagQuery(r, tagsKey, nil)
	if err != nil {
		return nil, err
	}

	req := listResourcesReq{
		token: r.Header.Get("Authorization"),
		pageMetadata: policies.PageMetadata{
			Offset:   o,
			Limit:    l,
			Name:     n,
			Order:    or,
			Dir:      d,
			Metadata: m,
			Tags:     t,
		},
	}

	return req, nil
}

func decodePolicyDuplicate(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errors.ErrUnsupportedContentType
	}

	req := duplicatePolicyReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(fleet.ErrMalformedEntity, err)
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
