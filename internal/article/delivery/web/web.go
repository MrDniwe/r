package delivery

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/article/usecase"
	"github.com/mrdniwe/r/internal/models"
	"github.com/mrdniwe/r/internal/server"
	"github.com/mrdniwe/r/internal/view"
	e "github.com/mrdniwe/r/pkg/errors"
	"github.com/mrdniwe/r/pkg/templator"
	"github.com/sirupsen/logrus"
)

type ArticleDelivery struct {
	Usecase usecase.ArticleUsecase
	T       *templator.Pages
	Srv     *server.Server
}

func NewDelivery(uc usecase.ArticleUsecase, r *mux.Router, srv *server.Server) {
	view := view.New()
	ad := &ArticleDelivery{uc, view, srv}
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
		http.ServeFile(w, r, "template/static/img/favicon.ico")
	}
}

func (ad *ArticleDelivery) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := ad.Usecase.LastArticles(ad.Srv.Conf.GetInt("pageAmount"), 0)
		if err != nil {
			ad.Srv.Logger.Info("err in last")
			e.HandleError(err, w, r)
			return
		}
		topArticle := articles[0]
		total, err := ad.Usecase.TotalPagesCount()
		if err != nil {
			ad.Srv.Logger.Info("err in countpage")
			e.HandleError(err, w, r)
			return
		}
		restArticles := models.ArticleList{
			articles[1:],
			total,
			1,
		}
		page := models.Page{topArticle.Header, topArticle.Lead}
		mp := &models.HomePage{page, topArticle, restArticles}
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
		page := models.Page{
			a.Header,
			a.Lead,
		}
		pp := models.ArticlePage{
			page,
			a,
		}
		ad.T.Items["post"].Execute(w, pp)
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
		pAmount := ad.Srv.Conf.GetInt("pageAmount")
		pNum, err := strconv.Atoi(vars["page"])
		if err != nil {
			ad.Srv.Logger.WithFields(logrus.Fields{
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
		page := models.Page{
			"Список статей, страница " + vars["page"],
			"Список статей, страница " + vars["page"],
		}
		list := models.ArticleList{
			articles,
			total,
			pNum,
		}
		lp := &models.ListPage{
			page,
			list,
		}
		ad.T.Items["list"].Execute(w, lp)
	}
}
