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

type Service struct {
	repo       Repo
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(repo Repo, jwtSecret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		repo:       repo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
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

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (s *Service) RegisterUser(ctx context.Context, input RegisterInput) (*User, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, errors.New("email and password required")
	}

	existing, _ := s.repo.GetUserByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	role := input.Role
	if role == "" {
		role = "USER"
	}

	user := &User{
		Email:         email,
		PasswordHash:  string(hashBytes),
		Role:          role,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, input LoginInput, userAgent, ip string) (*TokenPair, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, errors.New("email or password missing")
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

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

	err = s.repo.SaveRefreshToken(ctx, rt)
	if err != nil {
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

	err = s.repo.SaveSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken, userAgent, ip string) (*TokenPair, error) {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil || rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := s.generateAccessToken(user)
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

	err = s.repo.SaveRefreshToken(ctx, newRT)
	if err != nil {
		return nil, err
	}

	err = s.repo.RevokeRefreshToken(ctx, rt.ID)
	if err != nil {
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

	err = s.repo.SaveSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefresh,
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return errors.New("invalid token")
	}

	err = s.repo.RevokeRefreshToken(ctx, rt.ID)
	if err != nil {
		return err
	}

	err = s.repo.DeleteSessionByToken(ctx, rt.ID)
	if err != nil {
		return err
	}

	return nil
}

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

func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) GenerateResetToken(ctx context.Context, email string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("email not found")
	}

	token := uuid.NewString()
	err = s.repo.StoreResetToken(ctx, user.ID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
	user, err := s.repo.FindByResetToken(ctx, resetToken)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = s.repo.UpdatePassword(ctx, user.ID, newHash)
	if err != nil {
		return err
	}

	err = s.repo.ClearResetToken(ctx, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GenerateEmailVerificationToken(ctx context.Context, userID string) (string, error) {
	token := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)

	err := s.repo.StoreEmailVerification(ctx, userID, token, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.repo.VerifyEmailToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	err = s.repo.MarkEmailVerified(ctx, user.ID, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateToken(ctx context.Context, tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return claims, nil
	}

	return nil, errors.New("could not parse claims")
}
