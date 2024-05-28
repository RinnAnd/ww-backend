package server

import "github.com/gorilla/mux"

func (s *Server) RegisterRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/ping", s.Pong).Methods("GET")

	router.HandleFunc("/user", s.AllUsers).Methods("GET")
	router.HandleFunc("/user", s.Register).Methods("POST")
	router.HandleFunc("/user/auth", s.Login).Methods("POST")
	router.HandleFunc("/user/password", s.UpdatePassword).Methods("POST")

	router.HandleFunc("/finance", s.NewFinance).Methods("POST")
	router.HandleFunc("/finance/{id}", s.UserFinances).Methods("GET")
	router.HandleFunc("/expense", s.NewExpense).Methods("POST")
	router.HandleFunc("/saving", s.NewSaving).Methods("POST")

	s.httpServer.Handler = router
}
