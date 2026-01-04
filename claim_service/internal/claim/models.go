package claim

import "time"

type ClaimStatus string

const (
	ClaimStatusRequested ClaimStatus = "REQUESTED"
	ClaimStatusApproved  ClaimStatus = "APPROVED"
	ClaimStatusRejected  ClaimStatus = "REJECTED"
	ClaimStatusCancelled ClaimStatus = "CANCELLED"
	ClaimStatusPickedUp  ClaimStatus = "PICKED_UP"
	ClaimStatusDelivered ClaimStatus = "DELIVERED"
)

type Claim struct {
	ID string `json:"id" db:"id"`

	FoodListingID string `json:"foodListingId" db:"food_listing_id"`

	NGOUserID   string `json:"ngoUserId" db:"ngo_user_id"`
	DonorUserID string `json:"donorUserId" db:"donor_user_id"`

	Status ClaimStatus `json:"status" db:"status"`

	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

func (c Claim) IsActive() bool {
	switch c.Status {
	case ClaimStatusRequested, ClaimStatusApproved, ClaimStatusPickedUp:
		return true
	default:
		return false
	}
}

func (c Claim) IsTerminal() bool {
	switch c.Status {
	case ClaimStatusRejected, ClaimStatusCancelled, ClaimStatusDelivered:
		return true
	default:
		return false
	}
}

func (c Claim) CanTransitionTo(next ClaimStatus) bool {
	switch c.Status {
	case ClaimStatusRequested:
		return next == ClaimStatusApproved ||
			next == ClaimStatusRejected ||
			next == ClaimStatusCancelled

	case ClaimStatusApproved:
		return next == ClaimStatusPickedUp ||
			next == ClaimStatusCancelled

	case ClaimStatusPickedUp:
		return next == ClaimStatusDelivered

	default:
		return false
	}
}

type ActorRole string

const (
	ActorNGO   ActorRole = "NGO"
	ActorDonor ActorRole = "DONOR"
)

func (c Claim) CanBeModifiedBy(actor ActorRole, userID string) bool {
	switch actor {
	case ActorNGO:
		return c.NGOUserID == userID
	case ActorDonor:
		return c.DonorUserID == userID
	default:
		return false
	}
}
