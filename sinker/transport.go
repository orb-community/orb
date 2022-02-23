package sinker

import (
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func MakeHandler(svcName string) http.Handler {
	r := bone.New()
	r.Handle("/metrics", promhttp.Handler())
	return r
}