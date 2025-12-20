package auth

import "time"

// User represents an authenticated account
type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}


// Session represents a login session
type Session struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	RefreshTokenID string    `json:"refreshTokenId"`
	UserAgent      string    `json:"userAgent"`
	IPAddress      string    `json:"ipAddress"`
	IsCurrent      bool      `json:"isCurrent"`
	CreatedAt      time.Time `json:"createdAt"`
}

// PasswordReset represents password reset request
type PasswordReset struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"-"`
	ExpiresAt time.Time `json:"expiresAt"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
