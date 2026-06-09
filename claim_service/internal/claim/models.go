package claim

import "time"

/*
Claim Status
*/

type ClaimStatus string

const (
	ClaimStatusCreated   ClaimStatus = "CREATED"
	ClaimStatusAccepted  ClaimStatus = "ACCEPTED"
	ClaimStatusRejected  ClaimStatus = "REJECTED"
	ClaimStatusPickedUp  ClaimStatus = "PICKED_UP"
	ClaimStatusDelivered ClaimStatus = "DELIVERED"
	ClaimStatusCancelled ClaimStatus = "CANCELLED"
)

/*
Claim
*/

type Claim struct {
	ID string `json:"id" db:"id"`

	FoodListingID string `json:"foodListingId" db:"food_listing_id"`

	NGOUserID   string `json:"ngoUserId" db:"ngo_user_id"`
	DonorUserID string `json:"donorUserId" db:"donor_user_id"`

	Status ClaimStatus `json:"status" db:"status"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`

	AcceptedAt  *time.Time `json:"acceptedAt,omitempty" db:"accepted_at"`
	RejectedAt  *time.Time `json:"rejectedAt,omitempty" db:"rejected_at"`
	PickedUpAt  *time.Time `json:"pickedUpAt,omitempty" db:"picked_up_at"`
	DeliveredAt *time.Time `json:"deliveredAt,omitempty" db:"delivered_at"`
	CancelledAt *time.Time `json:"cancelledAt,omitempty" db:"cancelled_at"`
}

/*
Helpers
*/

// Active claims participate in the unique active-claim constraint.
func (c Claim) IsActive() bool {
	switch c.Status {
	case ClaimStatusCreated,
		ClaimStatusAccepted,
		ClaimStatusPickedUp:
		return true
	default:
		return false
	}
}

// Terminal states cannot transition further.
func (c Claim) IsTerminal() bool {
	switch c.Status {
	case ClaimStatusRejected,
		ClaimStatusCancelled,
		ClaimStatusDelivered:
		return true
	default:
		return false
	}
}

// Valid state transitions.
func (c Claim) CanTransitionTo(next ClaimStatus) bool {
	switch c.Status {

	case ClaimStatusCreated:
		return next == ClaimStatusAccepted ||
			next == ClaimStatusRejected ||
			next == ClaimStatusCancelled

	case ClaimStatusAccepted:
		return next == ClaimStatusPickedUp ||
			next == ClaimStatusCancelled

	case ClaimStatusPickedUp:
		return next == ClaimStatusDelivered

	default:
		return false
	}
}

/*
Actor Roles
*/

type ActorRole string

const (
	ActorNGO   ActorRole = "NGO"
	ActorDonor ActorRole = "DONOR"
)

/*
Authorization
*/

func (c Claim) CanBeModifiedBy(
	actor ActorRole,
	userID string,
) bool {

	switch actor {

	case ActorNGO:
		return c.NGOUserID == userID

	case ActorDonor:
		return c.DonorUserID == userID

	default:
		return false
	}
}