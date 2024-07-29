package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kh3rld/biasharaid/routes"
)

func main() {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	wrappedMux := routes.RouteChecker(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: wrappedMux,
	}

	fmt.Println("server running @http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
