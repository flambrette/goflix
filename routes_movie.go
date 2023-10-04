package main

import (
	"log"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
)

type jsonMovie struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Duration    int    `json:"duration"`
	TrailerUrl  string `json:"trailer_url"`
}

func (s *server) handleMovieList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movies, err := s.store.GetMovies()
		if err != nil {
			log.Printf("Cannot load movies. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}
		
		var response = make([]jsonMovie, len(movies))
		for i,m := range movies {
			response[i] = mapMovieToJson(m)
		}

		s.respond(w, r, response, http.StatusOK)
		return
	}
}

func (s *server) handleMovieById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Printf("Id is not a number (id:%v). err=%v\n", id, err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		movie, err := s.store.GetMovie(id)
		if err != nil {
			log.Printf("Cannot load movie for id:%v. err=%v\n", id, err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		if movie == nil {
			log.Printf("Cannot find movie for id:%v.\n", id)
			s.respond(w, r, nil, http.StatusNotFound)
			return
		}
		
		response := mapMovieToJson(movie)
		s.respond(w, r, response, http.StatusOK)
	}
}

func (s *server) handleMovieCreate() http.HandlerFunc {

	type movieRequest struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
		TrailerUrl  string `json:"trailer_url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//decode request
		req := movieRequest{}
		err := s.decode(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse movie request. err=%v\n", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		//create movie
		m := &Movie{
			Id: 0,
			Title: req.Title,
			ReleaseDate : req.ReleaseDate,
			Duration: req.Duration,
			TrailerUrl: req.TrailerUrl,
		}

		//store movie in db
		err = s.store.CreateMovie(m)
		if err != nil {
			log.Printf("Cannot create movie . err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = mapMovieToJson(m)
		s.respond(w, r, resp, http.StatusOK)

	}
}

func mapMovieToJson(m *Movie) jsonMovie {
	return jsonMovie{
		Id: m.Id,
		Title: m.Title,
		ReleaseDate : m.ReleaseDate,
		Duration: m.Duration,
		TrailerUrl: m.TrailerUrl,
	}
}