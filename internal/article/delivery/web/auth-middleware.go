package delivery

import (
	"context"
	"net/http"
)

func (ad *ArticleDelivery) CookieAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "isAuth", false)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (ad *ArticleDelivery) DrawIsAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ad.Srv.Logger.Info(ctx.Value("isAuth"))
		next.ServeHTTP(w, r)
	})
}

func isAuthorized(r *http.Request) bool {
	ctx := r.Context()
	return ctx.Value("isAuth").(bool)
}
