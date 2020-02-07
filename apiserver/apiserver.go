package apiserver

import (
	"github.com/UstinovV/wm_api/database"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type APIServer struct{
	config *Config
	router *mux.Router
	database *database.Database
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	s.configureRouter()

	if err := s.configureDatabase(); err != nil {
		return err
	}

	return http.ListenAndServe(s.config.Port, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/index", s.handleIndex())
}

func (s *APIServer) configureDatabase() error {
	db := database.New(s.config.Database)
	if err := db.Open(); err != nil {
		return err
	}

	s.database = db

	return nil
}


func (s *APIServer) handleIndex()  http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "index")
	}

}