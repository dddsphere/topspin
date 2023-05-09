package rest

import (
	"net/http"

	"github.com/dddsphere/topspin"
)

type (
	Server struct {
		*topspin.Server
		*topspin.CQRSManager
		Config *topspin.Config
		Router http.Handler
		log    topspin.Logger
	}
)

func NewServer(name string, cfg *topspin.Config, log topspin.Logger) (server *Server) {
	return &Server{
		Server: topspin.NewServer(name, log),
		Config: cfg,
	}
}
