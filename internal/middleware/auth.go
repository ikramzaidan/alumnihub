package middleware

import (
	"context"
	"net/http"

	"alumnihub/internal/auth"
)

func AuthRequired(authSvc *auth.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := authSvc.GetTokenFromHeaderAndVerify(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, auth.UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminRequired(authSvc *auth.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := authSvc.GetTokenFromHeaderAndVerify(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !claims.IsAdmin {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, auth.UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
