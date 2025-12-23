package middleware

import "context"

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

// UserID returns authenticated user id from context
func UserID(ctx context.Context) string {
	v, _ := ctx.Value(UserIDKey).(string)
	return v
}

// UserRole returns authenticated user role from context
func UserRole(ctx context.Context) string {
	v, _ := ctx.Value(UserRoleKey).(string)
	return v
}
