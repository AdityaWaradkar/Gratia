package user

import "time"

// DonorProfile represents a donor user profile in the system
type DonorProfile struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	Name      string     `json:"name"`
	Phone     *string    `json:"phone,omitempty"`
	Address   *string    `json:"address,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// NGOProfile represents an NGO user profile in the system
type NGOProfile struct {
	ID             string     `json:"id"`
	UserID         string     `json:"userId"`
	Organization   string     `json:"organization"`
	RegistrationNo string     `json:"registrationNo"`
	Verified       bool       `json:"verified"`
	VerifiedBy     *string    `json:"verifiedBy,omitempty"`
	VerifiedAt     *time.Time `json:"verifiedAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}