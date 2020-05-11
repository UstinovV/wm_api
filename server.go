package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
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
	s.router.HandleFunc("/offers", s.getOffersList)
	s.router.HandleFunc("/offer/{id:[a-zA-Z0-9]+}", s.getOffer)
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

func (s *Server) getOffer(w http.ResponseWriter, r *http.Request)  {
	//vars := mux.Vars(r)
}

func (s *Server) getOffersList(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	offersList := make([]Offer,0, 5)
	offer := Offer{}
	encoder := json.NewEncoder(w)

	query := "SELECT short_id, title, content, created_at from offer "
	where, limit, offset := "", "", ""

	if val := r.URL.Query().Get("query"); val != "" {
		where = where + fmt.Sprintf("title like %s", "'%" + val + "%'")
	}
	if val := r.URL.Query().Get("limit"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			log.Println("Bad request parameter `limit`", err)
		}
		limit = limit + fmt.Sprintf(" limit %d", i)
	}
	if val := r.URL.Query().Get("offset"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			log.Println("Bad request parameter `offset`", err)
		}
		offset = offset + fmt.Sprintf(" offset %d", i)
	}

	if len(where) > 0 {
		query = query + " where " + where
	}
	if len(limit) > 0 {
		query = query + limit
	} else {
		query = query + " limit 10"
	}
	if len(offset) > 0 {
		query = query + offset
	}

	rows, err := s.db.db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": err.Error()})
		return
		//log.Fatal("Error", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&offer.Id, &offer.Title, &offer.Content, &offer.CreatedAt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{"error": err.Error()})
			return
			//log.Fatal("Error scan")
		}
		offersList = append(offersList, offer)
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(offersList)

}