package user

import (
	"context"
	"errors"
	"strings"
)

// Service encapsulates the core business logic and rules for managing user profiles
type Service struct {
	repo Repository
}

// NewService initializes and returns a new user service instance with its repository dependency
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetDonorProfileByUserIDInternal safely retrieves a donor profile for internal service-to-service communication
func (s *Service) GetDonorProfileByUserIDInternal(ctx context.Context, userID string) (*DonorProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidUserID
	}
	return s.repo.GetDonorProfileByUserID(ctx, userID)
}

// GetNGOProfileByUserIDInternal safely retrieves an NGO profile for internal service-to-service communication
func (s *Service) GetNGOProfileByUserIDInternal(ctx context.Context, userID string) (*NGOProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidUserID
	}
	return s.repo.GetNGOProfileByUserID(ctx, userID)
}

// CreateDonorProfile registers a new donor profile while enforcing proper role boundaries and validation
func (s *Service) CreateDonorProfile(ctx context.Context, userID string, role string, name string, phone *string, address *string) (*DonorProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}

	if role != string(RoleDonor) && role != string(RoleUser) {
		return nil, ErrForbidden
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name is required")
	}

	existing, err := s.repo.GetDonorProfileByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, ErrDonorProfileExists
	}
	if err != nil && !errors.Is(err, ErrDonorProfileNotFound) {
		return nil, err
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

// GetMyDonorProfile retrieves the profile belonging to the currently authenticated donor
func (s *Service) GetMyDonorProfile(ctx context.Context, userID string) (*DonorProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}
	return s.repo.GetDonorProfileByUserID(ctx, userID)
}

// UpdateMyDonorProfile safely modifies fields on the authenticated user's active donor profile
func (s *Service) UpdateMyDonorProfile(ctx context.Context, userID string, name string, phone *string, address *string) error {
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

	return s.repo.UpdateDonorProfile(ctx, profile)
}

// CreateNGOProfile registers a new NGO organization profile with unverified initial status
func (s *Service) CreateNGOProfile(ctx context.Context, userID string, role string, organization string, registrationNo string) (*NGOProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}

	if role != string(RoleNGO) && role != string(RoleUser) {
		return nil, ErrForbidden
	}

	if strings.TrimSpace(organization) == "" || strings.TrimSpace(registrationNo) == "" {
		return nil, ErrInvalidNGODetails
	}

	existing, err := s.repo.GetNGOProfileByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, ErrNGOProfileExists
	}
	if err != nil && !errors.Is(err, ErrNGOProfileNotFound) {
		return nil, err
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

// GetMyNGOProfile retrieves the organization profile belonging to the authenticated NGO user
func (s *Service) GetMyNGOProfile(ctx context.Context, userID string) (*NGOProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}
	return s.repo.GetNGOProfileByUserID(ctx, userID)
}

// VerifyNGO executes administrative validation, marking an NGO profile as verified in the system
func (s *Service) VerifyNGO(ctx context.Context, adminUserID string, adminRole string, targetUserID string) error {
	if strings.TrimSpace(adminUserID) == "" {
		return ErrUnauthorized
	}

	if adminRole != string(RoleAdmin) {
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

	return s.repo.UpdateNGOVerification(ctx, targetUserID, true, &adminUserID)
}