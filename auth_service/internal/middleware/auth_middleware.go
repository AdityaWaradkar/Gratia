package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

/*
	Context keys are unexported to avoid collisions
*/
type contextKey string

const (
	ctxUserIDKey    contextKey = "userID"
	ctxUserEmailKey contextKey = "userEmail"
	ctxUserRoleKey  contextKey = "userRole"
)

type errorResponse struct {
	Message string `json:"message"`
}

/*
	Auth middleware validates JWT and injects user data into context
*/
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := extractBearerToken(r)
		if tokenString == "" {
			writeJSONError(w, "authorization token missing or invalid", http.StatusUnauthorized)
			return
		}

		claims, err := validateJWT(tokenString)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, ctxUserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, ctxUserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/* ===================== HELPERS ===================== */

type jwtClaims struct {
	UserID string
	Email  string
	Role   string
}

func extractBearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func validateJWT(tokenString string) (*jwtClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(config.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	userID, _ := mapClaims["sub"].(string)
	email, _ := mapClaims["email"].(string)
	role, _ := mapClaims["role"].(string)

	if userID == "" || email == "" || role == "" {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &jwtClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Message: message})
}

/* ===================== CONTEXT ACCESSORS ===================== */

func UserID(ctx context.Context) string {
	v, _ := ctx.Value(ctxUserIDKey).(string)
	return v
}

func UserEmail(ctx context.Context) string {
	v, _ := ctx.Value(ctxUserEmailKey).(string)
	return v
}

func UserRole(ctx context.Context) string {
	v, _ := ctx.Value(ctxUserRoleKey).(string)
	return v
}
