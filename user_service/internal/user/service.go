package user

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

/* ===================== INTERNAL ===================== */

// GetDonorProfileByUserIDInternal is used by other services
func (s *Service) GetDonorProfileByUserIDInternal(
	ctx context.Context,
	userID string,
) (*DonorProfile, error) {

	if userID == "" {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetDonorProfileByUserID(ctx, userID)
}

// GetNGOProfileByUserIDInternal is used by other services
func (s *Service) GetNGOProfileByUserIDInternal(
	ctx context.Context,
	userID string,
) (*NGOProfile, error) {

	if userID == "" {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetNGOProfileByUserID(ctx, userID)
}

/* ===================== DONOR PROFILE ===================== */

func (s *Service) CreateDonorProfile(
	ctx context.Context,
	userID string,
	role string,
	name string,
	phone *string,
	address *string,
) (*DonorProfile, error) {

	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	if role != "USER" {
		return nil, errors.New("only authenticated users can create donor profile")
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name is required")
	}

	donor := &DonorProfile{
		UserID:  userID,
		Name:    name,
		Phone:   phone,
		Address: address,
	}

	if err := s.repo.CreateDonorProfile(ctx, donor); err != nil {
		return nil, err
	}

	return donor, nil
}

func (s *Service) GetMyDonorProfile(ctx context.Context, userID string) (*DonorProfile, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	return s.repo.GetDonorProfileByUserID(ctx, userID)
}

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

func (s *Service) CreateNGOProfile(
	ctx context.Context,
	userID string,
	role string,
	organization string,
	registrationNo string,
) (*NGOProfile, error) {

	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	// IMPORTANT FIX: NGO is NOT a role
	if role != "USER" {
		return nil, errors.New("only authenticated users can apply as NGO")
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

func (s *Service) GetMyNGOProfile(ctx context.Context, userID string) (*NGOProfile, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	return s.repo.GetNGOProfileByUserID(ctx, userID)
}

/* ===================== ADMIN ===================== */

func (s *Service) VerifyNGO(
	ctx context.Context,
	adminUserID string,
	adminRole string,
	targetUserID string,
) error {

	if adminUserID == "" {
		return errors.New("unauthorized")
	}

	if adminRole != "ADMIN" {
		return errors.New("forbidden")
	}

	if targetUserID == "" {
		return errors.New("target user id required")
	}

	return s.repo.UpdateNGOVerification(ctx, targetUserID, true, &adminUserID)
}
