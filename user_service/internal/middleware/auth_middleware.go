package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/adityawaradkar/gratia/user_service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
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

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		tokenStr := parts[1]

		claims := &jwtClaims{}

		token, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(token *jwt.Token) (interface{}, error) {

				// Ensure HMAC signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}

				return []byte(config.AppConfig.JWTSecret), nil
			},
		)

		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		if claims.Subject == "" || claims.Role == "" {
			writeError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			claims.Subject,
		)

		ctx = context.WithValue(
			ctx,
			UserRoleKey,
			claims.Role,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}