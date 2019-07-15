package delivery

import (
	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/models"
	"net/http"
)

func (ad *ArticleDelivery) ErrHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		errtype := vars["errtype"]
		switch errtype {
		case "badrequest":
			site := models.Site{"Heкорректные параметры запроса", "Heкорректные параметры запроса"}
			ad.T.Items["badrequest"].Execute(w, site)
		case "notfound":
			site := models.Site{"Страница не найдена", "Страница не найдена"}
			ad.T.Items["notfound"].Execute(w, site)
		case "server":
			site := models.Site{"Ошибка сервера", "Проблемы на стороне сервера"}
			ad.T.Items["server"].Execute(w, site)
		case "forbidden":
			site := models.Site{"Доступ запрещен", "Доступ на данную страницу для вас запрещен"}
			ad.T.Items["forbidden"].Execute(w, site)
		default:
			site := models.Site{"Неизвестная ошибка", "Неизвестная ошибка"}
			ad.T.Items["unknown"].Execute(w, site)
		}
	}
}
