package middleware

import (
    "context"
    "net/http"
)

type contextKey string

const (
    UserIDKey    contextKey = "user_id"
    UserRoleKey  contextKey = "user_role"
)

// UserID extracts and returns the authenticated user id from the request context
func UserID(ctx context.Context) string {
    v, _ := ctx.Value(UserIDKey).(string)
    return v
}

// UserRole extracts and returns the authenticated user role from the request context
func UserRole(ctx context.Context) string {
    v, _ := ctx.Value(UserRoleKey).(string)
    return v
}

// RequireRole enforces that the caller possesses a specific role
func RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userRole := UserRole(r.Context())
            if userRole == "" {
                writeError(w, http.StatusUnauthorized, "unauthorized")
                return
            }
            if userRole != role {
                writeError(w, http.StatusForbidden, "forbidden")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// RequireAnyRole enforces that the caller possesses at least one of the permitted roles
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
    roleSet := make(map[string]struct{})
    for _, r := range roles {
        roleSet[r] = struct{}{}
    }

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userRole := UserRole(r.Context())
            if userRole == "" {
                writeError(w, http.StatusUnauthorized, "unauthorized")
                return
            }
            if _, ok := roleSet[userRole]; !ok {
                writeError(w, http.StatusForbidden, "forbidden")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}