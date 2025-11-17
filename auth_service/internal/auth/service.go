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

func (s *Service) RegisterUser(ctx context.Context, input RegisterInput, callerRole string) (*User, error) {
	now := time.Now()
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

	role := strings.ToUpper(input.Role)
	if role == "" {
		role = "USER"
	} else {
		if role != "USER" && callerRole != "ADMIN" {
			return nil, errors.New("unauthorized role assignment")
		}
		if !ValidRoles[role] {
			return nil, errors.New("invalid role")
		}
	}

	user := &User{
		Email:         email,
		PasswordHash:  string(hashBytes),
		Role:          role,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, input LoginInput, userAgent, ip string) (*TokenPair, error) {
	now := time.Now()
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
		ExpiresAt: now.Add(s.refreshTTL),
		Revoked:   false,
		CreatedAt: now,
	}

	err = s.repo.SaveRefreshToken(ctx, rt)
	if err != nil {
		return nil, err
	}

	sessions, _ := s.repo.GetSessionsByID(ctx, user.ID)
	for _, sess := range sessions {
		_ = s.repo.DeleteSessionByID(ctx, sess.ID)
	}

	session := &Session{
		UserID:         user.ID,
		RefreshTokenID: rt.ID,
		UserAgent:      userAgent,
		IPAddress:      ip,
		IsCurrent:      true,
		CreatedAt:      now,
	}

	err = s.repo.SaveSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken, userAgent, ip string) (*TokenPair, error) {
	now := time.Now()

	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil || rt.Revoked || rt.ExpiresAt.Before(now) {
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

	newRefreshToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	newRT := &RefreshToken{
		UserID:    user.ID,
		TokenHash: newRefreshToken,
		ExpiresAt: now.Add(s.refreshTTL),
		Revoked:   false,
		CreatedAt: now,
	}

	err = s.repo.SaveRefreshToken(ctx, newRT)
	if err != nil {
		return nil, err
	}

	err = s.repo.RevokeRefreshToken(ctx, rt.ID)
	if err != nil {
		return nil, err
	}

	sessions, _ := s.repo.GetSessionsByID(ctx, user.ID)
	for _, sess := range sessions {
		if sess.RefreshTokenID == rt.ID {
			_ = s.repo.DeleteSessionByID(ctx, sess.ID)
		}
	}

	session := &Session{
		UserID:         user.ID,
		RefreshTokenID: newRT.ID,
		UserAgent:      userAgent,
		IPAddress:      ip,
		IsCurrent:      true,
		CreatedAt:      now,
	}

	err = s.repo.SaveSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
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

	return s.repo.DeleteSessionByToken(ctx, rt.ID)
}

func (s *Service) generateAccessToken(user *User) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   now.Add(s.accessTTL).Unix(),
		"iat":   now.Unix(),
		"jti":   uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, errors.New("email required")
	}
	return s.repo.GetUserByEmail(ctx, email)
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

	return s.repo.ClearResetToken(ctx, user.ID)
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

	return s.repo.MarkEmailVerified(ctx, user.ID, token)
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

func (s *Service) GetSessions(ctx context.Context, userID string) ([]*Session, error) {
	return s.repo.GetSessionsByID(ctx, userID)
}

func (s *Service) DeleteSession(ctx context.Context, userID, sessionID string) error {
	sessions, err := s.repo.GetSessionsByID(ctx, userID)
	if err != nil {
		return err
	}

	var found bool
	for _, session := range sessions {
		if session.ID == sessionID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("session not found or unauthorized")
	}

	return s.repo.DeleteSessionByID(ctx, sessionID)
}

func (s *Service) UpdateUserRole(ctx context.Context, callerRole, userID, newRole string) error {
	if callerRole != "ADMIN" {
		return errors.New("only admin can update roles")
	}

	newRole = strings.ToUpper(newRole)
	if !ValidRoles[newRole] {
		return errors.New("invalid role")
	}

	return s.repo.UpdateUserRole(ctx, userID, newRole)
}
