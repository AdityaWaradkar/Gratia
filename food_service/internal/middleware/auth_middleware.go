package middleware

import (
    "context"
    "encoding/json"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type errorResponse struct {
    Message string `json:"message"`
}

type jwtClaims struct {
    Role string `json:"role"`
    jwt.RegisteredClaims
}

// Auth validates JWT and injects user context into the request using dependency injection for the secret
func Auth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                writeError(w, http.StatusUnauthorized, "authorization header missing")
                return
            }

            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
                writeError(w, http.StatusUnauthorized, "invalid authorization header format")
                return
            }

            tokenStr := parts[1]
            claims := &jwtClaims{}

            token, err := jwt.ParseWithClaims(
                tokenStr,
                claims,
                func(token *jwt.Token) (interface{}, error) {
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, jwt.ErrTokenSignatureInvalid
                    }
                    return []byte(jwtSecret), nil
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

            ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
            ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// writeError writes a structured JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(errorResponse{
        Message: message,
    })
}