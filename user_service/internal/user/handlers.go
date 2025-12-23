package user

import (
	"encoding/json"
	"net/http"

	"github.com/adityawaradkar/gratia/user_service/internal/middleware"
)

// Handler handles HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

/* ===================== REQUEST MODELS ===================== */

type CreateDonorProfileRequest struct {
	Name    string  `json:"name"`
	Phone   *string `json:"phone,omitempty"`
	Address *string `json:"address,omitempty"`
}

type UpdateDonorProfileRequest struct {
	Name    string  `json:"name"`
	Phone   *string `json:"phone,omitempty"`
	Address *string `json:"address,omitempty"`
}

type CreateNGOProfileRequest struct {
	Organization   string `json:"organization"`
	RegistrationNo string `json:"registrationNo"`
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

/* ===================== INTERNAL ===================== */

// GetDonorProfileInternal returns donor profile by user ID
func (h *Handler) GetDonorProfileInternal(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user id required")
		return
	}

	profile, err := h.service.GetDonorProfileByUserIDInternal(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "donor profile not found")
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

// GetNGOProfileInternal returns NGO profile by user ID
func (h *Handler) GetNGOProfileInternal(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user id required")
		return
	}

	profile, err := h.service.GetNGOProfileByUserIDInternal(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "ngo profile not found")
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

/* ===================== DONOR PROFILE ===================== */

// CreateDonorProfile creates donor profile
func (h *Handler) CreateDonorProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	role := middleware.UserRole(r.Context())

	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateDonorProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.service.CreateDonorProfile(
		r.Context(),
		userID,
		role,
		req.Name,
		req.Phone,
		req.Address,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, profile)
}

// GetMyDonorProfile returns donor profile of authenticated user
func (h *Handler) GetMyDonorProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.service.GetMyDonorProfile(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

// UpdateMyDonorProfile updates donor profile
func (h *Handler) UpdateMyDonorProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdateDonorProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.UpdateMyDonorProfile(
		r.Context(),
		userID,
		req.Name,
		req.Phone,
		req.Address,
	); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* ===================== NGO PROFILE ===================== */

// CreateNGOProfile creates NGO profile
func (h *Handler) CreateNGOProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	role := middleware.UserRole(r.Context())

	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateNGOProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ngo, err := h.service.CreateNGOProfile(
		r.Context(),
		userID,
		role,
		req.Organization,
		req.RegistrationNo,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, ngo)
}

// GetMyNGOProfile returns NGO profile of authenticated user
func (h *Handler) GetMyNGOProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ngo, err := h.service.GetMyNGOProfile(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, ngo)
}

/* ===================== ADMIN ===================== */

// VerifyNGO verifies NGO profile (admin only)
func (h *Handler) VerifyNGO(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.UserID(r.Context())
	adminRole := middleware.UserRole(r.Context())

	targetUserID := r.URL.Query().Get("userId")
	if targetUserID == "" {
		writeError(w, http.StatusBadRequest, "userId required")
		return
	}

	if err := h.service.VerifyNGO(
		r.Context(),
		adminID,
		adminRole,
		targetUserID,
	); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
