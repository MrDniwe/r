package delivery

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/article/usecase"
	"github.com/mrdniwe/r/internal/server"
	"github.com/mrdniwe/r/internal/view"
	"github.com/mrdniwe/r/pkg/templator"
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
	r.HandleFunc("/recovery-request", ad.RecoveryRequest()).Methods("GET")
	r.HandleFunc("/recovery-request-notfound", ad.RecoveryRequestNotFound()).Methods("GET")
	r.HandleFunc("/recovery-request-delay", ad.RecoveryRequestDelay()).Methods("GET")
	r.HandleFunc("/recovery-submit", ad.RecoverySubmit()).Methods("POST")
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

func (ad *ArticleDelivery) Static() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ad.T.Items["static"].Execute(w, r)
	}
}
