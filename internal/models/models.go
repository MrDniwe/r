package models

import (
	"html/template"
	"time"
)

type Page struct {
	Title       string
	Description string
}

type ArticleList struct {
	Articles    []*Article
	TotalPages  int
	CurrentPage int
}

type Article struct {
	Id       string
	Visible  bool
	Header   string
	Lead     string
	Text     template.HTML
	Date     time.Time
	Photo    string
	Views    int
	Comments []Comment
}

type Comment struct {
	Id     string
	UserId string
	Text   string
	Date   time.Time
	User   User
}

type User struct {
	Id    string
	Login string
	Email string
}

type ListPage struct {
	Page Page
	List ArticleList
}

type HomePage struct {
	Page    Page
	Article *Article
	Rest    ArticleList
}

type ArticlePage struct {
	Page    Page
	Article *Article
}

type ErrorPage struct {
	Page             Page
	ErrorTitle       string
	ErrorDescription string
}

func (lp *ArticleList) NextPage() int {
	return lp.CurrentPage + 1
}

func (lp *ArticleList) PrevPage() int {
	return lp.CurrentPage - 1
}

func (lp *ArticleList) HasPrev() bool {
	return lp.CurrentPage > 1
}

func (lp *ArticleList) HasNext() bool {
	return lp.CurrentPage < lp.TotalPages
}
