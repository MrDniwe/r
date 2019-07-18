package delivery

import (
	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/article/usecase"
	"github.com/mrdniwe/r/internal/models"
	"github.com/mrdniwe/r/internal/view"
	e "github.com/mrdniwe/r/pkg/errors"
	"github.com/mrdniwe/r/pkg/templator"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

type ArticleDelivery struct {
	Usecase usecase.ArticleUsecase
	L       *logrus.Logger
	T       *templator.Pages
	V       *viper.Viper
}

func NewDelivery(uc usecase.ArticleUsecase, l *logrus.Logger, r *mux.Router, v *viper.Viper) {
	view := view.New()
	ad := &ArticleDelivery{uc, l, view, v}
	r.NotFoundHandler = ad
	r.HandleFunc("/", ad.Home()).Methods("GET")
	r.HandleFunc("/post/{id}", ad.Post()).Methods("GET")
	r.HandleFunc("/info/{page}", ad.Static()).Methods("GET")
	r.HandleFunc("/list/{page}", ad.List()).Methods("GET")
	r.HandleFunc("/favicon.ico", ad.Favicon()).Methods("GET")
	// errors
	errh := r.PathPrefix("/errors/").Subrouter()
	errh.HandleFunc("/{errtype}", ad.ErrHandler())
	// Static
	static := http.FileServer(http.Dir("template/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", static))
}

func (ad *ArticleDelivery) Favicon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "template/static/favicon.ico")
	}
}

func (ad *ArticleDelivery) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := ad.Usecase.LastArticles(ad.V.GetInt("pageAmount"), 0)
		if err != nil {
			e.HandleError(err, w, r)
			return
		}
		topArticle := *articles[0]
		total, err := ad.Usecase.TotalPagesCount()
		if err != nil {
			e.HandleError(err, w, r)
			return
		}
		mp := &models.ListPage{
			articles[1:],
			topArticle,
			total,
			1,
		}
		ad.T.Items["mainpage"].Execute(w, mp)
	}
}

func (ad *ArticleDelivery) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		a, err := ad.Usecase.SingleArticle(vars["id"])
		if err != nil {
			e.HandleError(err, w, r)
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

func (ad *ArticleDelivery) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pAmount := ad.V.GetInt("pageAmount")
		pNum, err := strconv.Atoi(vars["page"])
		if err != nil {
			ad.L.WithFields(logrus.Fields{
				"type":  e.ValidationError,
				"in":    "Page number",
				"given": vars["page"],
			}).Info(err)
			e.HandleError(e.BadRequestErr, w, r)
			return
		}
		total, err := ad.Usecase.TotalPagesCount()
		if err != nil {
			e.HandleError(err, w, r)
			return
		}
		if pNum > total {
			e.HandleError(e.NotFoundErr, w, r)
			return
		}
		offset := pAmount * (pNum - 1)
		articles, err := ad.Usecase.LastArticles(pAmount, offset)
		if err != nil {
			e.HandleError(err, w, r)
			return
		}
		lp := &models.ListPage{
			articles,
			models.Article{
				Header: "Список статей, страница " + vars["page"],
				Lead:   "Список статей, страница " + vars["page"],
			},
			total,
			pNum,
		}
		ad.T.Items["list"].Execute(w, lp)
	}
}
