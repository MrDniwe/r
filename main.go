package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type templates map[string]*template.Template

var (
	tmpl templates
	r    *mux.Router
)

func init() {
	// Template and router init
	tmpl = make(templates)
	tmpl["dummy"] = template.Must(tmplFactory().ParseFiles("templates/dummy.html"))
	tmpl["home"] = template.Must(tmplFactory().ParseFiles("templates/home.html"))
	tmpl["post"] = template.Must(tmplFactory().ParseFiles("templates/post.html"))
	tmpl["info"] = template.Must(tmplFactory().ParseFiles("templates/static.html"))
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

func tmplFactory() *template.Template {
	return template.Must(template.ParseFiles(
		"templates/header.html",
		"templates/footer.html",
		"templates/navigation.html",
		"templates/page-header.html",
		"templates/page-footer.html",
		"templates/layout.html"))
}
