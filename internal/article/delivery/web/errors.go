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
			page := models.ErrorPage{
				models.Page{"Heкорректные параметры запроса", "Heкорректные параметры запроса"},
				"Ошибка 400",
				"Некорректные параметры запроса",
			}
			ad.T.Items["error"].Execute(w, page)
		case "notfound":
			page := models.ErrorPage{
				models.Page{"Страница не найдена", "Страница не найдена"},
				"Ошибка 404",
				"Ничего не найдено",
			}
			ad.T.Items["error"].Execute(w, page)
		case "server":
			page := models.ErrorPage{
				models.Page{"Ошибка сервера", "Проблемы на стороне сервера"},
				"Ошибка 500",
				"Проблемы на стороне сервера",
			}
			ad.T.Items["error"].Execute(w, page)
		case "forbidden":
			page := models.ErrorPage{
				models.Page{"Доступ запрещен", "Доступ на данную страницу для вас запрещен"},
				"Ошибка 403",
				"Доступ на данную страницу вам запрещён",
			}
			ad.T.Items["error"].Execute(w, page)
		default:
			page := models.ErrorPage{
				models.Page{"Неизвестная ошибка", "Неизвестная ошибка"},
				"Ошибка",
				"Неизвестная ошибка",
			}
			ad.T.Items["error"].Execute(w, page)
		}
	}
}

// метод для замены дефолтной страницы 404
func (ad *ArticleDelivery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := models.ErrorPage{
		models.Page{"Страница не найдена", "Страница не найдена"},
		"Ошибка 404",
		"Ничего не найдено",
	}
	ad.T.Items["error"].Execute(w, page)
}
