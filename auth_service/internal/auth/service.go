package auth

import (
    "context"
    "errors"
    "strings"
    "time"

    "github.com/adityawaradkar/gratia/auth_service/internal/utils"
    "github.com/google/uuid"
)

// Service orchestrates all business logic and database interactions for authentication
type Service struct {
    repo       Repository
    jwtSecret  []byte
    accessTTL  time.Duration
    refreshTTL time.Duration
}

// RegisterInput encapsulates the required data payload for new user onboarding
type RegisterInput struct {
    Email    string
    Password string
    Role     Role
}

// LoginInput encapsulates the required data payload for authenticating existing users
type LoginInput struct {
    Email    string
    Password string
}

// TokenPair delivers the stateless access token and stateful refresh token back to the client
type TokenPair struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
}

// NewService initializes the authentication service with injected dependencies and configurations
func NewService(repo Repository, jwtSecret string, accessTTL, refreshTTL time.Duration) *Service {
    return &Service{
        repo:       repo,
        jwtSecret:  []byte(jwtSecret),
        accessTTL:  accessTTL,
        refreshTTL: refreshTTL,
    }
}

// RegisterUser securely creates a new account and enforces role-based access restrictions
func (s *Service) RegisterUser(ctx context.Context, input RegisterInput, callerRole Role) (*User, error) {
    email := strings.TrimSpace(strings.ToLower(input.Email))
    
    if err := ValidateEmail(email); err != nil {
        return nil, err
    }
    if err := ValidatePassword(input.Password); err != nil {
        return nil, err
    }

    if user, _ := s.repo.GetUserByEmail(ctx, email); user != nil {
        return nil, errors.New("user already exists")
    }

    role := input.Role
    if role == "" {
        role = RoleUser
    }

    if !ValidRoles[role] {
        return nil, errors.New("invalid role")
    }

    if role == RoleAdmin && callerRole != RoleAdmin {
        return nil, errors.New("unauthorized role assignment")
    }

    hash, err := utils.HashPassword(input.Password)
    if err != nil {
        return nil, errors.New("failed to hash password")
    }

    user := &User{
        Email:         email,
        PasswordHash:  hash,
        Role:          role,
        EmailVerified: false,
    }

    if err := s.repo.CreateUser(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

// LoginUser verifies credentials and establishes a new secure session with token generation
func (s *Service) LoginUser(ctx context.Context, input LoginInput, userAgent, ip string) (*TokenPair, error) {
    email := strings.TrimSpace(strings.ToLower(input.Email))
    
    user, err := s.repo.GetUserByEmail(ctx, email)
    if err != nil || user == nil {
        return nil, errors.New("invalid credentials")
    }

    if err := utils.ComparePassword(user.PasswordHash, input.Password); err != nil {
        return nil, errors.New("invalid credentials")
    }

    accessToken, err := utils.GenerateAccessToken(user.ID, string(user.Role), string(s.jwtSecret), s.accessTTL)
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
        Token:     refreshToken,
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

    return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// RefreshTokens rotates long-lived credentials to maintain continuous secure access
func (s *Service) RefreshTokens(ctx context.Context, refreshToken, userAgent, ip string) (*TokenPair, error) {
    rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
    if err != nil || rt == nil || rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
        return nil, errors.New("invalid refresh token")
    }

    user, err := s.repo.GetUserByID(ctx, rt.UserID)
    if err != nil || user == nil {
        return nil, errors.New("user not found")
    }

    accessToken, err := utils.GenerateAccessToken(user.ID, string(user.Role), string(s.jwtSecret), s.accessTTL)
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
        Token:     newRefreshToken,
        ExpiresAt: now.Add(s.refreshTTL),
        Revoked:   false,
        CreatedAt: now,
    }

    if err := s.repo.RevokeRefreshToken(ctx, rt.ID); err != nil {
        return nil, err
    }
    if err := s.repo.DeleteSessionByToken(ctx, rt.ID); err != nil {
        return nil, err
    }
    if err := s.repo.SaveRefreshToken(ctx, newRT); err != nil {
        return nil, err
    }

    session := &Session{
        UserID:         user.ID,
        RefreshTokenID: newRT.ID,
        UserAgent:      userAgent,
        IPAddress:      ip,
        IsCurrent:      true,
        CreatedAt:      now,
    }

    if err := s.repo.SaveSession(ctx, session); err != nil {
        return nil, err
    }

    return &TokenPair{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

// Logout terminates an active session by actively revoking the provided refresh token
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
    rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
    if err != nil || rt == nil {
        return errors.New("invalid token")
    }

    if err := s.repo.RevokeRefreshToken(ctx, rt.ID); err != nil {
        return err
    }
    return s.repo.DeleteSessionByToken(ctx, rt.ID)
}

// GenerateResetToken initializes the secure account recovery process
func (s *Service) GenerateResetToken(ctx context.Context, email string) (string, error) {
    user, err := s.repo.GetUserByEmail(ctx, email)
    if err != nil || user == nil {
        return "", errors.New("user not found")
    }

    token := uuid.NewString()
    if err := s.repo.StoreResetToken(ctx, user.ID, token); err != nil {
        return "", err
    }

    return token, nil
}

// ResetPassword safely overwrites credentials using a verified recovery token
func (s *Service) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
    if err := ValidatePassword(newPassword); err != nil {
        return err
    }

    user, err := s.repo.FindByResetToken(ctx, resetToken)
    if err != nil || user == nil {
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

// GetUserByID fetches a user's domain profile using their unique identifier
func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
    return s.repo.GetUserByID(ctx, userID)
}