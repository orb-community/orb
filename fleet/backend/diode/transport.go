package diode

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

func MakeDiodeHandler(tracer opentracing.Tracer, dio diodeBackend, opts []kithttp.ServerOption, r *bone.Mux) {

	r.Get("/agents/backends/diode/configs", kithttp.NewServer(
		kitot.TraceServer(tracer, "view_agent_backend_configs")(viewAgentBackendConfigsEndpoint(dio)),
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
