package delivery

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/models"
	e "github.com/mrdniwe/r/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (ad *ArticleDelivery) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		isAuth := isAuthorized(r)
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
			isAuth,
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
