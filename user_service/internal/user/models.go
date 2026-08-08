package user

import "time"

// DonorProfile represents a food donor's profile linked to an authenticated account
type DonorProfile struct {
    ID        string    `json:"id"`
    UserID    string    `json:"userId"`
    Name      string    `json:"name"`
    Phone     *string   `json:"phone,omitempty"`
    Address   *string   `json:"address,omitempty"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// NGOProfile represents a non-governmental organization profile requiring administrative verification
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