package delivery

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/models"
	e "github.com/mrdniwe/r/pkg/errors"
)

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
