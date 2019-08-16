package models

import (
	"html/template"
	"time"
)

// Page is a set of common params for all viwable HTML template pages
type Page struct {
	Title        string
	Description  string
	IsAuthorized bool
}

// ArticleList represents a slice of articles with params for drawing pagination
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

// SignIn - data set for rendering sign in page with a probable error
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

// AuthData represents a combination of tokens for user to be successfully authorized
type AuthData struct {
	AccessToken  string
	RefreshToken string
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
