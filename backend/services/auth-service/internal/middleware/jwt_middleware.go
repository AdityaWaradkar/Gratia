package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/adityawaradkar/gratia-auth/internal/utils"
)

// ContextKey is a custom type to avoid context key collisions
type ContextKey string

const (
	UserIDKey  ContextKey = "user_id"
	UserRoleKey ContextKey = "user_role"
)

// JWTAuthMiddleware validates JWT tokens from Authorization header and
// stores user ID and role into request context.
func JWTAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := utils.ValidateJWT(tokenStr, secret)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok || sub == "" {
				http.Error(w, "Invalid token: subject claim missing", http.StatusUnauthorized)
				return
			}

			role, _ := claims["role"].(string) // Role may be empty or missing

			// Store user ID and role in context for handlers to use
			ctx := context.WithValue(r.Context(), UserIDKey, sub)
			ctx = context.WithValue(ctx, UserRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
