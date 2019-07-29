package delivery

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/models"
)

func (ad *ArticleDelivery) ErrHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		errtype := vars["errtype"]
		isAuth := isAuthorized(r)
		switch errtype {
		case "badrequest":
			page := models.Page{"Heкорректные параметры запроса", "Heкорректные параметры запроса", isAuth}
			o := models.ErrorPage{
				page,
				"Ошибка 400",
				"Некорректные параметры запроса",
			}
			ad.T.Items["error"].Execute(w, o)
		case "notfound":
			page := models.ErrorPage{
				models.Page{"Страница не найдена", "Страница не найдена", isAuth},
				"Ошибка 404",
				"Ничего не найдено",
			}
			ad.T.Items["error"].Execute(w, page)
		case "server":
			page := models.ErrorPage{
				models.Page{"Ошибка сервера", "Проблемы на стороне сервера", isAuth},
				"Ошибка 500",
				"Проблемы на стороне сервера",
			}
			ad.T.Items["error"].Execute(w, page)
		case "forbidden":
			page := models.ErrorPage{
				models.Page{"Доступ запрещен", "Доступ на данную страницу для вас запрещен", isAuth},
				"Ошибка 403",
				"Доступ на данную страницу вам запрещён",
			}
			ad.T.Items["error"].Execute(w, page)
		default:
			page := models.ErrorPage{
				models.Page{"Неизвестная ошибка", "Неизвестная ошибка", isAuth},
				"Ошибка",
				"Неизвестная ошибка",
			}
			ad.T.Items["error"].Execute(w, page)
		}
	}
}

// метод для замены дефолтной страницы 404
func (ad *ArticleDelivery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isAuth := isAuthorized(r)
	page := models.ErrorPage{
		models.Page{"Страница не найдена", "Страница не найдена", isAuth},
		"Ошибка 404",
		"Ничего не найдено",
	}
	ad.T.Items["error"].Execute(w, page)
}
