package user

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Service contains business logic for user_service
type Service struct {
	repo Repository
}

// NewService creates a new service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

/* ===================== DONOR PROFILE ===================== */

// CreateDonorProfile creates donor profile for authenticated user
func (s *Service) CreateDonorProfile(
	ctx context.Context,
	userID string,
	role string,
	name string,
	phone *string,
	address *string,
) (*DonorProfile, error) {

	if userID == "" || role == "" {
		return nil, errors.New("unauthorized")
	}

	if role != "USER" {
		return nil, errors.New("only donor users can create donor profile")
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name is required")
	}

	donor := &DonorProfile{
		UserID:     userID,
		Name:       name,
		Phone:      phone,
		Address:    address,
		IsVerified: false,
	}

	if err := s.repo.CreateDonorProfile(ctx, donor); err != nil {
		return nil, err
	}

	return donor, nil
}

// GetMyDonorProfile returns donor profile of authenticated user
func (s *Service) GetMyDonorProfile(
	ctx context.Context,
	userID string,
) (*DonorProfile, error) {

	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	return s.repo.GetDonorProfileByUserID(ctx, userID)
}

// UpdateMyDonorProfile updates donor profile
func (s *Service) UpdateMyDonorProfile(
	ctx context.Context,
	userID string,
	name string,
	phone *string,
	address *string,
) error {

	if userID == "" {
		return errors.New("unauthorized")
	}

	profile, err := s.repo.GetDonorProfileByUserID(ctx, userID)
	if err != nil {
		return errors.New("donor profile not found")
	}

	if strings.TrimSpace(name) != "" {
		profile.Name = name
	}

	profile.Phone = phone
	profile.Address = address
	profile.UpdatedAt = time.Now()

	return s.repo.UpdateDonorProfile(ctx, profile)
}

/* ===================== NGO PROFILE ===================== */

// CreateNGOProfile creates NGO profile for NGO user
func (s *Service) CreateNGOProfile(
	ctx context.Context,
	userID string,
	role string,
	organization string,
	registrationNo string,
) (*NGOProfile, error) {

	if userID == "" || role == "" {
		return nil, errors.New("unauthorized")
	}

	if role != "NGO" {
		return nil, errors.New("only NGO users can create ngo profile")
	}

	if strings.TrimSpace(organization) == "" || strings.TrimSpace(registrationNo) == "" {
		return nil, errors.New("invalid ngo details")
	}

	ngo := &NGOProfile{
		UserID:         userID,
		Organization:   organization,
		RegistrationNo: registrationNo,
		Verified:       false,
	}

	if err := s.repo.CreateNGOProfile(ctx, ngo); err != nil {
		return nil, err
	}

	return ngo, nil
}

// GetMyNGOProfile returns NGO profile of authenticated user
func (s *Service) GetMyNGOProfile(
	ctx context.Context,
	userID string,
) (*NGOProfile, error) {

	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	return s.repo.GetNGOProfileByUserID(ctx, userID)
}

/* ===================== ADMIN ===================== */

// VerifyNGO allows admin to verify NGO
func (s *Service) VerifyNGO(
	ctx context.Context,
	adminUserID string,
	adminRole string,
	targetUserID string,
) error {

	if adminUserID == "" || adminRole == "" {
		return errors.New("unauthorized")
	}

	if adminRole != "ADMIN" {
		return errors.New("forbidden")
	}

	if targetUserID == "" {
		return errors.New("target user id required")
	}

	return s.repo.UpdateNGOVerification(
		ctx,
		targetUserID,
		true,
		&adminUserID,
	)
}
