package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/adityawaradkar/gratia/user_service/internal/middleware"
)

// Handler handles HTTP requests for user profile operations
type Handler struct {
	service *Service
}

// NewHandler creates a new user handler instance
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Request models for user endpoints
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

// Helper functions for HTTP responses
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrForbidden):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrDonorProfileExists),
		errors.Is(err, ErrNGOProfileExists),
		errors.Is(err, ErrNGOAlreadyVerified):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrDonorProfileNotFound),
		errors.Is(err, ErrNGOProfileNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrInvalidUserID),
		errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrInvalidNGODetails),
		errors.Is(err, ErrTargetUserIDRequired):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

// GetDonorProfileInternal returns a donor profile by user ID for internal service calls
func (h *Handler) GetDonorProfileInternal(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	profile, err := h.service.GetDonorProfileByUserIDInternal(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

// GetNGOProfileInternal returns an NGO profile by user ID for internal service calls
func (h *Handler) GetNGOProfileInternal(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	profile, err := h.service.GetNGOProfileByUserIDInternal(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

// CreateDonorProfile registers a donor profile for the authenticated user
func (h *Handler) CreateDonorProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	role := middleware.UserRole(r.Context())

	if userID == "" || role == "" {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req CreateDonorProfileRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.service.CreateDonorProfile(r.Context(), userID, role, req.Name, req.Phone, req.Address)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, profile)
}

// GetMyDonorProfile returns the donor profile belonging to the authenticated user
func (h *Handler) GetMyDonorProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())

	if userID == "" {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	profile, err := h.service.GetMyDonorProfile(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

// UpdateMyDonorProfile updates the donor profile of the authenticated user
func (h *Handler) UpdateMyDonorProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())

	if userID == "" {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req UpdateDonorProfileRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.UpdateMyDonorProfile(r.Context(), userID, req.Name, req.Phone, req.Address)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CreateNGOProfile registers an NGO profile for the authenticated user
func (h *Handler) CreateNGOProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	role := middleware.UserRole(r.Context())

	if userID == "" || role == "" {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req CreateNGOProfileRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ngo, err := h.service.CreateNGOProfile(r.Context(), userID, role, req.Organization, req.RegistrationNo)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, ngo)
}

// GetMyNGOProfile returns the NGO profile belonging to the authenticated user
func (h *Handler) GetMyNGOProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())

	if userID == "" {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	ngo, err := h.service.GetMyNGOProfile(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ngo)
}

// VerifyNGO executes administrative verification of an NGO profile
func (h *Handler) VerifyNGO(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.UserID(r.Context())
	adminRole := middleware.UserRole(r.Context())

	if adminID == "" || adminRole == "" {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	targetUserID := r.URL.Query().Get("userId")
	if targetUserID == "" {
		writeError(w, http.StatusBadRequest, "userId query parameter is required")
		return
	}

	err := h.service.VerifyNGO(r.Context(), adminID, adminRole, targetUserID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck provides a lightweight endpoint for container orchestration readiness probes
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "user_service"})
}