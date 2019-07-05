package templator

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Page is a pseudo-class for working with complex templates
// represents a single page to render
type Page struct {
	t    *template.Template
	path string
	name string
}

// New - initialization of complex template
// it clones root layout and adds one custom page to it
func (p *Page) New(source *template.Template, path string, name string) {
	p.path = path
	p.name = name
	p.t = template.Must(source.Clone())
	p.t = template.Must(p.t.ParseFiles(path))
}

// Execute - runs current template with stored name and given Reader and Writer
func (p *Page) Execute(w http.ResponseWriter, data interface{}) {
	err := p.t.ExecuteTemplate(w, p.name, data)
	if err != nil {
		log.Println(err)
	}
}

// Pages is a storage of complex templates
// it stores a collection of pages as map with human-readable names
type Pages struct {
	Items map[string]*Page
	root  *template.Template
}

// New - generates the root page layout template with given paths
// and the empty map of custom pages within
func (t *Pages) New(paths []string) {
	t.root = template.Must(template.ParseFiles(paths...))
	t.Items = make(map[string]*Page)
}

// Add - right way to add another complex template in our set
// adds a custom template to the page set and gives a human-friendly name to it
func (t *Pages) Add(path string, n string) {
	name := filepath.Base(path)
	t.Items[n] = &Page{}
	t.Items[n].New(t.root, path, name)
}
