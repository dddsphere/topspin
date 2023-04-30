package rest

import (
	"net/http"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo/config"
)

type (
	Server struct {
		*topspin.Server
		*topspin.CQRSManager
		Config *config.Config
		Router http.Handler
		log    topspin.Logger
	}
)

func NewServer(name string, cfg *config.Config, log topspin.Logger) (server *Server) {
	return &Server{
		Server: topspin.NewServer(name, log),
		Config: cfg,
	}
}
