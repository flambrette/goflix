package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStore struct {
	movieId int64
	movies  []*Movie
}

func (t *testStore) Open() error {
	return nil
}

func (t *testStore) Close() error {
	return nil
}

func (t *testStore) GetMovies() ([]*Movie, error) {
	return t.movies, nil
}

func (t *testStore) GetMovie(id int64) (*Movie, error) {
	for _, m := range t.movies {
		if m.Id == id {
			return m, nil
		}
	}
	return nil, nil
}

func (t *testStore) CreateMovie(m *Movie) error {
	t.movieId++
	m.Id = t.movieId
	t.movies = append(t.movies, m)
	return nil
}

func (t *testStore) UpdateMovie(newM *Movie) error {
	var result []*Movie

	for _, m := range t.movies {
		if m.Id != newM.Id {
			result = append(result, m)
		}
	}
	result = append(result, newM)
	t.movies = result
	return nil

}

func (t *testStore) DeleteMovie(id int64) error {
	var result []*Movie

	for _, m := range t.movies {
		if m.Id != id {
			result = append(result, m)
		}
	}
	t.movies = result
	return nil
}

func TestMovieCreateUnit(t *testing.T) {
	//create server
	srv := newServer()
	srv.store = &testStore{}

	//prepare json body
	p := struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"releaseDate"`
		Duration    int    `json:"duration"`
		TrailerUrl  string `json:"trailerUrl"`
	}{
		Title:       "inception",
		ReleaseDate: "2010-07-18",
		Duration:    148,
		TrailerUrl:  "http://url",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/movie", &buf)
	w := httptest.NewRecorder()

	//call create with fake request
	f := srv.handleMovieCreate()
	f(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMovieCreateIntegration(t *testing.T) {
	//create server
	srv := newServer()
	srv.store = &testStore{}

	//prepare json body
	p := struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"releaseDate"`
		Duration    int    `json:"duration"`
		TrailerUrl  string `json:"trailerUrl"`
	}{
		Title:       "inception",
		ReleaseDate: "2010-07-18",
		Duration:    148,
		TrailerUrl:  "http://url",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/movie", &buf)
	w := httptest.NewRecorder()

	srv.serveHttp(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
