package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adityawaradkar/gratia/auth_service/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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
	user, err := h.service.RegisterUser(context.Background(), RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
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
	if h.service.Logout(context.Background(), req.RefreshToken) != nil {
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
	if h.service.ResetPassword(context.Background(), req.ResetToken, req.NewPassword) != nil {
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
	if h.service.VerifyEmail(r.Context(), req.Token) != nil {
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
