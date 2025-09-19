package auth

import "time"

// User represents a user account
type User struct {
	ID            string    `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	Role          string    `json:"role" db:"role"`
	EmailVerified bool      `json:"emailVerified" db:"email_verified"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

// RefreshToken represents a refresh token for session persistence
type RefreshToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	TokenHash string    `json:"-" db:"token"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	Revoked   bool      `json:"revoked" db:"revoked"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// Session tracks user logins across devices
type Session struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"userId" db:"user_id"`
	RefreshTokenID string    `json:"refreshTokenId" db:"refresh_token_id"`
	UserAgent      string    `json:"userAgent" db:"user_agent"`
	IPAddress      string    `json:"ipAddress" db:"ip_address"`
	IsCurrent      bool      `json:"isCurrent" db:"is_current"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
}

// EmailVerification stores email confirmation tokens
type EmailVerification struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Token     string    `json:"-" db:"token"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// PasswordReset stores password reset tokens
type PasswordReset struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Token     string    `json:"-" db:"token"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// AuditLog stores security/audit events
type AuditLog struct {
	ID        string    `json:"id" db:"id"`
	UserID    *string   `json:"userId,omitempty" db:"user_id"`
	Action    string    `json:"action" db:"action"`
	IPAddress string    `json:"ipAddress" db:"ip_address"`
	UserAgent string    `json:"userAgent" db:"user_agent"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
