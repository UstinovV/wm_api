package server

import (
	"github.com/UstinovV/wm_api/database"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	config *Config
	router *mux.Router
	db     *database.DB
	logger *zap.SugaredLogger
}

func NewServer(config *Config, logger *zap.SugaredLogger) *Server {
	return &Server{
		config: config,
		router: mux.NewRouter(),
		logger: logger,
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
	s.router.Use(prepareResponseMiddleware)
	s.router.Use(s.logRequestMiddleware)
	s.router.HandleFunc("/offers", s.getOffersListHandler)
	s.router.HandleFunc("/offer/{id:[a-zA-Z0-9]+}", s.getOfferHandler)
	s.router.HandleFunc("/companies", s.getCompaniesListHandler)
	s.router.HandleFunc("/company/{id:[a-zA-Z0-9']+}", s.getCompanyHandler)
}

func (s *Server) configureDatabase() error {
	db := database.NewDB(s.config.DBConnection)
	if err := db.Open(); err != nil {
		return err
	}

	s.db = db

	return nil
}

