package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/pkg/templator"
)

// global app vars
var (
	pgs *templator.Pages
	r   *mux.Router
)

func init() {
	// Template and router init
	pgs = &templator.Pages{}
	pgs.New()
	// pgs.add("templates/dummy.html", "dummy.html")
	// pgs.add("templates/home.html", "home.html")
	// pgs.add("templates/post.html", "post.html")
	// pgs.add("templates/static.html", "static.html")
	r = mux.NewRouter()

}

func main() {

	// content pages
	r.HandleFunc("/", h["home"]).Methods("GET")
	r.HandleFunc("/post/{id}", h["post"]).Methods("GET")
	r.HandleFunc("/info/{page}", h["info"]).Methods("GET")

	// API
	a := r.PathPrefix("/api/v1").Subrouter()
	a.HandleFunc("/test", h["dummy"]).Methods("GET", "POST")

	// Static
	static := http.FileServer(http.Dir("static"))

	// Middlewares
	r.Use(mwr["restUri"])

	// Handle and serve
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", static))

	fmt.Println("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
