package view

import "github.com/mrdniwe/r/pkg/templator"

func New() *templator.Pages {
	p := &templator.Pages{}
	root := []string{"./internal/views/layout/header.html",
		"./internal/views/layout/footer.html",
		"./internal/views/layout/navigation.html",
		"./internal/views/layout/page-header.html",
		"./internal/views/layout/page-footer.html",
		"./internal/views/layout/layout.html"}
	p.New(root)
	p.Add("./internal/views/pages/home.html", "home.html")
	return p
}
