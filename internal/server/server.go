package server

import (
	"fmt"
	"net/http"

	"github.com/jay-bhogayata/blogapi/internal/config"
	"github.com/jay-bhogayata/blogapi/internal/logger"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, router http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
			Handler: router,
		},
	}
}

func (s *Server) Start() error {

	logger.Log.Info("Starting the server...")

	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
