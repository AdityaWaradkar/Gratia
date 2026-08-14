package food

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/adityawaradkar/gratia/food_service/internal/middleware"
)

// Handler handles HTTP requests for food operations
type Handler struct {
	service *Service
}

// NewHandler creates a new food handler instance
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Request models for food endpoints
type CreateFoodListingRequest struct {
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	ExpiryTime  time.Time `json:"expiryTime"`
	Location    string    `json:"location"`
	ImageURL    *string   `json:"imageUrl,omitempty"`
}

type UpdateFoodListingRequest struct {
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	ExpiryTime  time.Time `json:"expiryTime"`
	Location    string    `json:"location"`
	ImageURL    *string   `json:"imageUrl,omitempty"`
}

type ClaimValidationResponse struct {
	FoodID    string `json:"foodId"`
	Claimable bool   `json:"claimable"`
	Reason    string `json:"reason,omitempty"`
}

// Helper functions for HTTP responses
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}

// handleServiceError maps domain errors to proper HTTP status codes
func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrForbidden):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrListingNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrStatusNotAvailable):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

// CreateFoodListing handles the creation of a new food listing
func (h *Handler) CreateFoodListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())

	if userID == "" {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var req CreateFoodListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	listing, err := h.service.CreateFoodListing(
		r.Context(),
		userID,
		req.Title,
		req.Description,
		req.Quantity,
		req.Unit,
		req.ExpiryTime,
		req.Location,
		req.ImageURL,
	)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, listing)
}

// GetFoodListing retrieves a specific food listing by ID
func (h *Handler) GetFoodListing(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "listing id required")
		return
	}

	listing, err := h.service.GetFoodListing(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, listing)
}

// ListAvailableFoodListings retrieves all available food listings
func (h *Handler) ListAvailableFoodListings(w http.ResponseWriter, r *http.Request) {
	listings, err := h.service.ListAvailableFoodListings(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, listings)
}

// UpdateFoodListing updates an existing food listing
func (h *Handler) UpdateFoodListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	id := r.PathValue("id")

	if userID == "" {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}
	if id == "" {
		writeError(w, http.StatusBadRequest, "listing id required")
		return
	}

	var req UpdateFoodListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	listing := &FoodListing{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		ExpiryTime:  req.ExpiryTime,
		Location:    req.Location,
		ImageURL:    req.ImageURL,
	}

	if err := h.service.UpdateFoodListing(r.Context(), userID, listing); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CancelFoodListing cancels a food listing
func (h *Handler) CancelFoodListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	id := r.PathValue("id")

	if userID == "" {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}
	if id == "" {
		writeError(w, http.StatusBadRequest, "listing id required")
		return
	}

	if err := h.service.CancelFoodListing(r.Context(), userID, id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ValidateFoodClaim validates if a food listing can be claimed (Internal endpoint)
func (h *Handler) ValidateFoodClaim(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "listing id required")
		return
	}

	claimable, reason, err := h.service.IsClaimable(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		ClaimValidationResponse{
			FoodID:    id,
			Claimable: claimable,
			Reason:    reason,
		},
	)
}

// MarkFoodClaimed marks a food listing as claimed (Internal endpoint)
func (h *Handler) MarkFoodClaimed(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "listing id required")
		return
	}

	if err := h.service.MarkClaimed(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck provides a lightweight endpoint for container orchestration readiness probes
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "food_service"})
}