package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/controllers"
	"github.com/mrdniwe/r/internal/view"
	"github.com/mrdniwe/r/pkg/templator"
)

// global app vars
var (
	pgs *templator.Pages
	r   *mux.Router
)

func init() {
	// Template and router init
	pgs = view.New()
	r = mux.NewRouter()

}

func main() {

	// content pages
	p := r.PathPrefix("/").Subrouter()
	controllers.Site(p, pgs)

	// API
	a := r.PathPrefix("/api/v1").Subrouter()
	controllers.Api(a, pgs)

	// Static
	static := http.FileServer(http.Dir("static"))

	// Middlewares
	// r.Use(mwr["restUri"])

	// Handle and serve
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", static))

	fmt.Println("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
