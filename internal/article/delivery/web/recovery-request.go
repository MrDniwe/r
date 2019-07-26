package delivery

import (
	"net/http"

	"github.com/mrdniwe/r/internal/models"
)

func (ad *ArticleDelivery) RecoveryRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := models.Page{
			"Восстановление пароля",
			"Восстановление пароля зарегистрированного пользователя",
		}
		p := models.RecoveryPage{page, "", ""}
		ad.T.Items["recovery-request"].Execute(w, p)
	}
}

func (ad *ArticleDelivery) RecoveryRequestNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := models.Page{
			"Восстановление пароля (email не найден)",
			"Восстановление пароля зарегистрированного пользователя",
		}
		email := r.FormValue("email")
		p := models.RecoveryPage{page, "Пользователь с указанной почтой у нас ещё не зарегистрирован", email}
		ad.T.Items["recovery-request"].Execute(w, p)
	}
}

func (ad *ArticleDelivery) RecoveryRequestDelay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := models.Page{
			"Восстановление пароля (слишком частые запросы)",
			"Восстановление пароля зарегистрированного пользователя",
		}
		email := r.FormValue("email")
		p := models.RecoveryPage{page, "Вы недавно отправляли запрос на восстановление пароля, новый можно будет отправить только через минуту", email}
		ad.T.Items["recovery-request"].Execute(w, p)
	}
}
