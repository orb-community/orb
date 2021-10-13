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
	"github.com/ns1labs/orb"
	"github.com/ns1labs/orb/internal/httputil"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"strings"
)

const (
	offsetKey   = "offset"
	limitKey    = "limit"
	nameKey     = "name"
	orderKey    = "order"
	dirKey      = "dir"
	metadataKey = "metadata"
	defOffset   = 0
	defLimit    = 10
)

func MakeHandler(tracer opentracing.Tracer, svcName string, svc sinks.SinkService) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r := bone.New()

	r.Post("/sinks", kithttp.NewServer(
		kitot.TraceServer(tracer, "create_sink")(addEndpoint(svc)),
		decodeAddRequest,
		types.EncodeResponse,
		opts...,
	))
	r.Put("/sinks/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "edit_sink")(updateSinkEndpoint(svc)),
		decodeEditRequest,
		types.EncodeResponse,
		opts...,
	))
	r.Get("/sinks", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_sinks")(listSinksEndpoint(svc)),
		decodeList,
		types.EncodeResponse,
		opts...,
	))
	r.Post("/sinks/validate", kithttp.NewServer(
		kitot.TraceServer(tracer, "validate_sink")(validateSinkEndpoint(svc)),
		decodeValidateRequest,
		types.EncodeResponse,
		opts...,
	))
	r.Get("/features/sinks", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_backends")(listBackendsEndpoint(svc)),
		decodeListBackends,
		types.EncodeResponse,
		opts...,
	))
	r.Get("/features/sinks/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_backend")(viewBackendEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...,
	))
	r.Get("/sinks/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_sink")(viewSinkEndpoint(svc)),
		decodeView,
		types.EncodeResponse,
		opts...,
	))
	r.Delete("/sinks/:id", kithttp.NewServer(
		kitot.TraceServer(tracer, "delete_sink")(deleteSinkEndpoint(svc)),
		decodeDeleteRequest,
		types.EncodeResponse,
		opts...,
	))
	r.Get("/sinks/statistics", kithttp.NewServer(
		kitot.TraceServer(tracer, "sink_statistics")(sinksStatisticsEndpoint(svc)),
		decodeSinksStatistics,
		types.EncodeResponse,
		opts...,
	))

	r.GetFunc("/version", orb.Version(svcName))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, errors.ErrUnsupportedContentType
	}

	req := addReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeEditRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, errors.ErrUnsupportedContentType
	}
	req := updateSinkReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}

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

func decodeListBackends(_ context.Context, r *http.Request) (interface{}, error) {
	req := listBackendsReq{token: r.Header.Get("Authorization")}
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

	req := listResourcesReq{
		token: r.Header.Get("Authorization"),
		pageMetadata: sinks.PageMetadata{
			Offset:   o,
			Limit:    l,
			Name:     n,
			Order:    or,
			Dir:      d,
			Metadata: m,
		},
	}

	return req, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := deleteSinkReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}

	return req, nil
}

func decodeValidateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, errors.ErrUnsupportedContentType
	}

	req := validateReq{token: r.Header.Get("Authorization")}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return req, nil
}

func decodeSinksStatistics(_ context.Context, r *http.Request) (interface{}, error) {
	req := sinksStatisticsReq{token: r.Header.Get("Authorization")}
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
		case errors.Contains(errorVal, sinks.ErrInvalidBackend):
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
