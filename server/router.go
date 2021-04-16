package server

import (
	"github.com/ns1labs/orb/server/api/v1"
)

func (s *Server) RegisterRoutes() {
	v1RouteGroup := s.engine.Group("/v1")
	v1.RegisterExampleHandler(v1RouteGroup.Group("/example"))
}
