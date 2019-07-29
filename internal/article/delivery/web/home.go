package delivery

import (
	"net/http"

	"github.com/mrdniwe/r/internal/models"
	e "github.com/mrdniwe/r/pkg/errors"
)

func (ad *ArticleDelivery) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuth := isAuthorized(r)
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
		page := models.Page{topArticle.Header, topArticle.Lead, isAuth}
		mp := &models.HomePage{page, topArticle, restArticles}
		ad.T.Items["mainpage"].Execute(w, mp)
	}
}
