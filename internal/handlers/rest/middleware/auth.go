package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/auth"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
)

func AuthMiddleware(tokenService auth.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				HandleError(w, r, apperrors.ErrAuthMissingToken)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				HandleError(w, r, apperrors.ErrAuthInvalidTokenFormat)
				return
			}

			tokenString := parts[1]
			claims, err := tokenService.ValidateToken(tokenString)
			if err != nil {
				HandleError(w, r, err)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "userID", claims.UserID)
			ctx = context.WithValue(ctx, "email", claims.Email)
			ctx = context.WithValue(ctx, "username", claims.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
