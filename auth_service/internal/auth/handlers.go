package auth

import (
	"encoding/json"
	"net/http"

	"github.com/adityawaradkar/gratia/auth_service/internal/middleware"
)

// Handler handles HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates auth handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

/* ===================== REQUEST MODELS ===================== */

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

/* ===================== HELPERS ===================== */

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"message": message})
}

/* ===================== AUTH ROUTES ===================== */

// RegisterUser handles user signup
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := ValidateEmail(req.Email); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ValidatePassword(req.Password); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	callerRole := middleware.UserRole(r.Context())

	user, err := h.service.RegisterUser(
		r.Context(),
		RegisterInput{
			Email:    req.Email,
			Password: req.Password,
			Role:     req.Role,
		},
		callerRole,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

// LoginUser handles login
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := h.service.LoginUser(
		r.Context(),
		LoginInput{
			Email:    req.Email,
			Password: req.Password,
		},
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

// RefreshTokens issues new tokens
func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := h.service.RefreshTokens(
		r.Context(),
		req.RefreshToken,
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

// Logout revokes refresh token
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* ===================== PASSWORD RESET ===================== */

// ForgotPassword generates reset token
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.service.GenerateResetToken(r.Context(), req.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"resetToken": token,
	})
}

// ResetPassword updates password
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.ResetPassword(
		r.Context(),
		req.Token,
		req.NewPassword,
	); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* ===================== INTERNAL ===================== */

// GetCurrentUser returns authenticated user
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

//Health Check
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
