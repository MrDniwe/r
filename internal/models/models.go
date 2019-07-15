package models

import (
	"html/template"
	"time"
)

type Site struct {
	Header string
	Lead   string
}

type Article struct {
	Id      string
	Visible bool
	Header  string
	Lead    string
	Text    template.HTML
	Date    time.Time
	Photo   string
	Views   int
}

type MainPage struct {
	Rest []*Article
	Article
}
