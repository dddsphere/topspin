package rest

import (
	"net/http"

	"github.com/dddsphere/topspin"
)

func (server *Server) InitRESTRouter(h http.Handler) {
	r := topspin.NewRouter("rest-router", server.Log())
	r.Mount("/api/v1", h)

	server.SetRouter(r)
}
