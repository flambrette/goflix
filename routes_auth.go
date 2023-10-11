package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

const JWT_KEY = "training.go"

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to goflix")
	}
}

func (s *server) handleTokenCreate() http.HandlerFunc {

	type request struct {
			Username string `json:"username"`
			Password string `json:"password"`
	}

	type response struct {
		Token string `json:"token"`
	}

	type responseError struct {
		Error string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		//Parsing request
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			msg := fmt.Sprintf("Cannot parse request, err=%v", err)
			s.respond(w, r, responseError{
				Error: msg,
			}, http.StatusBadRequest)
			return
		}

		userFound, err := s.store.AuthenticateUser(req.Username, req.Password)
		if err != nil {
			msg := fmt.Sprintf("Cannot find user, err=%v", err)
			s.respond(w, r, responseError{
				Error: msg,
			}, http.StatusInternalServerError)
			return
		}

		if ! userFound {
			s.respond(w, r, responseError{
				Error: "Invalid credentials",
			}, http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username" :req.Username,
			"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			"iat" :time.Now().Unix(),
		})
		
		//Generate token
		bytes := []byte(JWT_KEY)
		tokenStr, err := token.SignedString(bytes)
		if err != nil {
			msg := fmt.Sprintf("Cannot produce token, err=%v", err)
			s.respond(w, r, responseError{
				Error: msg,
			}, http.StatusInternalServerError)
			return
		}

		s.respond(w, r, response {
			Token: tokenStr,
		}, http.StatusOK)
	}
}
