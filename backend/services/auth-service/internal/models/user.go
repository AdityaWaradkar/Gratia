package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID  `db:"user_id" json:"user_id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	FullName     string     `db:"full_name" json:"full_name,omitempty"`
	PhoneNumber  string     `db:"phone_number" json:"phone_number,omitempty"`
	Role         string     `db:"role" json:"role"`
	LocationLat  *float64   `db:"location_lat" json:"location_lat,omitempty"`
	LocationLng  *float64   `db:"location_lng" json:"location_lng,omitempty"`
	IsActive     *bool      `db:"is_active" json:"is_active,omitempty"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
