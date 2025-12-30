package claim

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

/*
Domain errors
*/

var (
	ErrUnauthorized      = errors.New("unauthorized action")
	ErrInvalidState      = errors.New("invalid claim state transition")
	ErrActiveClaimExists = errors.New("active claim already exists")
	ErrNGONotVerified    = errors.New("ngo is not verified")
	ErrFoodNotOpen       = errors.New("food listing is not open")
)

/*
External service contracts
*/

type UserClient interface {
	IsNGOVerified(ctx context.Context, userID string) (bool, error)
}

type FoodClient interface {
	GetFoodForClaim(ctx context.Context, foodListingID string) (donorUserID string, status string, err error)
}

/*
Service
*/

type Service struct {
	db         *sqlx.DB
	repo       ClaimRepository
	userClient UserClient
	foodClient FoodClient
}

func NewService(
	db *sqlx.DB,
	repo ClaimRepository,
	userClient UserClient,
	foodClient FoodClient,
) *Service {
	return &Service{
		db:         db,
		repo:       repo,
		userClient: userClient,
		foodClient: foodClient,
	}
}

/*
Create claim (NGO)
*/

func (s *Service) CreateClaim(
	ctx context.Context,
	foodListingID string,
	ngoUserID string,
) (*Claim, error) {

	verified, err := s.userClient.IsNGOVerified(ctx, ngoUserID)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, ErrNGONotVerified
	}

	donorUserID, foodStatus, err := s.foodClient.GetFoodForClaim(ctx, foodListingID)
	if err != nil {
		return nil, err
	}
	if foodStatus != "OPEN" {
		return nil, ErrFoodNotOpen
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := s.repo.GetActiveByFoodID(ctx, foodListingID); err == nil {
		return nil, ErrActiveClaimExists
	}

	now := time.Now().UTC()

	claim := &Claim{
		ID:            uuid.NewString(),
		FoodListingID: foodListingID,
		NGOUserID:     ngoUserID,
		DonorUserID:   donorUserID,
		Status:        ClaimStatusRequested,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(ctx, tx, claim); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return claim, nil
}

/*
Donor actions
*/

func (s *Service) ApproveClaim(ctx context.Context, claimID, donorUserID string) error {
	return s.updateStatus(ctx, claimID, donorUserID, ActorDonor, ClaimStatusApproved)
}

func (s *Service) RejectClaim(ctx context.Context, claimID, donorUserID string) error {
	return s.updateStatus(ctx, claimID, donorUserID, ActorDonor, ClaimStatusRejected)
}

/*
NGO actions
*/

func (s *Service) CancelByNGO(ctx context.Context, claimID, ngoUserID string) error {
	return s.updateStatus(ctx, claimID, ngoUserID, ActorNGO, ClaimStatusCancelled)
}

func (s *Service) MarkPickedUp(ctx context.Context, claimID, ngoUserID string) error {
	return s.updateStatus(ctx, claimID, ngoUserID, ActorNGO, ClaimStatusPickedUp)
}

func (s *Service) MarkDelivered(ctx context.Context, claimID, ngoUserID string) error {
	return s.updateStatus(ctx, claimID, ngoUserID, ActorNGO, ClaimStatusDelivered)
}

/*
Shared state transition logic
*/

func (s *Service) updateStatus(
	ctx context.Context,
	claimID string,
	userID string,
	actor ActorRole,
	nextStatus ClaimStatus,
) error {

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

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

	if err := s.repo.UpdateStatus(ctx, tx, claimID, nextStatus); err != nil {
		return err
	}

	return tx.Commit()
}
