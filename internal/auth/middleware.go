package auth

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
)

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("auth-cookie")
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}
			userId, err := validateAndGetUserID(c)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), UserCtxKey, userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func WrapGraphQLHandler(srv *handler.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RespWriter, w)
		srv.ServeHTTP(w, r.WithContext(ctx))
	})
}
