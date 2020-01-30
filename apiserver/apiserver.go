package apiserver

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type APIServer struct{
	config *Config
	router *mux.Router
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	s.configureRouter()
	return http.ListenAndServe(s.config.Port, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/index", s.handleIndex())
}

func (s *APIServer) handleIndex()  http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "index")
	}

}