package user

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Service handles business logic for user operations
type Service struct {
	repo Repository
}

// NewService creates a new user service instance
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetDonorProfileByUserIDInternal retrieves a donor profile by user ID for internal use
func (s *Service) GetDonorProfileByUserIDInternal(
	ctx context.Context,
	userID string,
) (*DonorProfile, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidUserID
	}

	return s.repo.GetDonorProfileByUserID(ctx, userID)
}

// GetNGOProfileByUserIDInternal retrieves an NGO profile by user ID for internal use
func (s *Service) GetNGOProfileByUserIDInternal(
	ctx context.Context,
	userID string,
) (*NGOProfile, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidUserID
	}

	return s.repo.GetNGOProfileByUserID(ctx, userID)
}

// CreateDonorProfile creates a new donor profile for a user
func (s *Service) CreateDonorProfile(
	ctx context.Context,
	userID string,
	role string,
	name string,
	phone *string,
	address *string,
) (*DonorProfile, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}

	if role != RoleUser {
		return nil, ErrForbidden
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name is required")
	}

	existing, err := s.repo.GetDonorProfileByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, ErrDonorProfileExists
	}

	donor := &DonorProfile{
		UserID:  userID,
		Name:    strings.TrimSpace(name),
		Phone:   phone,
		Address: address,
	}

	if err := s.repo.CreateDonorProfile(ctx, donor); err != nil {
		return nil, err
	}

	return donor, nil
}

// GetMyDonorProfile retrieves the donor profile of the authenticated user
func (s *Service) GetMyDonorProfile(
	ctx context.Context,
	userID string,
) (*DonorProfile, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}

	return s.repo.GetDonorProfileByUserID(ctx, userID)
}

// UpdateMyDonorProfile updates the donor profile of the authenticated user
func (s *Service) UpdateMyDonorProfile(
	ctx context.Context,
	userID string,
	name string,
	phone *string,
	address *string,
) error {

	if strings.TrimSpace(userID) == "" {
		return ErrUnauthorized
	}

	profile, err := s.repo.GetDonorProfileByUserID(ctx, userID)
	if err != nil {
		return ErrDonorProfileNotFound
	}

	if strings.TrimSpace(name) != "" {
		profile.Name = strings.TrimSpace(name)
	}

	profile.Phone = phone
	profile.Address = address
	profile.UpdatedAt = time.Now()

	return s.repo.UpdateDonorProfile(ctx, profile)
}

// CreateNGOProfile creates a new NGO profile for a user
func (s *Service) CreateNGOProfile(
	ctx context.Context,
	userID string,
	role string,
	organization string,
	registrationNo string,
) (*NGOProfile, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}

	if role != RoleUser {
		return nil, ErrForbidden
	}

	if strings.TrimSpace(organization) == "" ||
		strings.TrimSpace(registrationNo) == "" {
		return nil, ErrInvalidNGODetails
	}

	existing, err := s.repo.GetNGOProfileByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, ErrNGOProfileExists
	}

	ngo := &NGOProfile{
		UserID:         userID,
		Organization:   strings.TrimSpace(organization),
		RegistrationNo: strings.TrimSpace(registrationNo),
		Verified:       false,
	}

	if err := s.repo.CreateNGOProfile(ctx, ngo); err != nil {
		return nil, err
	}

	return ngo, nil
}

// GetMyNGOProfile retrieves the NGO profile of the authenticated user
func (s *Service) GetMyNGOProfile(
	ctx context.Context,
	userID string,
) (*NGOProfile, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}

	return s.repo.GetNGOProfileByUserID(ctx, userID)
}

// VerifyNGO verifies an NGO profile (admin only)
func (s *Service) VerifyNGO(
	ctx context.Context,
	adminUserID string,
	adminRole string,
	targetUserID string,
) error {

	if strings.TrimSpace(adminUserID) == "" {
		return ErrUnauthorized
	}

	if adminRole != RoleAdmin {
		return ErrForbidden
	}

	if strings.TrimSpace(targetUserID) == "" {
		return ErrTargetUserIDRequired
	}

	ngo, err := s.repo.GetNGOProfileByUserID(ctx, targetUserID)
	if err != nil {
		return ErrNGOProfileNotFound
	}

	if ngo.Verified {
		return ErrNGOAlreadyVerified
	}

	return s.repo.UpdateNGOVerification(
		ctx,
		targetUserID,
		true,
		&adminUserID,
	)
}