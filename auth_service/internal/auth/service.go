package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Service provides user-related business logic
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type RegisterInput struct {
	Email    string
	Password string
	Role     string
}

type LoginInput struct {
	Email    string
	Password string
}

// RegisterUser creates a new user
func (s *Service) RegisterUser(ctx context.Context, input RegisterInput) (*User, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, errors.New("email and password required")
	}

	// check if user exists
	if existing, _ := s.repo.GetUserByEmail(ctx, email); existing != nil {
		return nil, errors.New("user already exists")
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &User{
		Email:         email,
		PasswordHash:  string(hashed),
		Role:          "USER",
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// LoginUser checks credentials
func (s *Service) LoginUser(ctx context.Context, input LoginInput) (*User, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, errors.New("email and password required")
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
