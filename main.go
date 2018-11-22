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

func main() {

	// Template and router init
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	r = mux.NewRouter()

	// content pages
	r.HandleFunc("/", emptyHandler).Methods("GET")
	r.HandleFunc("/post/{id}", emptyHandler).Methods("GET")
	r.HandleFunc("/about", emptyHandler).Methods("GET")

	// API
	a := r.PathPrefix("/api/v1").Subrouter()
	a.HandleFunc("/test", emptyHandler).Methods("GET", "POST")

	// Middlewares
	r.Use(mwr["restUri"])

	// Handle and serve
	http.Handle("/", r)
	fmt.Println("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
