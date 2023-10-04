package main

import (
	"log"
	"net/http"
)

/*
* log all http requests
*/
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%v] %v", r.Method, r.RequestURI)
		next.ServeHTTP(w,r)
	}
}