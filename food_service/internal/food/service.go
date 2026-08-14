package food

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Domain-level sentinel errors for uniform HTTP status mapping
var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrListingNotFound    = errors.New("food listing not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrStatusNotAvailable = errors.New("only available listings can be modified")
)

// Service contains the core business logic for food operations and cross-service validation
type Service struct {
	repo       Repository
	userClient UserClient
}

// NewService creates and returns a new food service instance
func NewService(repo Repository, userClient UserClient) *Service {
	return &Service{
		repo:       repo,
		userClient: userClient,
	}
}

// CreateFoodListing validates the user's donor status via the User Service before creating a listing
func (s *Service) CreateFoodListing(
	ctx context.Context,
	userID string,
	title string,
	description *string,
	quantity int,
	unit string,
	expiryTime time.Time,
	location string,
	imageURL *string,
) (*FoodListing, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}

	// Synchronous cross-service check to ensure the user has a registered Donor profile
	isDonor, err := s.userClient.IsDonor(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify donor status: %w", err)
	}
	if !isDonor {
		return nil, ErrForbidden
	}

	// Input validation
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than zero", ErrInvalidInput)
	}
	if strings.TrimSpace(unit) == "" {
		return nil, fmt.Errorf("%w: unit is required", ErrInvalidInput)
	}
	if strings.TrimSpace(location) == "" {
		return nil, fmt.Errorf("%w: location is required", ErrInvalidInput)
	}
	if expiryTime.Before(time.Now()) {
		return nil, fmt.Errorf("%w: expiry time must be in the future", ErrInvalidInput)
	}

	listing := &FoodListing{
		DonorUserID: userID,
		Title:       strings.TrimSpace(title),
		Description: description,
		Quantity:    quantity,
		Unit:        strings.TrimSpace(unit),
		ExpiryTime:  expiryTime,
		Location:    strings.TrimSpace(location),
		ImageURL:    imageURL,
		Status:      FoodStatusAvailable,
	}

	if err := s.repo.CreateFoodListing(ctx, listing); err != nil {
		return nil, err
	}

	return listing, nil
}

// GetFoodListing retrieves a single food listing by ID
func (s *Service) GetFoodListing(ctx context.Context, id string) (*FoodListing, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: listing id required", ErrInvalidInput)
	}

	listing, err := s.repo.GetFoodListingByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if listing == nil {
		return nil, ErrListingNotFound
	}

	return listing, nil
}

// ListAvailableFoodListings retrieves all food listings currently marked as AVAILABLE
func (s *Service) ListAvailableFoodListings(ctx context.Context) ([]*FoodListing, error) {
	return s.repo.ListFoodListings(ctx, FoodStatusAvailable)
}

// UpdateFoodListing allows a donor to update their own active listing
func (s *Service) UpdateFoodListing(ctx context.Context, userID string, listing *FoodListing) error {
	if strings.TrimSpace(userID) == "" {
		return ErrUnauthorized
	}

	existing, err := s.repo.GetFoodListingByID(ctx, listing.ID)
	if err != nil {
		return ErrListingNotFound
	}
	if existing == nil {
		return ErrListingNotFound
	}

	if existing.DonorUserID != userID {
		return ErrForbidden
	}

	if existing.Status != FoodStatusAvailable {
		return ErrStatusNotAvailable
	}

	// Re-validate inputs before update
	if strings.TrimSpace(listing.Title) == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if listing.Quantity <= 0 {
		return fmt.Errorf("%w: quantity must be greater than zero", ErrInvalidInput)
	}
	if strings.TrimSpace(listing.Unit) == "" {
		return fmt.Errorf("%w: unit is required", ErrInvalidInput)
	}
	if strings.TrimSpace(listing.Location) == "" {
		return fmt.Errorf("%w: location is required", ErrInvalidInput)
	}
	if listing.ExpiryTime.Before(time.Now()) {
		return fmt.Errorf("%w: expiry time must be in the future", ErrInvalidInput)
	}

	return s.repo.UpdateFoodListing(ctx, listing)
}

// CancelFoodListing allows a donor to safely cancel their own available listing
func (s *Service) CancelFoodListing(ctx context.Context, userID string, listingID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrUnauthorized
	}

	listing, err := s.repo.GetFoodListingByID(ctx, listingID)
	if err != nil {
		return ErrListingNotFound
	}
	if listing == nil {
		return ErrListingNotFound
	}

	if listing.DonorUserID != userID {
		return ErrForbidden
	}

	if listing.Status != FoodStatusAvailable {
		return ErrStatusNotAvailable
	}

	return s.repo.UpdateFoodStatus(ctx, listingID, FoodStatusCancelled)
}

// IsClaimable evaluates if a food listing is currently open for NGOs to claim
func (s *Service) IsClaimable(ctx context.Context, id string) (bool, string, error) {
	listing, err := s.repo.GetFoodListingByID(ctx, id)
	if err != nil {
		return false, "NOT_FOUND", nil
	}
	if listing == nil {
		return false, "NOT_FOUND", nil
	}

	if listing.Status != FoodStatusAvailable {
		return false, listing.Status, nil
	}

	if listing.ExpiryTime.Before(time.Now()) {
		return false, FoodStatusExpired, nil
	}

	return true, "", nil
}

// MarkClaimed locks a food listing by changing its status to CLAIMED
func (s *Service) MarkClaimed(ctx context.Context, id string) error {
	listing, err := s.repo.GetFoodListingByID(ctx, id)
	if err != nil {
		return ErrListingNotFound
	}
	if listing == nil {
		return ErrListingNotFound
	}

	if listing.Status != FoodStatusAvailable {
		return ErrStatusNotAvailable
	}

	return s.repo.UpdateFoodStatus(ctx, id, FoodStatusClaimed)
}

// ExpireListings triggers the repository to sweep and expire any outdated listings
func (s *Service) ExpireListings(ctx context.Context) error {
	return s.repo.ExpireFoodListings(ctx)
}