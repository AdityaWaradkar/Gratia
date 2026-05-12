package service

import (
	"context"
	"errors"

	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	"github.com/adityawaradkar/gratia/auth_service/internal/model"
	"github.com/adityawaradkar/gratia/auth_service/internal/repository"
	"github.com/adityawaradkar/gratia/auth_service/internal/utils"
)

type AuthService struct {
	UserRepository *repository.UserRepository
	Config         *config.Config
}

func NewAuthService(
	userRepository *repository.UserRepository,
	cfg *config.Config,
) *AuthService {

	return &AuthService{
		UserRepository: userRepository,
		Config:         cfg,
	}
}

func (s *AuthService) Signup(
	ctx context.Context,
	request *model.SignupRequest,
) error {

	_, err := s.UserRepository.GetUserByEmail(
		ctx,
		request.Email,
	)

	if err == nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(
		request.Password,
	)

	if err != nil {
		return err
	}

	user := &model.User{
		Email:        request.Email,
		PasswordHash: hashedPassword,
		Role:         request.Role,
	}

	err = s.UserRepository.CreateUser(
		ctx,
		user,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(
	ctx context.Context,
	request *model.LoginRequest,
) (*model.AuthResponse, error) {

	user, err := s.UserRepository.GetUserByEmail(
		ctx,
		request.Email,
	)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = utils.ComparePassword(
		user.PasswordHash,
		request.Password,
	)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(
		user.UserID,
		user.Role,
		s.Config.JWTSecret,
		s.Config.JWTExpirationMinutes,
	)

	if err != nil {
		return nil, err
	}

	response := &model.AuthResponse{
		AccessToken: token,
		ExpiresIn:   s.Config.JWTExpirationMinutes * 60,
	}

	return response, nil
}

func (s *AuthService) ValidateToken(
	tokenString string,
) (*model.JWTClaims, error) {

	claims, err := utils.ValidateJWT(
		tokenString,
		s.Config.JWTSecret,
	)

	if err != nil {
		return nil, err
	}

	return claims, nil
}