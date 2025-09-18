package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/adityawaradkar/gratia/auth_service/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service provides user-related business logic
type Service struct {
	repo       RepositoryInterface
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewService creates a new Service instance
func NewService(repo RepositoryInterface, jwtSecret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		repo:       repo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// Input Types
type RegisterInput struct {
	Email    string
	Password string
	Role     string
}

type LoginInput struct {
	Email    string
	Password string
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RegisterUser registers a new user
func (s *Service) RegisterUser(ctx context.Context, input RegisterInput) (*User, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, errors.New("email and password required")
	}

	if existing, _ := s.repo.GetUserByEmail(ctx, email); existing != nil {
		return nil, errors.New("user already exists")
	}

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

// LoginUser authenticates a user and returns tokens
func (s *Service) LoginUser(ctx context.Context, input LoginInput, userAgent, ip string) (*TokenPair, error) {
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

	// Generate access token
	access, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	rt := &RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshToken,
		ExpiresAt: time.Now().Add(s.refreshTTL),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.SaveRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	session := &Session{
		UserID:         user.ID,
		RefreshTokenID: rt.ID,
		UserAgent:      userAgent,
		IPAddress:      ip,
		IsCurrent:      true,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.SaveSession(ctx, session); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: access, RefreshToken: refreshToken}, nil
}

// RefreshTokens refreshes access and refresh tokens
func (s *Service) RefreshTokens(ctx context.Context, refreshToken string, userAgent, ip string) (*TokenPair, error) {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil || rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	access, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefresh, err := utils.GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	newRT := &RefreshToken{
		UserID:    user.ID,
		TokenHash: newRefresh,
		ExpiresAt: time.Now().Add(s.refreshTTL),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.SaveRefreshToken(ctx, newRT); err != nil {
		return nil, err
	}

	if err := s.repo.RevokeRefreshToken(ctx, rt.ID); err != nil {
		return nil, err
	}

	session := &Session{
		UserID:         user.ID,
		RefreshTokenID: newRT.ID,
		UserAgent:      userAgent,
		IPAddress:      ip,
		IsCurrent:      true,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.SaveSession(ctx, session); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: access, RefreshToken: newRefresh}, nil
}

// Logout revokes a refresh token and deletes the session
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return errors.New("invalid token")
	}

	if err := s.repo.RevokeRefreshToken(ctx, rt.ID); err != nil {
		return err
	}

	if err := s.repo.DeleteSessionByToken(ctx, rt.ID); err != nil {
		return err
	}

	return nil
}

// generateAccessToken creates a JWT access token
func (s *Service) generateAccessToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.accessTTL).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// GetUserByID fetches a user by ID
func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}
