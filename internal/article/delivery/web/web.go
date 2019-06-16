package delivery

import (
	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/article/usecase"
	"github.com/mrdniwe/r/internal/view"
	"github.com/mrdniwe/r/pkg/templator"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ArticleDelivery struct {
	Usecase usecase.ArticleUsecase
	L       *logrus.Logger
	T       *templator.Pages
}

func NewDelivery(uc usecase.ArticleUsecase, l *logrus.Logger, r *mux.Router) {
	view := view.New()
	ad := &ArticleDelivery{uc, l, view}
	r.HandleFunc("/", ad.Home()).Methods("GET")
	r.HandleFunc("/post/{id}", ad.Post()).Methods("GET")
	r.HandleFunc("/info/{page}", ad.Static()).Methods("GET")
}

func (ad *ArticleDelivery) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ad.T.Items["mainpage"].Execute(w, r)
	}
}

func (ad *ArticleDelivery) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ad.T.Items["post"].Execute(w, r)
	}
}

func (ad *ArticleDelivery) Static() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ad.T.Items["static"].Execute(w, r)
	}
}

func (ad *ArticleDelivery) Dummy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ad.T.Items["dummy"].Execute(w, r)
	}
}
