package middleware

import (
	"context"
	"net/http"
	"strings"
)

type AuthOAuth interface {
	ContextAccessToken(ctx context.Context, token string) context.Context
}

func Auth(oa AuthOAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return authBearer(oa)(next)
	}
}

func authBearer(oa AuthOAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			bToken, bExist := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if bExist && bToken != "" {
				r = r.WithContext(oa.ContextAccessToken(ctx, bToken))
			}

			next.ServeHTTP(w, r)
		})
	}
}
