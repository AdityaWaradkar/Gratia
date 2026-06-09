package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type errorResponse struct {
	Error string `json:"error"`
}

type ContextKey string

const (
	UserIDKey   ContextKey = "userID"
	UserRoleKey ContextKey = "role"
)

type Middleware struct {
	jwtSecret []byte
}

func NewMiddleware(jwtSecret string) *Middleware {
	return &Middleware{
		jwtSecret: []byte(jwtSecret),
	}
}

/*
RequireAuth

Validates JWT and injects:

- user id
- user role

into request context.
*/

func (m *Middleware) RequireAuth(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(
				w,
				http.StatusUnauthorized,
				"authorization header missing",
			)
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 ||
			strings.ToLower(parts[0]) != "bearer" {

			writeError(
				w,
				http.StatusUnauthorized,
				"invalid authorization header",
			)
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(
			tokenStr,
			func(token *jwt.Token) (interface{}, error) {

				// Ensure HMAC signing method.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(
						"unexpected signing method: %v",
						token.Header["alg"],
					)
				}

				return m.jwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			writeError(
				w,
				http.StatusUnauthorized,
				"invalid or expired token",
			)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(
				w,
				http.StatusUnauthorized,
				"invalid token claims",
			)
			return
		}

		userID, _ := claims["sub"].(string)
		role, _ := claims["role"].(string)

		if userID == "" || role == "" {
			writeError(
				w,
				http.StatusUnauthorized,
				"invalid token claims",
			)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			userID,
		)

		ctx = context.WithValue(
			ctx,
			UserRoleKey,
			role,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

/*
Helpers
*/

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(
		errorResponse{
			Error: message,
		},
	)
}