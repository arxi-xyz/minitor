package transport

import (
	"context"

	"minitor/transport/http"
)

type Server struct {
	Http *http.Http
}

func NewServer() *Server {
	return &Server{
		Http: http.NewHttp(),
	}
}

func (s *Server) Run(ctx context.Context) error {
	return s.Http.Run(ctx)
}
