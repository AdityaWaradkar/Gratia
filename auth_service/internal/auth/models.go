package auth

import "time"

/*
	User represents an authenticated account in the system.
*/
type User struct {
	ID            string    `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	Role          string    `json:"role" db:"role"`
	EmailVerified bool      `json:"emailVerified" db:"email_verified"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

/*
	RefreshToken is used to issue new access tokens.
*/
type RefreshToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	TokenHash string    `json:"-" db:"token"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	Revoked   bool      `json:"revoked" db:"revoked"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

/*
	Session tracks active login sessions per user.
*/
type Session struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"userId" db:"user_id"`
	RefreshTokenID string    `json:"refreshTokenId" db:"refresh_token_id"`
	UserAgent      string    `json:"userAgent" db:"user_agent"`
	IPAddress      string    `json:"ipAddress" db:"ip_address"`
	IsCurrent      bool      `json:"isCurrent" db:"is_current"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
}

/*
	EmailVerification is used to verify user email addresses.
*/
type EmailVerification struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Token     string    `json:"-" db:"token"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

/*
	PasswordReset handles forgot-password flows.
*/
type PasswordReset struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Token     string    `json:"-" db:"token"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

/*
	AuditLog records security-sensitive actions.
	(Optional / future-use, safe to keep.)
*/
type AuditLog struct {
	ID        string    `json:"id" db:"id"`
	UserID    *string   `json:"userId,omitempty" db:"user_id"`
	Action    string    `json:"action" db:"action"`
	IPAddress string    `json:"ipAddress" db:"ip_address"`
	UserAgent string    `json:"userAgent" db:"user_agent"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
