package server

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

func (s *Server) getOfferHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	if val, ok := vars["id"]; ok {
		offer, err := s.db.GetOffer(val)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(map[string]string{})
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{"error": err.Error()})
		} else {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(offer)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": "Missing parameter ID"})
	}

}

func (s *Server) getOffersListHandler(w http.ResponseWriter, r *http.Request) {

	allowedParams := []string{"title", "limit", "offset"}
	parsedString, queryArgs := parseQueryParams(r.URL.Query(), allowedParams)
	encoder := json.NewEncoder(w)
	result, err := s.db.GetOffersList(parsedString, queryArgs)
	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		//log error
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(result)

}

func (s *Server) getCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(w)

	if val, ok := vars["id"]; ok {
		company, err := s.db.GetCompany(val)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(map[string]string{})
			return
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		encoder.Encode(company)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": "Missing parameter ID"})
		return
	}

}

func (s *Server) getCompaniesListHandler(w http.ResponseWriter, r *http.Request) {

	encoder := json.NewEncoder(w)

	allowedParams := []string{"name", "limit", "offset"}
	parsedString, queryArgs := parseQueryParams(r.URL.Query(), allowedParams)
	companiesList, err := s.db.GetCompaniesList(parsedString, queryArgs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(map[string]string{"error": err.Error()})
		return
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