package pktvisor

import (
	"context"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

func MakePktvisorHandler(tracer opentracing.Tracer, pkt pktvisorBackend, opts []kithttp.ServerOption, r *bone.Mux) {

	r.Get("/backends", kithttp.NewServer(
		kitot.TraceServer(tracer, "list_backends")(listAgentBackendsEndpoint(pkt)),
		decodeListBackends,
		types.EncodeResponse,
		opts...))
	r.Get("/backends/:id/handler", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_handler")(viewAgentBackendHandlerEndpoint(pkt)),
		decodeBackendView,
		types.EncodeResponse,
		opts...))
	r.Get("/backends/:id/input", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_input")(viewAgentBackendInputEndpoint(pkt)),
		decodeBackendView,
		types.EncodeResponse,
		opts...))
	r.Get("/backends/:id/taps", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_taps")(viewAgentBackendTapsEndpoint(pkt)),
		decodeBackendView,
		types.EncodeResponse,
		opts...))
}

func decodeBackendView(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewResourceReq{
		token: r.Header.Get("Authorization"),
		id:    bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeListBackends(_ context.Context, r *http.Request) (interface{}, error) {
	req := listAgentBackendsReq{token: r.Header.Get("Authorization")}
	return req, nil
}
