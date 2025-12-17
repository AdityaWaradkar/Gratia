package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/adityawaradkar/gratia/auth_service/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/* ===================== SERVICE ===================== */

type Service struct {
	repo Repo

	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration

	userServiceURL string
	httpClient     *http.Client
	logger         *log.Logger
}

func NewService(
	repo Repo,
	jwtSecret string,
	accessTTL, refreshTTL time.Duration,
	userServiceURL string,
	logger *log.Logger,
) *Service {
	return &Service{
		repo:           repo,
		jwtSecret:      []byte(jwtSecret),
		accessTTL:      accessTTL,
		refreshTTL:     refreshTTL,
		userServiceURL: userServiceURL,
		httpClient:     &http.Client{Timeout: 5 * time.Second},
		logger:         logger,
	}
}

/* ===================== INPUT MODELS ===================== */

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

/* ===================== ROLE NORMALIZATION ===================== */

// Converts Auth roles → User Service roles
func normalizeUserServiceRole(role string) string {
	switch role {
	case RoleDonor:
		return "Donor"
	case RoleNGO:
		return "NGO"
	case RoleAdmin:
		return "Admin"
	default:
		return "Donor"
	}
}

/* ===================== AUTH ===================== */

func (s *Service) RegisterUser(
	ctx context.Context,
	input RegisterInput,
	callerRole string,
) (*User, error) {

	now := time.Now()
	email := strings.TrimSpace(strings.ToLower(input.Email))

	if email == "" || input.Password == "" {
		return nil, errors.New("email and password required")
	}

	existing, _ := s.repo.GetUserByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	role := strings.ToUpper(input.Role)
	if role == "" {
		role = RoleDonor
	} else {
		if role != RoleDonor && role != RoleNGO && callerRole != RoleAdmin {
			return nil, errors.New("unauthorized role assignment")
		}
	}

	user := &User{
		Email:         email,
		PasswordHash:  string(passwordHash),
		Role:          role,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Async, best-effort user-profile creation
	go s.notifyUserService(
		user.ID,
		normalizeUserServiceRole(role),
	)

	return user, nil
}

func (s *Service) LoginUser(
	ctx context.Context,
	input LoginInput,
	userAgent, ip string,
) (*TokenPair, error) {

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, errors.New("email or password missing")
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	); err != nil {
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

	now := time.Now()

	rt := &RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshToken,
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

	if err := s.repo.SaveSession(ctx, session); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshTokens(
	ctx context.Context,
	refreshToken, userAgent, ip string,
) (*TokenPair, error) {

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

	if err := s.repo.SaveRefreshToken(ctx, newRT); err != nil {
		return nil, err
	}

	_ = s.repo.RevokeRefreshToken(ctx, rt.ID)

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

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return errors.New("invalid token")
	}

	if err := s.repo.RevokeRefreshToken(ctx, rt.ID); err != nil {
		return err
	}

	return s.repo.DeleteSessionByToken(ctx, rt.ID)
}

/* ===================== PASSWORD RESET ===================== */

func (s *Service) GenerateResetToken(ctx context.Context, email string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("email not found")
	}

	token := uuid.NewString()
	if err := s.repo.StoreResetToken(ctx, user.ID, token); err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) ResetPassword(
	ctx context.Context,
	resetToken, newPassword string,
) error {

	user, err := s.repo.FindByResetToken(ctx, resetToken)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, user.ID, newHash); err != nil {
		return err
	}

	return s.repo.ClearResetToken(ctx, user.ID)
}

/* ===================== INTERNAL ===================== */

func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
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

/* ===================== USER SERVICE INTEGRATION ===================== */

func (s *Service) notifyUserService(userID, role string) {
	s.logger.Println("calling user service to create profile for:", userID)

	payload := map[string]string{
		"user_id": userID,
		"role":    role,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		s.logger.Println("user-service marshal error:", err)
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.userServiceURL+"/internal/users",
		bytes.NewBuffer(body),
	)
	if err != nil {
		s.logger.Println("user-service request creation failed:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Println("user-service call failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		s.logger.Println("user-service unexpected status:", resp.StatusCode)
	}
}
