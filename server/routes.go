package server

import "github.com/gorilla/mux"

func (s *Server) RegisterRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/user", s.Register).Methods("POST")

	s.httpServer.Handler = router
}
