package routes

import (
	"net/http"
	"strings"

	"github.com/kh3rld/biasharaid/internals/handlers"
	"github.com/kh3rld/biasharaid/internals/renders"
)

var allowedRoutes = map[string]bool{
	"/":        true,
	"/verify":  true,
	"/details": true,
	"/dummy":   true,
	"/contact": true,
	"/test":    true,
	"/signup":  true,
	"/addpage": true,
	"/about":   true,
	"/help":    true,
}

func RegisterRoutes(mux *http.ServeMux) {
	staticDir := renders.GetProjectRoot("views", "static")
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r)
	})

	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		handlers.VerifyHandler(w, r)
	})

	mux.HandleFunc("/details", func(w http.ResponseWriter, r *http.Request) {
		handlers.Details(w, r)
	})

	mux.HandleFunc("/dummy", func(w http.ResponseWriter, r *http.Request) {
		handlers.DummyHandler(w, r)
	})

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		handlers.DummyHandler(w, r)
	})

	mux.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		handlers.Contact(w, r)
	})

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handlers.Add(w, r)
	})

	mux.HandleFunc("/addpage", func(w http.ResponseWriter, r *http.Request) {
		handlers.Addpage(w, r)
	})

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		handlers.About(w, r)
	})

	mux.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
		handlers.Help(w, r)
	})
}

func RouteChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		if _, ok := allowedRoutes[r.URL.Path]; !ok {
			handlers.NotFoundHandler(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
