package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// pseudo-class for working with complex templates
type page struct {
	t    *template.Template
	path string
	name string
}

// initialization of complex template
func (p *page) new(source *template.Template, path string, name string) {
	p.path = path
	p.name = name
	p.t = template.Must(source.Clone())
	p.t = template.Must(p.t.ParseFiles(path))
}

// complex template executor
func (p *page) execute(w http.ResponseWriter, r *http.Request) {
	err := p.t.ExecuteTemplate(w, p.name, r)
	if err != nil {
		log.Println(err)
	}
}

// storage of complex templates
type pages struct {
	items map[string]*page
	root  *template.Template
}

func (t *pages) new() {
	t.root = template.Must(template.ParseFiles(
		"templates/header.html",
		"templates/footer.html",
		"templates/navigation.html",
		"templates/page-header.html",
		"templates/page-footer.html",
		"templates/layout.html"))
	t.items = make(map[string]*page)
}

// right way to add another complex template in our set
func (t *pages) add(path string, name string) {
	t.items[name] = &page{}
	t.items[name].new(t.root, path, name)
}

// global app vars
var (
	pgs *pages
	r   *mux.Router
)

func init() {
	// Template and router init
	pgs = &pages{}
	pgs.new()
	pgs.add("templates/dummy.html", "dummy.html")
	pgs.add("templates/home.html", "home.html")
	pgs.add("templates/post.html", "post.html")
	pgs.add("templates/static.html", "static.html")
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
