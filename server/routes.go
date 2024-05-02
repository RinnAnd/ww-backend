package server

import (
	"net/http"

	"github.com/RinnAnd/ww-backend/services"
)

type Server struct {
	httpServer *http.Server
	gateWay    *Gateway
}

type Gateway struct {
	UserService services.UserService
}

// func (s *Server) Serve(addr string) {
// 	server := &Server{}
// }
