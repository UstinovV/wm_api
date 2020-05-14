package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Server struct {
	config *Config
	router *mux.Router
	db     *DB
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
	s.router.HandleFunc("/offers", s.getOffersList)
	s.router.HandleFunc("/offer/{id:[a-zA-Z0-9]+}", s.getOffer)
	s.router.HandleFunc("/companies", s.getCompaniesList)
	s.router.HandleFunc("/company/{id:[a-zA-Z0-9]+}", s.getCompany)
}

func (s *Server) configureDatabase() error {
	db := NewDB(s.config.DBConnection)
	if err := db.Open(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Server) getOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	offer := Offer{}
	encoder := json.NewEncoder(w)

	query := "SELECT short_id, title, content, created_at from offer where short_id = $1"
	if val, ok := vars["id"]; ok {
		err := s.db.db.QueryRow(query, val).Scan(&offer.Id, &offer.Title, &offer.Content, &offer.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusOK)
				encoder.Encode(map[string]string{})
			} else {
				w.WriteHeader(http.StatusBadRequest)
				encoder.Encode(map[string]string{"error": err.Error()})
			}
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": "Missing parameter ID"})
		return
	}
	w.WriteHeader(http.StatusOK)
	encoder.Encode(offer)
}

func (s *Server) getOffersList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	offersList := make([]Offer, 0, 10)
	offer := Offer{}
	encoder := json.NewEncoder(w)

	query := "SELECT short_id, title, content, created_at from offer "
	allowedParams := []string{"title", "limit", "offset"}
	parsedString, queryArgs := parseQueryParams(r.URL.Query(), allowedParams)

	rows, err := s.db.db.Query(query + parsedString, queryArgs...)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": err.Error()})
		return
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

func (s *Server) getCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	company := Company{}
	encoder := json.NewEncoder(w)

	query := "SELECT c.short_id, ci.name, ci.description from company c where short_id = $1 left join company_info ci on c.info_id = ci.id"
	if val, ok := vars["id"]; ok {
		err := s.db.db.QueryRow(query, val).Scan(&company.Id, &company.Name, &company.Description)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusOK)
				encoder.Encode(map[string]string{})
			} else {
				w.WriteHeader(http.StatusBadRequest)
				encoder.Encode(map[string]string{"error": err.Error()})
			}
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": "Missing parameter ID"})
		return
	}
	w.WriteHeader(http.StatusOK)
	encoder.Encode(company)
}

func (s *Server) getCompaniesList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	companiesList := make([]Company, 0, 10)
	company := Company{}
	encoder := json.NewEncoder(w)

	query := "SELECT c.short_id, ci.name, ci.description from company c left join company_info ci on c.info_id = ci.id"

	allowedParams := []string{"name", "limit", "offset"}
	parsedString, queryArgs := parseQueryParams(r.URL.Query(), allowedParams)

	rows, err := s.db.db.Query(query + parsedString, queryArgs...)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&company.Id, &company.Name, &company.Description)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{"error": err.Error()})
			return
			//log.Fatal("Error scan")
		}
		companiesList = append(companiesList, company)
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(companiesList)

}


func parseQueryParams(values url.Values, requestedParams []string) (string, []interface{}) {
	result, where, limit, offset := "", "", "", ""
	args := []interface{}{}
	for _, param := range requestedParams {
		if val := values.Get(param); val != "" {
			switch param {
			case "title":
				where = where + fmt.Sprintf("title like $%d", len(args)+1)
				args = append(args, "%"+val+"%")
			case "name":
				where = where + fmt.Sprintf("ci.name like $%d", len(args)+1)
				args = append(args, "%"+val+"%")
			case "limit":
				i, err := strconv.Atoi(val)
				if err != nil {
					log.Println("Bad request parameter `limit`", err)
				} else {
					limit = limit + fmt.Sprintf(" limit $%d", len(args)+1)
					args = append(args, i)
				}
			case "offset":
				i, err := strconv.Atoi(val)
				if err != nil {
					log.Println("Bad request parameter `offset`", err)
				}
				offset = offset + fmt.Sprintf(" offset $%d", i)
				args = append(args, i)
			}
		}
	}

	if len(where) > 0 {
		result = result + " where " + where
	}
	if len(limit) > 0 {
		result = result + limit
	} else {
		result = result + " limit 10"
	}
	if len(offset) > 0 {
		result = result + offset
	}

	return result, args
}