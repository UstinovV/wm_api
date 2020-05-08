package main

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Server struct{
	config *Config
	router *mux.Router
	db *DB
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.configureRouter()

	if err := s.configureDatabase(); err != nil {
		return err
	}

	return http.ListenAndServe(s.config.Port, s.router)
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/index", s.handleIndex())
}

func (s *Server) configureDatabase() error {
	db := NewDB(s.config.DBConnection)
	if err := db.Open(); err != nil {
		return err
	}

	s.db = db

	return nil
}


func (s *Server) handleIndex()  http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "index")
	}

}