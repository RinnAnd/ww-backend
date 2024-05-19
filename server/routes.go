package server

import "github.com/gorilla/mux"

func (s *Server) RegisterRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/ping", s.Pong).Methods("GET")

	router.HandleFunc("/user", s.Register).Methods("POST")
	router.HandleFunc("/user/auth", s.Login).Methods("POST")

	s.httpServer.Handler = router
}
