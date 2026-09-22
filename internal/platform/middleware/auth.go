package middleware

import (
	"context"
	"ecommerce-api/internal/platform/security"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

func Authenticate(jwtManager *security.JWTManager) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHead := r.Header.Get("Authorization")
			if authHead == "" {
				http.Error(w, "authorization header required", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHead, "Bearer ") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHead, "Bearer ")

			claims, err := jwtManager.ValidateAccessToken(tokenString)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				claims.UserID,
			)

			ctx = context.WithValue(
				ctx,
				RoleKey,
				claims.Role,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
