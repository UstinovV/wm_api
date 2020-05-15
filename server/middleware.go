package server

import (
	"go.uber.org/zap"
	"net/http"
)

func prepareResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(
			"Request accepted",
			zap.String("url", r.URL.String()),
			zap.String("user-agent", r.UserAgent()),
		)
		next.ServeHTTP(w, r)
	})
}