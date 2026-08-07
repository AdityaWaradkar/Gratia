package auth

import (
    "encoding/json"
    "net/http"

    "github.com/adityawaradkar/gratia/auth_service/internal/middleware"
)

// Handler maps HTTP requests directly into the core business service
type Handler struct {
    service *Service
}

// NewHandler initializes the HTTP delivery layer with its required service dependency
func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// RegisterRequest structures the incoming JSON payload for new account creation
type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     Role   `json:"role,omitempty"`
}

// LoginRequest structures the incoming JSON payload for authentication attempts
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// RefreshRequest structures the incoming JSON payload for token rotation
type RefreshRequest struct {
    RefreshToken string `json:"refreshToken"`
}

// ForgotPasswordRequest structures the incoming JSON payload for account recovery
type ForgotPasswordRequest struct {
    Email string `json:"email"`
}

// ResetPasswordRequest structures the incoming JSON payload for credential replacement
type ResetPasswordRequest struct {
    Token       string `json:"token"`
    NewPassword string `json:"newPassword"`
}

// writeJSON centralizes response formatting to enforce a uniform API contract
func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if data != nil {
        _ = json.NewEncoder(w).Encode(data)
    }
}

// writeError wraps error messages into consistent JSON objects for client consumption
func writeError(w http.ResponseWriter, status int, message string) {
    writeJSON(w, status, map[string]string{"error": message})
}

// RegisterUser processes client onboarding and translates domain errors to HTTP statuses
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid json payload")
        return
    }

    callerRole := Role(middleware.UserRole(r.Context()))
    user, err := h.service.RegisterUser(r.Context(), RegisterInput{
        Email:    req.Email,
        Password: req.Password,
        Role:     req.Role,
    }, callerRole)

    if err != nil {
        writeError(w, http.StatusConflict, err.Error())
        return
    }

    writeJSON(w, http.StatusCreated, user)
}

// LoginUser issues secure credentials based on verified client inputs
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid json payload")
        return
    }

    tokens, err := h.service.LoginUser(r.Context(), LoginInput{
        Email:    req.Email,
        Password: req.Password,
    }, r.UserAgent(), r.RemoteAddr)

    if err != nil {
        writeError(w, http.StatusUnauthorized, err.Error())
        return
    }

    writeJSON(w, http.StatusOK, tokens)
}

// RefreshTokens handles stateless session extension without requiring raw passwords
func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
    var req RefreshRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid json payload")
        return
    }

    tokens, err := h.service.RefreshTokens(r.Context(), req.RefreshToken, r.UserAgent(), r.RemoteAddr)
    if err != nil {
        writeError(w, http.StatusUnauthorized, err.Error())
        return
    }

    writeJSON(w, http.StatusOK, tokens)
}

// Logout actively tears down sessions to prevent stolen refresh token exploitation
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
    var req RefreshRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid json payload")
        return
    }

    if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
        writeError(w, http.StatusUnauthorized, err.Error())
        return
    }

    writeJSON(w, http.StatusNoContent, nil)
}

// ForgotPassword generates safe recovery links and maintains an opaque success response
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
    var req ForgotPasswordRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid json payload")
        return
    }

    token, err := h.service.GenerateResetToken(r.Context(), req.Email)
    if err != nil {
        // Always return 200 to prevent user enumeration attacks
        writeJSON(w, http.StatusOK, map[string]string{"message": "if the email exists, a reset link was sent"})
        return
    }

    // In a real application, you would email this token instead of returning it directly
    writeJSON(w, http.StatusOK, map[string]string{"resetToken": token})
}

// ResetPassword consumes a recovery token to establish new credentials
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
    var req ResetPasswordRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid json payload")
        return
    }

    if err := h.service.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }

    writeJSON(w, http.StatusNoContent, nil)
}

// GetCurrentUser returns the profile data associated with the incoming JWT
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
    userID := middleware.UserID(r.Context())
    user, err := h.service.GetUserByID(r.Context(), userID)
    
    if err != nil || user == nil {
        writeError(w, http.StatusNotFound, "user not found")
        return
    }

    writeJSON(w, http.StatusOK, user)
}

// HealthCheck provides a lightweight endpoint for container orchestration readiness probes
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "auth_service"})
}