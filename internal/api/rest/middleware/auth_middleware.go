// middleware/auth_middleware.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/httputil"
	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/internal/lib"
)

func AuthMiddleware(jwtManager *lib.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				status, message := AuthErrorHandler(r, apperrors.ErrAuthMissingToken)
				httputil.RespondWithError(w, status, message)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				status, message := AuthErrorHandler(r, apperrors.ErrAuthInvalidTokenFormat)
				httputil.RespondWithError(w, status, message)
				return
			}

			tokenString := parts[1]

			// Use JWTManager to validate token
			claims, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				// Since ValidateToken returns specific apperrors, we can use them directly
				status, message := AuthErrorHandler(r, err)
				httputil.RespondWithError(w, status, message)
				return
			}

			// Extract user information from claims
			userID, ok := claims["user_id"]
			if !ok {
				status, message := AuthErrorHandler(r, apperrors.ErrAuthTokenInvalid)
				httputil.RespondWithError(w, status, message)
				return
			}

			// Create new Context
			ctx := context.WithValue(r.Context(), "userID", userID)

			// Add user info to context
			if email, ok := claims["email"].(string); ok {
				ctx = context.WithValue(r.Context(), "email", email)
			}
			if username, ok := claims["username"].(string); ok {
				ctx = context.WithValue(r.Context(), "username", username)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
