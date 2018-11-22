package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	tmpl *template.Template
	r    *mux.Router
)

func init() {
	// Template and router init
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	r = mux.NewRouter()
}

func main() {

	// content pages
	r.HandleFunc("/", h["dummy"]).Methods("GET")
	r.HandleFunc("/post/{id}", h["dummy"]).Methods("GET")
	r.HandleFunc("/about", h["dummy"]).Methods("GET")

	// API
	a := r.PathPrefix("/api/v1").Subrouter()
	a.HandleFunc("/test", h["dummy"]).Methods("GET", "POST")

	// Middlewares
	r.Use(mwr["restUri"])

	// Handle and serve
	http.Handle("/", r)
	fmt.Println("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
