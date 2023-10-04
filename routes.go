package main

func(s *server) routes(){
	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.router.HandleFunc("/movie", s.handleMovieList()).Methods("GET")
	s.router.HandleFunc("/movie", s.handleMovieCreate()).Methods("POST")
	s.router.HandleFunc("/movie/{id:[0-9]+}", s.handleMovieById()).Methods("GET")
}

