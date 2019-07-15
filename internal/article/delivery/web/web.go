package delivery

import (
	"github.com/gorilla/mux"
	"github.com/mrdniwe/pkg/errors"
	"github.com/mrdniwe/r/internal/article/usecase"
	"github.com/mrdniwe/r/internal/models"
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
	// Static
	static := http.FileServer(http.Dir("template/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", static))
}

func (ad *ArticleDelivery) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := ad.Usecase.LastArticles(11)
		if err != nil {
			ad.L.Fatal(err)
		}
		topArticle := *articles[0]
		mp := models.MainPage{
			articles[1:],
			topArticle,
		}
		ad.T.Items["mainpage"].Execute(w, mp)
	}
}

func (ad *ArticleDelivery) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		a, err := ad.Usecase.SingleArticle(vars["id"])
		if err != nil {
			errors.HandleError(err, w, r)
			return
		}
		ad.T.Items["post"].Execute(w, a)
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
