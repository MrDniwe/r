package delivery

import (
	"net/http"
	"time"

	"github.com/mrdniwe/r/internal/models"
	e "github.com/mrdniwe/r/pkg/errors"
)

// SignIn - a handler for sigin-in page
func (ad *ArticleDelivery) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuth := isAuthorized(r)
		if isAuth {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		page := models.Page{
			"Аутентификация",
			"Аутентификация - вход в профиль ранее зарегистрированного пользователя",
			isAuth,
		}
		var signIn models.SignIn
		signIn = models.SignIn{}
		// если был POST
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				e.HandleError(e.ServerErr, w, r)
				return
			}
			email := r.Form.Get("email")
			password := r.Form.Get("password")
			// если всё ок, ставим куки и редиректаем на главную и выходим
			auth, err := ad.Usecase.UserAuth(email, password)
			if err == nil {
				ad.SetAuthCookies(w, auth)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			// если не ок, заполняем ошибки и отрисовываем форму и выходим
			switch err {
			case e.NotFoundErr:
				signIn = models.SignIn{
					Email:         email,
					EmailError:    "Указанный email не найден",
					PasswordError: "",
				}
			case e.WrongPasswordErr:
				signIn = models.SignIn{
					Email:         email,
					EmailError:    "",
					PasswordError: "Указан неверный пароль",
				}
			default:
				e.HandleError(err, w, r)
				return
			}
		}
		// отрисовываем форму
		p := models.SignInPage{Page: page, SignIn: signIn}
		ad.T.Items["sign-in"].Execute(w, p)
	}
}

// SetAuthCookies - creates auth cookies for user
func (ad *ArticleDelivery) SetAuthCookies(w http.ResponseWriter, auth models.AuthData) {
	//TODO: добавить проверку домена из конфига
	accessTokenCookie := &http.Cookie{
		Name:     "r57AT",
		Value:    auth.AccessToken,
		Expires:  time.Now().Add(time.Minute * 5),
		MaxAge:   5 * 60,
		HttpOnly: true,
	}
	refreshTokenCookie := &http.Cookie{
		Name:     "r57RT",
		Value:    auth.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 180),
		MaxAge:   60 * 60 * 24 * 180,
		HttpOnly: true,
	}
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
}

// DropAuthCookies - removes auth cookies due logout
func (ad *ArticleDelivery) DropAuthCookies(w http.ResponseWriter) {
	//TODO: добавить проверку домена из конфига
	accessTokenCookie := &http.Cookie{
		Name:     "r57AT",
		Expires:  time.Now().Add(time.Minute * 5),
		MaxAge:   -1,
		HttpOnly: true,
	}
	refreshTokenCookie := &http.Cookie{
		Name:     "r57RT",
		Expires:  time.Now().Add(time.Hour * 24 * 180),
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
}

// SignOut - logs the user out
func (ad *ArticleDelivery) SignOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// если пользователь не авторизован, то игнорируем и перебрасываем на главную
		isAuth := isAuthorized(r)
		if !isAuth {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// ломимся в базу за логаутом
		// удаляем кукизы
	}
}

// ReadAuthCookies - tries to read auth cookies and returns tokens from there as AuthData
func (ad *ArticleDelivery) ReadAuthCookies(r *http.Request) (models.AuthData, error) {
	atCookie, err := r.Cookie("r57AT")
	if err != nil {
		return models.AuthData{}, err
	}
	rtCookie, err := r.Cookie("r57RT")
	if err != nil {
		return models.AuthData{}, err
	}
	return models.AuthData{
		AccessToken:  atCookie.Value,
		RefreshToken: rtCookie.Value,
	}, nil
}
