package otelinf

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-zoo/bone"
	"github.com/opentracing/opentracing-go"
	"github.com/orb-community/orb/pkg/types"

	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeOtelinfHandler(tracer opentracing.Tracer, dio otelinfBackend, opts []kithttp.ServerOption, r *bone.Mux) {

	r.Get("/agents/backends/otelinf/handlers", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_handler")(viewAgentBackendHandlerEndpoint(dio)),
		decodeBackendView,
		types.EncodeResponse,
		opts...))
	r.Get("/agents/backends/otelinf/inputs", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_input")(viewAgentBackendInputEndpoint(dio)),
		decodeBackendView,
		types.EncodeResponse,
		opts...))
	r.Get("/agents/backends/otelinf/taps", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_taps")(viewAgentBackendTapsEndpoint(dio)),
		decodeBackendView,
		types.EncodeResponse,
		opts...))
}

func decodeBackendView(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewResourceReq{
		token: parseJwt(r),
	}
	return req, nil
}

func parseJwt(r *http.Request) (token string) {
	if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		token = r.Header.Get("Authorization")[7:]
	}
	return
}
