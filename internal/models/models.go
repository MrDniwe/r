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

type ListPage struct {
	Rest []*Article
	Article
	TotalPages  int
	CurrentPage int
}

func (lp *ListPage) NextPage() int {
	return lp.CurrentPage + 1
}

func (lp *ListPage) PrevPage() int {
	return lp.CurrentPage - 1
}

func (lp *ListPage) HasPrev() bool {
	return lp.CurrentPage > 1
}

func (lp *ListPage) HasNext() bool {
	return lp.CurrentPage < lp.TotalPages
}
