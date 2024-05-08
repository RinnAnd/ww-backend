package server

import (
	"fmt"
	"net/http"

	"github.com/RinnAnd/ww-backend/services"
)

type Server struct {
	httpServer *http.Server
	gateWay    *Gateway
}

type Gateway struct {
	UserService *services.UserService
}

func NewGateway(us *services.UserService) *Gateway {
	return &Gateway{
		UserService: us,
	}
}

func NewServer(addr string) *Server {
	userService := *services.MakeUserService()
	return &Server{
		httpServer: &http.Server{
			Addr: addr,
		},
		gateWay: NewGateway(&userService),
	}
}

func (s *Server) Start() error {
	fmt.Println("[APP] is now up and listening on port", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	s.gateWay.UserService.CreateUser(w, r)
}
