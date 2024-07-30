package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kh3rld/biasharaid/blockchain"
	"github.com/kh3rld/biasharaid/internals/renders"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "home.page.html", nil)
}

func Verification(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "verify.page.html", nil)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "contact.page.html", nil)
}

func Details(w http.ResponseWriter, r *http.Request) {

}

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	resp := blockchain.BlockchainInstance.Blocks
	renders.RenderTemplate(w, "dummy.page.html", resp)
}
func TestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renders.RenderTemplate(w, "test.page.html", nil)
		return
	}
	nationalID := r.URL.Query().Get("national_id")
	fmt.Println("National ID:", r)
	if nationalID == "" {
		BadRequestHandler(w, r)
		return
	}

	var block *blockchain.Block
	for _, b := range blockchain.BlockchainInstance.Blocks {
		if b.Data.NationalID == nationalID {
			block = b
			break
		}
	}

	if block == nil {
		renders.RenderTemplate(w, "not_found.page.html", nil)
		return
	}
	renders.RenderTemplate(w, "test.page.html", block)
}

func Add(w http.ResponseWriter, r *http.Request) {
	var entrepreneur blockchain.Entrepreneur
	if err := json.NewDecoder(r.Body).Decode(&entrepreneur); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if entrepreneur.FirstName == "" || entrepreneur.SecondName == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	blockchain.BlockchainInstance.AddBlock(entrepreneur)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Block added successfully"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	renders.RenderTemplate(w, "404.page.html", nil)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	renders.RenderTemplate(w, "400.page.html", nil)
}

func ServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	renders.RenderTemplate(w, "500.page.html", nil)
}
