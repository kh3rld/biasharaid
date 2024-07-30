package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kh3rld/biasharaid/blockchain"
	"github.com/kh3rld/biasharaid/internals/routes"
)

func main() {
	blockchain.BlockchainInstance = blockchain.InitializeBlockchain()
	blockchain.LoadData("../data.json")

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	// wrappedMux := routes.RouteChecker(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("server running @http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
