package models

import "html/template"

type Article struct {
	Id     string
	Cat    int
	Access int
	Hidden int
	Header string
	Pre    string
	Text   template.HTML
	Date   int64
	Photo  string
	Views  int
	Yandex int
}

type MainPage struct {
	Rest []*Article
	Article
}
