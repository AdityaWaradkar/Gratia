package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Context keys for authenticated user
type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UserRoleKey  contextKey = "user_role"
)

// errorResponse represents auth error
type errorResponse struct {
	Message string `json:"message"`
}

// Auth validates JWT and injects user context
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusUnauthorized, "authorization header missing")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		if userID == "" || email == "" || role == "" {
			writeError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserEmailKey, email)
		ctx = context.WithValue(ctx, UserRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserID returns user id from context
func UserID(ctx context.Context) string {
	v, _ := ctx.Value(UserIDKey).(string)
	return v
}

// UserEmail returns email from context
func UserEmail(ctx context.Context) string {
	v, _ := ctx.Value(UserEmailKey).(string)
	return v
}

// UserRole returns role from context
func UserRole(ctx context.Context) string {
	v, _ := ctx.Value(UserRoleKey).(string)
	return v
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Message: message})
}
