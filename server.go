package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	store Store

}

func newServer() *server {
	s := &server{
		router: mux.NewRouter(),
	}

	s.routes()
	return s
}

func (s *server) serveHttp(w http.ResponseWriter, r *http.Request) {
	//log the request and then execute the request provided in param
	logRequest(s.router.ServeHTTP).ServeHTTP(w,r)
}

func (s *server) respond(w http.ResponseWriter, _ *http.Request, data interface{}, status int) {
	w.Header().Add("Content-Type", "application.json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Cannot format json. err=%v\n", err)
	}
}

func (s *server) decode(w http.ResponseWriter, r *http.Request, data interface{}) (error) {
	return json.NewDecoder(r.Body).Decode(data)
}