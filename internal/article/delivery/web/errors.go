package delivery

import (
	"github.com/mrdniwe/r/internal/models"
	"net/http"
)

func (ad *ArticleDelivery) NotFoundErr() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		site := models.Site{"Страница не найдена", "Страница не найдена"}
		ad.T.Items["notfound"].Execute(w, site)
	}
}

func (ad *ArticleDelivery) BadRequestErr() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		site := models.Site{"Некорректные параметры запроса", "Некорректные параметры запроса"}
		ad.T.Items["badrequest"].Execute(w, site)
	}
}

//TODO остальные ошибки
