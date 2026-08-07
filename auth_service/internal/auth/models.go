package auth

import "time"

// User represents an authenticated account mapping directly to the database schema
type User struct {
    ID            string    `json:"id"`
    Email         string    `json:"email"`
    PasswordHash  string    `json:"-"`
    Role          Role      `json:"role"`
    EmailVerified bool      `json:"emailVerified"`
    CreatedAt     time.Time `json:"createdAt"`
    UpdatedAt     time.Time `json:"updatedAt"`
}

// RefreshToken tracks long-lived session credentials for token rotation
type RefreshToken struct {
    ID        string    `json:"id"`
    UserID    string    `json:"userId"`
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expiresAt"`
    Revoked   bool      `json:"revoked"`
    CreatedAt time.Time `json:"createdAt"`
}

// Session records device and network details for active user logins
type Session struct {
    ID             string    `json:"id"`
    UserID         string    `json:"userId"`
    RefreshTokenID string    `json:"refreshTokenId"`
    UserAgent      string    `json:"userAgent"`
    IPAddress      string    `json:"ipAddress"`
    IsCurrent      bool      `json:"isCurrent"`
    CreatedAt      time.Time `json:"createdAt"`
}

// PasswordReset handles secure token lifecycles for account recovery
type PasswordReset struct {
    ID        string    `json:"id"`
    UserID    string    `json:"userId"`
    Token     string    `json:"-"`
    ExpiresAt time.Time `json:"expiresAt"`
    Used      bool      `json:"used"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}