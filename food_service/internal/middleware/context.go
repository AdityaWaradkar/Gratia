package middleware

import "context"

// contextKey is the type for context keys to avoid string collisions
type contextKey string

const (
    UserIDKey   contextKey = "user_id"
    UserRoleKey contextKey = "user_role"
)

// UserID returns the authenticated user ID from context safely
func UserID(ctx context.Context) string {
    v, _ := ctx.Value(UserIDKey).(string)
    return v
}

// UserRole returns the authenticated user role from context safely
func UserRole(ctx context.Context) string {
    v, _ := ctx.Value(UserRoleKey).(string)
    return v
}