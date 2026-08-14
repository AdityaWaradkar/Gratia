package claim

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Domain errors for claim operations
var (
	ErrUnauthorized      = errors.New("unauthorized action")
	ErrInvalidState      = errors.New("invalid claim state transition")
	ErrActiveClaimExists = errors.New("active claim already exists")
	ErrNGONotVerified    = errors.New("ngo is not verified")
	ErrFoodNotOpen       = errors.New("food listing is not open")
	ErrSelfClaim         = errors.New("ngo cannot claim its own food listing")
)

// Service handles business logic for claim operations
type Service struct {
	repo       Repository
	userClient *UserClient
	foodClient *FoodClient
}

// NewService creates a new claim service instance
func NewService(repo Repository, userClient *UserClient, foodClient *FoodClient) *Service {
	return &Service{
		repo:       repo,
		userClient: userClient,
		foodClient: foodClient,
	}
}

// CreateClaim creates a new claim for a food listing
func (s *Service) CreateClaim(ctx context.Context, foodListingID string, ngoUserID string) (*Claim, error) {
	// Validate NGO is verified
	verified, err := s.userClient.IsVerifiedNGO(ctx, ngoUserID)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, ErrNGONotVerified
	}

	// Get food listing details
	donorUserID, foodStatus, err := s.foodClient.GetFoodForClaim(ctx, foodListingID)
	if err != nil {
		return nil, err
	}

	// Validate food listing is available
	if foodStatus != "AVAILABLE" {
		return nil, ErrFoodNotOpen
	}

	// Prevent self-claiming
	if donorUserID == ngoUserID {
		return nil, ErrSelfClaim
	}

	// Check for existing active claim
	existing, err := s.repo.GetActiveByFoodID(ctx, foodListingID)
	if err != nil && !errors.Is(err, ErrClaimNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrActiveClaimExists
	}

	now := time.Now()

	claim := &Claim{
		ID:            uuid.NewString(),
		FoodListingID: foodListingID,
		NGOUserID:     ngoUserID,
		DonorUserID:   donorUserID,
		Status:        ClaimStatusCreated,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(ctx, claim); err != nil {
		return nil, err
	}

	return claim, nil
}

// ApproveClaim approves a claim by the donor
func (s *Service) ApproveClaim(ctx context.Context, claimID string, donorUserID string) error {
	return s.updateStatus(ctx, claimID, donorUserID, ActorDonor, ClaimStatusAccepted)
}

// RejectClaim rejects a claim by the donor
func (s *Service) RejectClaim(ctx context.Context, claimID string, donorUserID string) error {
	return s.updateStatus(ctx, claimID, donorUserID, ActorDonor, ClaimStatusRejected)
}

// CancelByNGO cancels a claim by the NGO
func (s *Service) CancelByNGO(ctx context.Context, claimID string, ngoUserID string) error {
	return s.updateStatus(ctx, claimID, ngoUserID, ActorNGO, ClaimStatusCancelled)
}

// MarkPickedUp marks a claim as picked up by the NGO
func (s *Service) MarkPickedUp(ctx context.Context, claimID string, ngoUserID string) error {
	return s.updateStatus(ctx, claimID, ngoUserID, ActorNGO, ClaimStatusPickedUp)
}

// MarkDelivered marks a claim as delivered by the NGO
func (s *Service) MarkDelivered(ctx context.Context, claimID string, ngoUserID string) error {
	return s.updateStatus(ctx, claimID, ngoUserID, ActorNGO, ClaimStatusDelivered)
}

// updateStatus handles the shared state transition logic for claims
func (s *Service) updateStatus(ctx context.Context, claimID string, userID string, actor ActorRole, nextStatus ClaimStatus) error {
	claim, err := s.repo.GetByID(ctx, claimID)
	if err != nil {
		return err
	}

	if !claim.CanBeModifiedBy(actor, userID) {
		return ErrUnauthorized
	}

	if !claim.CanTransitionTo(nextStatus) {
		return ErrInvalidState
	}

	if err := s.repo.UpdateStatus(ctx, claimID, nextStatus); err != nil {
		return err
	}

	// If claim is accepted, mark food as claimed
	if nextStatus == ClaimStatusAccepted {
		if err := s.foodClient.MarkFoodClaimed(ctx, claim.FoodListingID); err != nil {
			return err
		}
	}

	return nil
}

// GetClaimByID retrieves a claim by ID
func (s *Service) GetClaimByID(ctx context.Context, claimID string) (*Claim, error) {
	return s.repo.GetByID(ctx, claimID)
}

// GetClaimsByNGO retrieves all claims for an NGO
func (s *Service) GetClaimsByNGO(ctx context.Context, ngoUserID string) ([]Claim, error) {
	return s.repo.GetByNGOUserID(ctx, ngoUserID)
}

// GetClaimsByDonor retrieves all claims for a donor
func (s *Service) GetClaimsByDonor(ctx context.Context, donorUserID string) ([]Claim, error) {
	return s.repo.GetByDonorUserID(ctx, donorUserID)
}

// GetClaimsByFood retrieves all claims for a food listing
func (s *Service) GetClaimsByFood(ctx context.Context, foodListingID string) ([]Claim, error) {
	return s.repo.GetByFoodID(ctx, foodListingID)
}