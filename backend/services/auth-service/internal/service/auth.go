package service

import (
	"context"
	"errors"
	"time"

	"github.com/adityawaradkar/gratia-auth/internal/models"
	"github.com/adityawaradkar/gratia-auth/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrUserExists         = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register returns user ID on success
func (s *AuthService) Register(ctx context.Context, user *models.User, password string) (uuid.UUID, error) {
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, user.Email)
	if existingUser != nil {
		return uuid.Nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}
	user.PasswordHash = string(hash)

	if user.IsActive == nil {
		active := true
		user.IsActive = &active
	}

	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return user.UserID, nil
}

// Login returns accessToken and refreshToken
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", ErrInvalidCredentials
	}

	if user.IsActive != nil && !*user.IsActive {
		return "", "", errors.New("user account is disabled")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := s.generateJWT(user, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateJWT(user, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) generateJWT(user *models.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.UserID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
