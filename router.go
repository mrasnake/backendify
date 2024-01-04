package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// InitServer defines all the routes and attaches them to the Server object.
func (s *Server) InitServer() {
	s.Router = mux.NewRouter()

	s.Router.HandleFunc("/status", s.GetStatus).Methods(http.MethodGet)
	s.Router.HandleFunc("/company", s.GetCompany).Methods(http.MethodGet)
}