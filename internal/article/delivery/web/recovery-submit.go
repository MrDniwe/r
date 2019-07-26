package delivery

import (
	"fmt"
	"net/http"

	"github.com/mrdniwe/r/internal/models"
	"github.com/mrdniwe/r/pkg/errors"
)

func (ad *ArticleDelivery) RecoverySubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// проверяем наличие, валидность и наличие емейла в БД
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/errors/server", http.StatusMovedPermanently)
		}
		email := r.Form.Get("email")
		exists, err := ad.Usecase.CheckEmailExists(email)
		if err != nil {
			errors.HandleError(err, w, r)
		}
		if !exists {
			http.Redirect(w, r, fmt.Sprintf("/recovery-request-notfound?email=%v", email), http.StatusMovedPermanently)
			return
		}
		// пробуем сгенерировать код и отправить на почту
		var userdata models.RecoveryData
		userdata, err = ad.Usecase.NewRecoveryHash(email)
		if err != nil {
			switch err {
			case errors.DelayErr:
				http.Redirect(w, r, "/recovery-request-delay", http.StatusMovedPermanently)
				return
			default:
				errors.HandleError(err, w, r)
				return
			}
		}
		err = ad.Srv.Mailer.SendRecovery(userdata.Login, userdata.Email, userdata.Code)
		if err != nil {
			errors.HandleError(err, w, r)
			return
		}
		page := models.Page{
			"Восстановление пароля",
			"Восстановление пароля зарегистрированного пользователя",
		}
		p := models.SubmitPage{page, "", ""}
		ad.T.Items["recovery-submit"].Execute(w, p)
	}
}
