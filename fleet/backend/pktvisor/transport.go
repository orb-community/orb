package pktvisor

import (
	"context"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"strings"
)

func MakePktvisorHandler(tracer opentracing.Tracer, pkt pktvisorBackend, opts []kithttp.ServerOption, r *bone.Mux) {

	r.Get("/agents/backends/pktvisor/handlers", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_handler")(viewAgentBackendHandlerEndpoint(pkt)),
		decodeBackendView,
		types.EncodeResponse,
		opts...))
	r.Get("/agents/backends/pktvisor/inputs", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_input")(viewAgentBackendInputEndpoint(pkt)),
		decodeBackendView,
		types.EncodeResponse,
		opts...))
	r.Get("/agents/backends/pktvisor/taps", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_taps")(viewAgentBackendTapsEndpoint(pkt)),
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
