package auth

import (
	"context"
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

type JSONError struct {
	Message string `json:"message"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	ResetToken  string `json:"resetToken"`
	NewPassword string `json:"newPassword"`
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONError{Message: msg})
}

func isValidRole(role string) bool {
	switch role {
	case RoleUser, RoleAdmin, RoleMod:
		return true
	default:
		return false
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
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

	role := req.Role
	if role == "" {
		role = RoleUser
	} else if !isValidRole(role) {
		writeJSONError(w, http.StatusBadRequest, "invalid role")
		return
	}

	callerRole := middleware.UserRole(r.Context())
	if callerRole == "" {
		callerRole = RoleUser
	}

	user, err := h.service.RegisterUser(
		context.Background(),
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
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
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

	tokens, err := h.service.LoginUser(
		context.Background(),
		LoginInput{Email: req.Email, Password: req.Password},
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	tokens, err := h.service.RefreshTokens(
		context.Background(),
		req.RefreshToken,
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err := h.service.Logout(context.Background(), req.RefreshToken)
	if err != nil {
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	token, err := h.service.GenerateResetToken(context.Background(), req.Email)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"resetToken": token,
		"note":       "This would be emailed in production",
	})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.ResetToken == "" || req.NewPassword == "" {
		writeJSONError(w, http.StatusBadRequest, "token and password required")
		return
	}

	err := h.service.ResetPassword(context.Background(), req.ResetToken, req.NewPassword)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid reset token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Token == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err := h.service.VerifyEmail(r.Context(), req.Token)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Token == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	claims, err := h.service.ValidateToken(r.Context(), req.Token)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}

func (h *Handler) GenerateEmailVerification(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Email == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "user not found")
		return
	}

	token, err := h.service.GenerateEmailVerificationToken(r.Context(), user.ID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token = strings.ReplaceAll(token, "-", "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"verificationToken": token,
		"note":              "This would be emailed in production",
	})
}

func (h *Handler) GetSessions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	sessions, err := h.service.GetSessions(r.Context(), userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not fetch sessions")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	prefix := "/auth/sessions/"
	sessionID := r.URL.Path[len(prefix):]

	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "session ID required")
		return
	}

	userID := middleware.UserID(r.Context())
	if userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.service.DeleteSession(r.Context(), userID, sessionID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ResendEmailVerification(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Email == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "user not found")
		return
	}

	if user.EmailVerified {
		writeJSONError(w, http.StatusBadRequest, "email already verified")
		return
	}

	token, err := h.service.GenerateEmailVerificationToken(r.Context(), user.ID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"verificationToken": token,
		"note":              "This would be emailed in production",
	})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := h.db.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status":"error","details":"database unreachable"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
