package handlers

import (
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
	switch r.Method {
	case "GET":
		renders.RenderTemplate(w, "test.page.html", nil)
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}

		nationalID := r.FormValue("national_id")

		if nationalID == "" {
			BadRequestHandler(w, r)
			return
		}

		var block *blockchain.Block
		fmt.Println(block)
		for _, b := range blockchain.BlockchainInstance.Blocks {
			fmt.Println(b.Data.NationalID)
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
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "add.page.html", nil)
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
