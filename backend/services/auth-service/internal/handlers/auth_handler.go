package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adityawaradkar/gratia-auth/internal/middleware"
	"github.com/adityawaradkar/gratia-auth/internal/utils"
	"github.com/adityawaradkar/gratia-auth/internal/models"
	"github.com/adityawaradkar/gratia-auth/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration requests
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FullName    string `json:"full_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Role        string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:       req.Email,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Role:        req.Role,
	}

	userID, err := h.authService.Register(r.Context(), user, req.Password)
	if err != nil {
		if err == service.ErrUserExists {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := map[string]interface{}{
		"message": "User registered successfully.",
		"user_id": userID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Login authenticates a user and returns JWT access and refresh tokens
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			if err.Error() == "user account is disabled" {
				http.Error(w, err.Error(), http.StatusForbidden)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}
		return
	}

	resp := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900, // 15 minutes
		"token_type":    "Bearer",
	}

	json.NewEncoder(w).Encode(resp)
}

// Me returns the currently authenticated user's profile info
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.UserIDKey)
	if userIDRaw == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Logout invalidates the user's session (currently a dummy implementation)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Logged out successfully."}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Refresh issues a new access token given a valid refresh token
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	claims, err := utils.ValidateJWT(req.RefreshToken, h.authService.Secret())
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.authService.GenerateJWT(user, 15*time.Minute)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"access_token": accessToken,
		"expires_in":   900,
		"token_type":   "Bearer",
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement password reset email sending
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Forgot password feature not implemented yet",
	})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement password reset using token and new password
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Reset password feature not implemented yet",
	})
}

// Routes defines all auth-related routes and applies JWT middleware to protected routes
func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Public routes
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/forgot-password", h.ForgotPassword)
	r.Post("/reset-password", h.ResetPassword)

	// Protected routes
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JWTAuthMiddleware(h.authService.Secret()))
		protected.Post("/logout", h.Logout)
		protected.Get("/me", h.Me)
	})

	return r
}

