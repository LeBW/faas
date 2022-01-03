package handlers

import (
	"log"
	"net/http"
)

func MakePredictHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("In MakePredictHandler. r.Host: %s, r.RemoteAddress: %s\n", r.Host, r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}
