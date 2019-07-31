package delivery

import (
	"net/http"

	"github.com/mrdniwe/r/internal/models"
)

func (ad *ArticleDelivery) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuth := isAuthorized(r)
		if isAuth {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		page := models.Page{
			"Аутентификация",
			"Аутентификация - вход в профиль ранее зарегистрированного пользователя",
			isAuth,
		}
		var signIn models.SignIn
		// если был POST
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Redirect(w, r, "/errors/server", http.StatusTemporaryRedirect)
			}
			email := r.Form.Get("email")
			email := r.Form.Get("password")
			// если всё ок, ставим куки и редиректаем на главную и выходим
			auth, err := ad.Usecase.UserAuth(email, password)
			// если не ок, заполняем ошибки и отрисовываем форму и выходим
		}
		// отрисовываем пустую форму
		signIn = models.SignIn{}
		p := models.SignInPage{page, signIn}
		ad.T.Items["sign-in"].Execute(w, p)
	}
}
