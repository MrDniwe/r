package delivery

import (
	"context"
	"net/http"
)

func (ad *ArticleDelivery) CookieAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// читаем кукизы
		authData, err := ad.ReadAuthCookies(r)
		if err != nil {
			ctx = context.WithValue(ctx, "access_token", false)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}
		// если есть - сверяем в бд
		if err = ad.Usecase.CheckAccessToken(authData.AccessToken); err == nil {
			ctx = context.WithValue(ctx, "access_token", authData.AccessToken)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}
		// если устарели - пробуем рефрешнуть
		authData, err = ad.Usecase.RefreshToken(authData.RefreshToken)
		if err != nil {
			ctx = context.WithValue(ctx, "access_token", false)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return

		}
		ad.SetAuthCookies(w, authData)
		ctx = context.WithValue(ctx, "access_token", authData.AccessToken)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}

func isAuthorized(r *http.Request) bool {
	ctx := r.Context()
	switch ctx.Value("access_token").(type) {
	case string:
		return true
	default:
		return false
	}
}
