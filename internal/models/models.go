package models

import (
	"html/template"
	"time"
)

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
