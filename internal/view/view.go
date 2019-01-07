package view

import "github.com/mrdniwe/r/pkg/templator"

func New() *templator.Pages {
	p := &templator.Pages{}
	root := []string{"./internal/view/layout/header.html",
		"./internal/view/layout/footer.html",
		"./internal/view/layout/navigation.html",
		"./internal/view/layout/page-header.html",
		"./internal/view/layout/page-footer.html",
		"./internal/view/layout/layout.html"}
	p.New(root)
	p.Add("./internal/view/pages/home.html", "mainpage")
	p.Add("./internal/view/pages/post.html", "post")
	p.Add("./internal/view/pages/static.html", "static")
	p.Add("./internal/view/pages/dummy.html", "dummy")
	return p
}
