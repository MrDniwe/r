package models

import (
	"html/template"
	"time"
)

type Page struct {
	Title        string
	Description  string
	IsAuthorized bool
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
	Comments Comments
}

type Comment struct {
	Id     string    `json:"uuid"`
	UserId string    `json:"user_uuid"`
	Text   string    `json:"message"`
	Date   time.Time `json:"created_at"`
	User   User
}

type Comments []Comment

func (c Comments) Presents() bool {
	return len(c) > 0
}

type User struct {
	Id    string `json:"uuid"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type EmptyPage struct {
	Page Page
}

type RecoveryPage struct {
	Page  Page
	Error string
	Email string
}

type SubmitPage struct {
	Page  Page
	Error string
	Code  string
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

type SignIn struct {
	Email         string
	EmailError    string
	PasswordError string
}

type SignInPage struct {
	Page   Page
	SignIn SignIn
}

type RecoveryData struct {
	Login string
	Email string
	Code  string
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
