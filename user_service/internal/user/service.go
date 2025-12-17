package user

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMyProfile(ctx context.Context, userID string) (*UserProfile, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) UpdateMyProfile(
	ctx context.Context,
	userID string,
	name string,
	phone string,
	address string,
) (*UserProfile, error) {

	return s.repo.UpdateByUserID(ctx, userID, name, phone, address)
}

func (s *Service) CreateProfile(
	ctx context.Context,
	userID string,
	role string,
) (*UserProfile, error) {

	// default empty profile
	return s.repo.CreateProfile(ctx, userID, role)
}
