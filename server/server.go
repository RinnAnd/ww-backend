package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/RinnAnd/ww-backend/services"
)

type Server struct {
	httpServer *http.Server
	gateWay    *Gateway
}

func (s *Server) Pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pong")
}

type Gateway struct {
	UserService    *services.UserService
	FinanceService *services.FinanceService
}

func NewGateway(us *services.UserService, fs *services.FinanceService) *Gateway {
	return &Gateway{
		UserService:    us,
		FinanceService: fs,
	}
}

func NewServer(addr string, pool *sql.DB) *Server {
	userService := *services.MakeUserService(pool)
	financeService := *services.MakeFinanceService(pool)
	return &Server{
		httpServer: &http.Server{
			Addr: addr,
		},
		gateWay: NewGateway(&userService, &financeService),
	}
}

func (s *Server) Start() error {
	fmt.Println("[APP] is now up and listening on port", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	s.gateWay.UserService.CreateUser(w, r)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	s.gateWay.UserService.Auth(w, r)
}

func (s *Server) AllUsers(w http.ResponseWriter, r *http.Request) {
	s.gateWay.UserService.GetUsers(w, r)
}

func (s *Server) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	s.gateWay.UserService.ChangePassword(w, r)
}

func (s *Server) NewFinance(w http.ResponseWriter, r *http.Request) {
	s.gateWay.FinanceService.CreateFinance(w, r)
}

func (s *Server) UserFinances(w http.ResponseWriter, r *http.Request) {
	s.gateWay.FinanceService.GetUserFinances(w, r)
}

func (s *Server) NewExpense(w http.ResponseWriter, r *http.Request) {
	s.gateWay.FinanceService.CreateExpense(w, r)
}
