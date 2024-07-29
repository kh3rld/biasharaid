package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", nil).Methods("GET")

	log.Println("Listening on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
