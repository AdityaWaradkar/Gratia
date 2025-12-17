package user

import "time"

/*
UserProfile represents profile data for any user
*/
type UserProfile struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	Role      string     `json:"role"`
	Name      string     `json:"name"`
	Phone     *string    `json:"phone,omitempty"`
	Address   *string    `json:"address,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
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
