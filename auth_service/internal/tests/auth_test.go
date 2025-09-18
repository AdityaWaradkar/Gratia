package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/adityawaradkar/gratia/auth_service/internal/auth"
	"github.com/adityawaradkar/gratia/auth_service/internal/utils"
)

// MockRepository implements the auth.Repository interface for testing
type MockRepository struct {
	users         map[string]*auth.User
	refreshTokens map[string]*auth.RefreshToken
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users:         make(map[string]*auth.User),
		refreshTokens: make(map[string]*auth.RefreshToken),
	}
}

func (m *MockRepository) CreateUser(ctx context.Context, user *auth.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	user.ID = "mockID"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.Email] = user
	return nil
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockRepository) GetUserByID(ctx context.Context, id string) (*auth.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockRepository) UpdateEmailVerified(ctx context.Context, userID string) error {
	for _, u := range m.users {
		if u.ID == userID {
			u.EmailVerified = true
			u.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockRepository) UpdatePassword(ctx context.Context, userID string, newHash string) error {
	for _, u := range m.users {
		if u.ID == userID {
			u.PasswordHash = newHash
			u.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockRepository) SaveRefreshToken(ctx context.Context, rt *auth.RefreshToken) error {
	rt.ID = "mockRTID"
	m.refreshTokens[rt.TokenHash] = rt
	return nil
}

func (m *MockRepository) GetRefreshToken(ctx context.Context, token string) (*auth.RefreshToken, error) {
	rt, ok := m.refreshTokens[token]
	if !ok {
		return nil, errors.New("refresh token not found")
	}
	return rt, nil
}

func (m *MockRepository) RevokeRefreshToken(ctx context.Context, id string) error {
	for k, v := range m.refreshTokens {
		if v.ID == id {
			v.Revoked = true
			m.refreshTokens[k] = v
			return nil
		}
	}
	return errors.New("token not found")
}

func (m *MockRepository) SaveSession(ctx context.Context, s *auth.Session) error {
	s.ID = "mockSessionID"
	return nil
}

func (m *MockRepository) DeleteSessionByToken(ctx context.Context, refreshTokenID string) error {
	return nil
}

// -------- Unit Tests --------

func TestRegisterUser(t *testing.T) {
	repo := NewMockRepository()
	service := auth.NewService(repo, "testsecret", 15*time.Minute, 24*time.Hour)

	user, err := service.RegisterUser(context.Background(), auth.RegisterInput{
		Email:    "test@example.com",
		Password: "Test@1234",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email)
	}
}

func TestRegisterDuplicateUser(t *testing.T) {
	repo := NewMockRepository()
	service := auth.NewService(repo, "testsecret", 15*time.Minute, 24*time.Hour)

	_, _ = service.RegisterUser(context.Background(), auth.RegisterInput{
		Email:    "dup@example.com",
		Password: "Test@1234",
	})
	_, err := service.RegisterUser(context.Background(), auth.RegisterInput{
		Email:    "dup@example.com",
		Password: "Test@1234",
	})
	if err == nil {
		t.Fatal("expected error for duplicate user, got nil")
	}
}

func TestLoginUser(t *testing.T) {
	repo := NewMockRepository()
	service := auth.NewService(repo, "testsecret", 15*time.Minute, 24*time.Hour)

	_, _ = service.RegisterUser(context.Background(), auth.RegisterInput{
		Email:    "login@example.com",
		Password: "Test@1234",
	})

	tokens, err := service.LoginUser(context.Background(), auth.LoginInput{
		Email:    "login@example.com",
		Password: "Test@1234",
	}, "test-agent", "127.0.0.1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Error("expected non-empty access and refresh tokens")
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "MyP@ssword123"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("hashing failed: %v", err)
	}

	if err := utils.VerifyPassword(hash, password); err != nil {
		t.Errorf("password verification failed: %v", err)
	}
}
