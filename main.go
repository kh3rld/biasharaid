package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// homeHandler is a handler function for the root path
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")

	log.Println("Listening on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
