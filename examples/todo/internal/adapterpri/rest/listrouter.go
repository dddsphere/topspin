package rest

import (
	"net/http"

	"github.com/dddsphere/topspin"
)

func (server *Server) InitRESTRouter(h http.Handler) {
	rr := topspin.NewRouter("rest-router", server.Log())
	rr.Mount("/api/v1", h)

	server.SetRouter(rr)
}
