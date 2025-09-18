package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Context keys for storing user info
type contextKey string

const (
	ContextUserID    contextKey = "userID"
	ContextUserEmail contextKey = "userEmail"
	ContextUserRole  contextKey = "userRole"
)

// JSONError represents a standard error response
type JSONError struct {
	Message string `json:"message"`
}

// AuthMiddleware validates JWT access tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, "authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeError(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		// Parse JWT token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(config.GetJWTSecret()), nil
		})
		if err != nil || !token.Valid {
			writeError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok1 := claims["sub"].(string)
		userEmail, ok2 := claims["email"].(string)
		userRole, ok3 := claims["role"].(string)
		if !ok1 || !ok2 || !ok3 {
			writeError(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		// Store in request context
		ctx := context.WithValue(r.Context(), ContextUserID, userID)
		ctx = context.WithValue(ctx, ContextUserEmail, userEmail)
		ctx = context.WithValue(ctx, ContextUserRole, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper functions to retrieve user info from context
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(ContextUserID).(string); ok {
		return id
	}
	return ""
}

func GetUserEmail(ctx context.Context) string {
	if email, ok := ctx.Value(ContextUserEmail).(string); ok {
		return email
	}
	return ""
}

func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(ContextUserRole).(string); ok {
		return role
	}
	return ""
}

// Utility to return JSON errors
func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONError{Message: message})
}
