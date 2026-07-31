package food

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Service contains business logic for food operations
type Service struct {
	repo       Repository
	userClient UserClient
}

// NewService creates a new food service instance
func NewService(
	repo Repository,
	userClient UserClient,
) *Service {
	return &Service{
		repo:       repo,
		userClient: userClient,
	}
}

// CreateFoodListing allows only DONOR users to create food listings
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

	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	isDonor, err := s.userClient.IsDonor(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !isDonor {
		return nil, errors.New("only donors can create food listings")
	}

	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title is required")
	}

	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	if strings.TrimSpace(unit) == "" {
		return nil, errors.New("unit is required")
	}

	if expiryTime.Before(time.Now().UTC()) {
		return nil, errors.New("expiry time must be in the future")
	}

	if strings.TrimSpace(location) == "" {
		return nil, errors.New("location is required")
	}

	listing := &FoodListing{
		DonorUserID: userID,
		Title:       title,
		Description: description,
		Quantity:    quantity,
		Unit:        unit,
		ExpiryTime:  expiryTime,
		Location:    location,
		ImageURL:    imageURL,
		Status:      FoodStatusAvailable,
	}

	if err := s.repo.CreateFoodListing(ctx, listing); err != nil {
		return nil, err
	}

	return listing, nil
}

// GetFoodListing retrieves a single food listing by ID
func (s *Service) GetFoodListing(
	ctx context.Context,
	id string,
) (*FoodListing, error) {

	if id == "" {
		return nil, errors.New("listing id required")
	}

	return s.repo.GetFoodListingByID(ctx, id)
}

// ListAvailableFoodListings retrieves all available food listings
func (s *Service) ListAvailableFoodListings(
	ctx context.Context,
) ([]*FoodListing, error) {

	return s.repo.ListFoodListings(
		ctx,
		FoodStatusAvailable,
	)
}

// UpdateFoodListing allows a donor to update their own listing
func (s *Service) UpdateFoodListing(
	ctx context.Context,
	userID string,
	listing *FoodListing,
) error {

	if userID == "" {
		return errors.New("unauthorized")
	}

	existing, err := s.repo.GetFoodListingByID(
		ctx,
		listing.ID,
	)
	if err != nil {
		return errors.New("food listing not found")
	}

	if existing.DonorUserID != userID {
		return errors.New("forbidden")
	}

	if existing.Status != FoodStatusAvailable {
		return errors.New("only available listings can be updated")
	}

	if strings.TrimSpace(listing.Title) == "" {
		return errors.New("title is required")
	}

	if listing.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if strings.TrimSpace(listing.Unit) == "" {
		return errors.New("unit is required")
	}

	if strings.TrimSpace(listing.Location) == "" {
		return errors.New("location is required")
	}

	if listing.ExpiryTime.Before(time.Now().UTC()) {
		return errors.New("expiry time must be in the future")
	}

	return s.repo.UpdateFoodListing(
		ctx,
		listing,
	)
}

// CancelFoodListing allows a donor to cancel their own listing
func (s *Service) CancelFoodListing(
	ctx context.Context,
	userID string,
	listingID string,
) error {

	if userID == "" {
		return errors.New("unauthorized")
	}

	listing, err := s.repo.GetFoodListingByID(
		ctx,
		listingID,
	)
	if err != nil {
		return errors.New("food listing not found")
	}

	if listing.DonorUserID != userID {
		return errors.New("forbidden")
	}

	if listing.Status != FoodStatusAvailable {
		return errors.New("only available listings can be cancelled")
	}

	return s.repo.UpdateFoodStatus(
		ctx,
		listingID,
		FoodStatusCancelled,
	)
}

// IsClaimable validates if a food listing can be claimed
func (s *Service) IsClaimable(
	ctx context.Context,
	id string,
) (bool, string, error) {

	listing, err := s.repo.GetFoodListingByID(
		ctx,
		id,
	)
	if err != nil {
		return false, "NOT_FOUND", nil
	}

	if listing.Status != FoodStatusAvailable {
		return false, listing.Status, nil
	}

	if listing.ExpiryTime.Before(time.Now().UTC()) {
		return false, FoodStatusExpired, nil
	}

	return true, "", nil
}

// MarkClaimed marks a food listing as claimed
func (s *Service) MarkClaimed(
	ctx context.Context,
	id string,
) error {

	listing, err := s.repo.GetFoodListingByID(
		ctx,
		id,
	)
	if err != nil {
		return err
	}

	if listing.Status != FoodStatusAvailable {
		return errors.New("food listing is not available")
	}

	return s.repo.UpdateFoodStatus(
		ctx,
		id,
		FoodStatusClaimed,
	)
}

// ExpireListings expires food listings that have passed their expiry time
func (s *Service) ExpireListings(
	ctx context.Context,
) error {

	return s.repo.ExpireFoodListings(ctx)
}