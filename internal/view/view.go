package view

import "github.com/mrdniwe/r/pkg/templator"

func New() *templator.Pages {
	p := &templator.Pages{}
	root := []string{"template/layout/header.html",
		"template/layout/footer.html",
		"template/layout/navigation.html",
		"template/layout/page-header.html",
		"template/layout/page-footer.html",
		"template/layout/navbar-placeholder.html",
		"template/layout/comments.html",
		"template/layout/layout.html"}

	p.New(root)
	p.Add("template/pages/home.html", "mainpage")
	p.Add("template/pages/post.html", "post")
	p.Add("template/pages/static.html", "static")
	p.Add("template/pages/list.html", "list")
	p.Add("template/errors/error.html", "error")
	return p
}
