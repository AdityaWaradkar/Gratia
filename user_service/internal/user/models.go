package user

import "time"

/*
UserProfile represents profile data for any user
*/
type UserProfile struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"` // from Auth Service
	Role      string    `json:"role"`    // Donor | NGO | Admin

	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/*
NGOProfile holds NGO-specific data
*/
type NGOProfile struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`

	Organization   string    `json:"organization"`
	RegistrationNo string    `json:"registration_no"`

	Verified       bool      `json:"verified"`
	VerifiedBy     *string   `json:"verified_by,omitempty"`
	VerifiedAt     *time.Time`json:"verified_at,omitempty"`

	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
