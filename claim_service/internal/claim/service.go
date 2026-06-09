package claim

import (
	"context"
	"errors"
	"strings"
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
	ErrSelfClaim         = errors.New("ngo cannot claim its own food listing")
)

/*
External service contracts
*/

type UserClient interface {
	IsNGOVerified(ctx context.Context, userID string) (bool, error)
}

type FoodClient interface {
	GetFoodForClaim(
		ctx context.Context,
		foodListingID string,
	) (donorUserID string, status string, err error)
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
Create Claim
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

	donorUserID, foodStatus, err := s.foodClient.GetFoodForClaim(
		ctx,
		foodListingID,
	)
	if err != nil {
		return nil, err
	}

	if foodStatus != "OPEN" {
		return nil, ErrFoodNotOpen
	}

	// NGO cannot claim its own donation.
	if donorUserID == ngoUserID {
		return nil, ErrSelfClaim
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now().UTC()

	claim := &Claim{
		ID:            uuid.NewString(),
		FoodListingID: foodListingID,
		NGOUserID:     ngoUserID,
		DonorUserID:   donorUserID,
		Status:        ClaimStatusCreated,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(ctx, tx, claim); err != nil {

		// PostgreSQL unique partial index violation
		if isUniqueViolation(err) {
			return nil, ErrActiveClaimExists
		}

		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return claim, nil
}

/*
Donor Actions
*/

func (s *Service) ApproveClaim(
	ctx context.Context,
	claimID string,
	donorUserID string,
) error {
	return s.updateStatus(
		ctx,
		claimID,
		donorUserID,
		ActorDonor,
		ClaimStatusAccepted,
	)
}

func (s *Service) RejectClaim(
	ctx context.Context,
	claimID string,
	donorUserID string,
) error {
	return s.updateStatus(
		ctx,
		claimID,
		donorUserID,
		ActorDonor,
		ClaimStatusRejected,
	)
}

/*
NGO Actions
*/

func (s *Service) CancelByNGO(
	ctx context.Context,
	claimID string,
	ngoUserID string,
) error {
	return s.updateStatus(
		ctx,
		claimID,
		ngoUserID,
		ActorNGO,
		ClaimStatusCancelled,
	)
}

func (s *Service) MarkPickedUp(
	ctx context.Context,
	claimID string,
	ngoUserID string,
) error {
	return s.updateStatus(
		ctx,
		claimID,
		ngoUserID,
		ActorNGO,
		ClaimStatusPickedUp,
	)
}

func (s *Service) MarkDelivered(
	ctx context.Context,
	claimID string,
	ngoUserID string,
) error {
	return s.updateStatus(
		ctx,
		claimID,
		ngoUserID,
		ActorNGO,
		ClaimStatusDelivered,
	)
}

/*
Shared State Transition Logic
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

	defer func() {
		_ = tx.Rollback()
	}()

	claim, err := s.repo.GetByIDForUpdate(
		ctx,
		tx,
		claimID,
	)
	if err != nil {
		return err
	}

	if !claim.CanBeModifiedBy(actor, userID) {
		return ErrUnauthorized
	}

	if !claim.CanTransitionTo(nextStatus) {
		return ErrInvalidState
	}

	if err := s.repo.UpdateStatus(
		ctx,
		tx,
		claimID,
		nextStatus,
	); err != nil {
		return err
	}

	return tx.Commit()
}

/*
PostgreSQL Helpers
*/

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(
		strings.ToLower(err.Error()),
		"duplicate key value",
	)
}