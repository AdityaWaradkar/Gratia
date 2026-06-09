package food

import "time"

// FoodListing represents a food donation listing
type FoodListing struct {
	ID string `json:"id"`

	DonorUserID string `json:"donorUserId"`

	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`

	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`

	ExpiryTime time.Time `json:"expiryTime"`

	Location string `json:"location"`

	ImageURL *string `json:"imageUrl,omitempty"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

const (
	FoodStatusAvailable = "AVAILABLE"
	FoodStatusClaimed   = "CLAIMED"
	FoodStatusExpired   = "EXPIRED"
	FoodStatusCancelled = "CANCELLED"
)