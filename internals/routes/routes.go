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
}

func RegisterRoutes(mux *http.ServeMux) {
	staticDir := renders.GetProjectRoot("views", "static")
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r)
	})

	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		handlers.Verification(w, r)
	})

	mux.HandleFunc("/details", func(w http.ResponseWriter, r *http.Request) {
		handlers.Details(w, r)
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