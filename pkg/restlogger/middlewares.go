package restlogger

import (
	"log"
	"net/http"
)

// Map of middleware functions
// for simple usage by name
var mwr = map[string]func(next http.Handler) http.Handler{
	// just logs URI from request
	"restUri": func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RequestURI)
			next.ServeHTTP(w, r)
		})
	},
}
