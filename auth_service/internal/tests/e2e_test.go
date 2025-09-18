package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adityawaradkar/gratia/auth_service/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupHandler(t *testing.T) *auth.Handler {
	// Connect to a real test database (use a separate DB for tests!)
	dbURL := "postgres://username:password@localhost:5432/gratia_test?sslmode=disable"
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	repo := auth.NewRepository(dbpool)
	service := auth.NewService(repo, "testsecret", 15*time.Minute, 24*time.Hour)
	return auth.NewHandler(service)
}

func TestAuthE2E(t *testing.T) {
	handler := setupHandler(t)

	// --- REGISTER ---
	registerBody := map[string]string{
		"email":    "e2euser@example.com",
		"password": "Test@1234",
	}
	registerBytes, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(registerBytes))
	w := httptest.NewRecorder()
	handler.RegisterUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", w.Code)
	}

	var registeredUser auth.User
	json.NewDecoder(w.Body).Decode(&registeredUser)
	if registeredUser.Email != registerBody["email"] {
		t.Fatalf("expected email %s, got %s", registerBody["email"], registeredUser.Email)
	}

	// --- LOGIN ---
	loginBody := map[string]string{
		"email":    "e2euser@example.com",
		"password": "Test@1234",
	}
	loginBytes, _ := json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(loginBytes))
	req.Header.Set("User-Agent", "e2e-test-agent")
	w = httptest.NewRecorder()
	handler.LoginUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", w.Code)
	}

	var tokens auth.TokenPair
	json.NewDecoder(w.Body).Decode(&tokens)
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatal("expected non-empty access and refresh tokens")
	}

	// --- REFRESH ---
	refreshBody := map[string]string{"refreshToken": tokens.RefreshToken}
	refreshBytes, _ := json.Marshal(refreshBody)
	req = httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(refreshBytes))
	req.Header.Set("User-Agent", "e2e-test-agent")
	w = httptest.NewRecorder()
	handler.RefreshTokens(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK on refresh, got %d", w.Code)
	}

	var newTokens auth.TokenPair
	json.NewDecoder(w.Body).Decode(&newTokens)
	if newTokens.AccessToken == "" || newTokens.RefreshToken == "" {
		t.Fatal("expected new access and refresh tokens after refresh")
	}

	// --- LOGOUT ---
	logoutBody := map[string]string{"refreshToken": newTokens.RefreshToken}
	logoutBytes, _ := json.Marshal(logoutBody)
	req = httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader(logoutBytes))
	req.Header.Set("User-Agent", "e2e-test-agent")
	w = httptest.NewRecorder()
	handler.Logout(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 No Content on logout, got %d", w.Code)
	}
}
