package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/adityawaradkar/gratia/auth_service/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service *Service
	db      *pgxpool.Pool
}

func NewHandler(service *Service, db *pgxpool.Pool) *Handler {
	return &Handler{service: service, db: db}
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
	ResetToken  string `json:"resetToken"`
	NewPassword string `json:"newPassword"`
}

type JSONError struct {
	Message string `json:"message"`
}

/* ===================== HELPERS ===================== */

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(JSONError{Message: msg})
}

/* ===================== AUTH ===================== */

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if ValidateEmail(req.Email) != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid email")
		return
	}

	if ValidatePassword(req.Password) != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid password")
		return
	}

	role := strings.ToUpper(req.Role)

	callerRole := middleware.UserRole(r.Context())
	if callerRole == "" {
		callerRole = RoleDonor
	}

	user, err := h.service.RegisterUser(
		r.Context(),
		RegisterInput{
			Email:    req.Email,
			Password: req.Password,
			Role:     role,
		},
		callerRole,
	)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if ValidateEmail(req.Email) != nil || ValidatePassword(req.Password) != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	tokens, err := h.service.LoginUser(
		r.Context(),
		LoginInput{Email: req.Email, Password: req.Password},
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tokens)
}

func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := h.service.RefreshTokens(
		r.Context(),
		req.RefreshToken,
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tokens)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		writeJSONError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "user not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

/* ===================== PASSWORD ===================== */

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.service.GenerateResetToken(r.Context(), req.Email)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"resetToken": token,
		"note":       "This would be emailed in production",
	})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.ResetPassword(
		r.Context(),
		req.ResetToken,
		req.NewPassword,
	); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid reset token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* ===================== HEALTH ===================== */

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h.db.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status":"error","details":"database unreachable"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
