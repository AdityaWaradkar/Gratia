package middleware

import (
    "context"
    "encoding/json"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
    UserIDKey    contextKey = "user_id"
    UserEmailKey contextKey = "user_email"
    UserRoleKey  contextKey = "user_role"
)

type errorResponse struct {
    Message string `json:"message"`
}

// Auth creates a middleware closure that injects the JWT secret
func Auth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
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

            token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
                if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, jwt.ErrTokenSignatureInvalid
                }
                return []byte(jwtSecret), nil
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
            roleStr, _ := claims["role"].(string)

            if userID == "" || email == "" || roleStr == "" {
                writeError(w, http.StatusUnauthorized, "incomplete token claims")
                return
            }

            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            ctx = context.WithValue(ctx, UserEmailKey, email)
            ctx = context.WithValue(ctx, UserRoleKey, roleStr)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func UserID(ctx context.Context) string {
    v, _ := ctx.Value(UserIDKey).(string)
    return v
}

func UserEmail(ctx context.Context) string {
    v, _ := ctx.Value(UserEmailKey).(string)
    return v
}

// UserRole now returns a standard string to prevent import cycles
func UserRole(ctx context.Context) string {
    v, _ := ctx.Value(UserRoleKey).(string)
    return v
}

func writeError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(errorResponse{Message: message})
}