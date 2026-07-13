package transport

import (
	"context"

	"minitor/config"
	"minitor/transport/http"
)

type Server struct {
	Http *http.Http
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		Http: http.NewHttp(cfg),
	}
}

func (s *Server) Run(ctx context.Context) error {
	return s.Http.Run(ctx)
}
