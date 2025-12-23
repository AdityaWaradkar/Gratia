package food

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Service contains business logic for food_service
type Service struct {
	repo       Repository
	userClient UserClient
}

// NewService creates food service
func NewService(repo Repository, userClient UserClient) *Service {
	return &Service{
		repo:       repo,
		userClient: userClient,
	}
}

/* ===================== CREATE ===================== */

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

	// Validate donor via User Service
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

	if expiryTime.Before(time.Now()) {
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
		Status:      FoodStatusOpen,
	}

	if err := s.repo.CreateFoodListing(ctx, listing); err != nil {
		return nil, err
	}

	return listing, nil
}

/* ===================== READ ===================== */

// GetFoodListing returns a single listing
func (s *Service) GetFoodListing(
	ctx context.Context,
	id string,
) (*FoodListing, error) {

	if id == "" {
		return nil, errors.New("listing id required")
	}

	return s.repo.GetFoodListingByID(ctx, id)
}

// ListOpenFoodListings returns only OPEN listings
func (s *Service) ListOpenFoodListings(ctx context.Context) ([]*FoodListing, error) {
	return s.repo.ListFoodListings(ctx, FoodStatusOpen)
}

/* ===================== UPDATE ===================== */

// UpdateFoodListing allows donor to update own listing
func (s *Service) UpdateFoodListing(
	ctx context.Context,
	userID string,
	listing *FoodListing,
) error {

	if userID == "" {
		return errors.New("unauthorized")
	}

	existing, err := s.repo.GetFoodListingByID(ctx, listing.ID)
	if err != nil {
		return errors.New("food listing not found")
	}

	if existing.DonorUserID != userID {
		return errors.New("forbidden")
	}

	if existing.Status != FoodStatusOpen {
		return errors.New("only open listings can be updated")
	}

	if listing.ExpiryTime.Before(time.Now()) {
		return errors.New("expiry time must be in the future")
	}

	return s.repo.UpdateFoodListing(ctx, listing)
}

/* ===================== STATUS ===================== */

// MarkListingExpired is used by cron / background jobs
func (s *Service) MarkListingExpired(ctx context.Context, id string) error {
	return s.repo.UpdateFoodStatus(ctx, id, FoodStatusExpired)
}
