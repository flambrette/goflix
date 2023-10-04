package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("GoFlix")

	srv := newServer()
	srv.store = &dbStore{}

	err := srv.store.Open()
	HandleError(err)

	defer srv.store.Close()

	http.HandleFunc("/", srv.serveHttp)
	log.Printf("Serving http...")
	err = http.ListenAndServe(":9000", nil)
	HandleError(err)
}

func HandleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

