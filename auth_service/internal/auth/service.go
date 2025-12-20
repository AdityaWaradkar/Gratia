package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/adityawaradkar/gratia/auth_service/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service contains business logic for auth
type Service struct {
	repo       Repository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// RegisterInput represents registration input
type RegisterInput struct {
	Email    string
	Password string
	Role     string
}

// LoginInput represents login input
type LoginInput struct {
	Email    string
	Password string
}

// TokenPair represents JWT tokens
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// NewService creates auth service
func NewService(
	repo Repository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *Service {
	return &Service{
		repo:       repo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

/* ===================== AUTH ===================== */

// RegisterUser creates a new authenticated account
func (s *Service) RegisterUser(
	ctx context.Context,
	input RegisterInput,
	callerRole string,
) (*User, error) {

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, errors.New("email and password required")
	}

	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return nil, errors.New("user already exists")
	}

	role := strings.ToUpper(input.Role)
	if role == "" {
		role = RoleUser
	}

	if !ValidRoles[role] {
		return nil, errors.New("invalid role")
	}

	if role == RoleAdmin && callerRole != RoleAdmin {
		return nil, errors.New("unauthorized role assignment")
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &User{
		Email:         email,
		PasswordHash:  string(hash),
		Role:          role,
		EmailVerified: false,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser authenticates user and issues tokens
func (s *Service) LoginUser(
	ctx context.Context,
	input LoginInput,
	userAgent string,
	ip string,
) (*TokenPair, error) {

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, errors.New("invalid credentials")
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	rt := &RefreshToken{
		UserID:    user.ID,
		Token: refreshToken,
		ExpiresAt: now.Add(s.refreshTTL),
		Revoked:   false,
		CreatedAt: now,
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
		CreatedAt:      now,
	}

	_ = s.repo.SaveSession(ctx, session)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshTokens rotates refresh token and issues new access token
func (s *Service) RefreshTokens(
	ctx context.Context,
	refreshToken string,
	userAgent string,
	ip string,
) (*TokenPair, error) {

	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil || rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	newRT := &RefreshToken{
		UserID:    user.ID,
		Token: newRefreshToken,
		ExpiresAt: now.Add(s.refreshTTL),
		Revoked:   false,
		CreatedAt: now,
	}

	_ = s.repo.RevokeRefreshToken(ctx, rt.ID)
	_ = s.repo.DeleteSessionByToken(ctx, rt.ID)
	_ = s.repo.SaveRefreshToken(ctx, newRT)

	session := &Session{
		UserID:         user.ID,
		RefreshTokenID: newRT.ID,
		UserAgent:      userAgent,
		IPAddress:      ip,
		IsCurrent:      true,
		CreatedAt:      now,
	}

	_ = s.repo.SaveSession(ctx, session)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout revokes refresh token
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return errors.New("invalid token")
	}

	_ = s.repo.RevokeRefreshToken(ctx, rt.ID)
	_ = s.repo.DeleteSessionByToken(ctx, rt.ID)

	return nil
}

/* ===================== PASSWORD RESET ===================== */

// GenerateResetToken creates password reset token
func (s *Service) GenerateResetToken(ctx context.Context, email string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	token := uuid.NewString()
	if err := s.repo.StoreResetToken(ctx, user.ID, token); err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword updates password using reset token
func (s *Service) ResetPassword(
	ctx context.Context,
	resetToken string,
	newPassword string,
) error {

	user, err := s.repo.FindByResetToken(ctx, resetToken)
	if err != nil {
		return errors.New("invalid reset token")
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, user.ID, hash); err != nil {
		return err
	}

	return s.repo.ClearResetToken(ctx, user.ID)
}

/* ===================== INTERNAL ===================== */

// GetUserByID returns user details
func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// generateAccessToken creates JWT access token
func (s *Service) generateAccessToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(s.accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
