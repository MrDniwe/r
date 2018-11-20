package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	tmpl *template.Template
)

func main() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", emptyHandler)
	r.HandleFunc("/post/{id}", emptyHandler)
	r.Use(mwr["restUri"])
	http.Handle("/", r)

	fmt.Println("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
