package food

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adityawaradkar/gratia/food_service/internal/middleware"
)

/* ===================== HANDLER ===================== */

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

/* ===================== REQUEST MODELS ===================== */

type CreateFoodListingRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Quantity    int        `json:"quantity"`
	Unit        string     `json:"unit"`
	ExpiryTime  time.Time  `json:"expiryTime"`
	Location    string     `json:"location"`
	ImageURL    *string    `json:"imageUrl,omitempty"`
}

type UpdateFoodListingRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Quantity    int        `json:"quantity"`
	Unit        string     `json:"unit"`
	ExpiryTime  time.Time  `json:"expiryTime"`
	Location    string     `json:"location"`
	ImageURL    *string    `json:"imageUrl,omitempty"`
}

/* ===================== HELPERS ===================== */

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"message": msg})
}

/* ===================== CREATE ===================== */

// POST /food
func (h *Handler) CreateFoodListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
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
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, listing)
}

/* ===================== READ ===================== */

// GET /food/{id}
func (h *Handler) GetFoodListing(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "listing id required")
		return
	}

	listing, err := h.service.GetFoodListing(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "food listing not found")
		return
	}

	writeJSON(w, http.StatusOK, listing)
}

// GET /food
func (h *Handler) ListOpenFoodListings(w http.ResponseWriter, r *http.Request) {
	listings, err := h.service.ListOpenFoodListings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch listings")
		return
	}

	writeJSON(w, http.StatusOK, listings)
}

/* ===================== UPDATE ===================== */

// PUT /food/{id}
func (h *Handler) UpdateFoodListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.PathValue("id")
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
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
