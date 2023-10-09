package main

func(s *server) routes(){
	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.router.HandleFunc("/token", s.handleTokenCreate()).Methods("POST")
	s.router.HandleFunc("/movie", s.loggedOnly(s.handleMovieList())).Methods("GET")
	s.router.HandleFunc("/movie", s.loggedOnly(s.handleMovieCreate())).Methods("POST")
	s.router.HandleFunc("/movie/{id:[0-9]+}", s.loggedOnly(s.handleMovieUpdate())).Methods("PUT")
	s.router.HandleFunc("/movie/{id:[0-9]+}", s.loggedOnly(s.handleMovieDeletion())).Methods("DELETE")
	s.router.HandleFunc("/movie/{id:[0-9]+}", s.loggedOnly(s.handleMovieById())).Methods("GET")
}

