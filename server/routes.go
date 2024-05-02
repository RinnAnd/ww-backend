package server

import "github.com/gorilla/mux"

func (s *Server) RegisterRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/user", s.Register).Methods("GET")

	s.httpServer.Handler = router
}
