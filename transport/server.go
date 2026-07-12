package transport

import (
	"minitor/transport/http"
)

type Server struct {
	Http *http.Http
}

func NewServer() *Server {
	return &Server{
		Http: http.NewHttp()}
}

func (s *Server) Run() {
	s.Http.Run()
}
