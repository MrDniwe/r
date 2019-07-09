package models

import "html/template"

type Article struct {
	Id      string
	Visible bool
	Header  string
	Lead    string
	Text    template.HTML
	Date    int64
	Photo   string
	Views   int
}

type MainPage struct {
	Rest []*Article
	Article
}
