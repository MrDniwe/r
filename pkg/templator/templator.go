package templator

import (
	"html/template"
	"log"
	"net/http"
)

// Page is a pseudo-class for working with complex templates
type Page struct {
	t    *template.Template
	path string
	name string
}

// New - initialization of complex template
func (p *Page) New(source *template.Template, path string, name string) {
	p.path = path
	p.name = name
	p.t = template.Must(source.Clone())
	p.t = template.Must(p.t.ParseFiles(path))
}

// complex template executor
func (p *Page) Execute(w http.ResponseWriter, r *http.Request) {
	err := p.t.ExecuteTemplate(w, p.name, r)
	if err != nil {
		log.Println(err)
	}
}

// Pages is a storage of complex templates
type Pages struct {
	Items map[string]*Page
	root  *template.Template
}

func (t *Pages) new() {
	t.root = template.Must(template.ParseFiles(
		"templates/header.html",
		"templates/footer.html",
		"templates/navigation.html",
		"templates/page-header.html",
		"templates/page-footer.html",
		"templates/layout.html"))
	t.Items = make(map[string]*Page)
}

// right way to add another complex template in our set
func (t *Pages) add(path string, name string) {
	t.Items[name] = &Page{}
	t.Items[name].New(t.root, path, name)
}
