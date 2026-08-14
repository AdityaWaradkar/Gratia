package claim

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/adityawaradkar/gratia/claim_service/internal/middleware"
)

// Handler handles HTTP requests for claim operations
type Handler struct {
	service *Service
}

// NewHandler creates a new claim handler instance
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Request DTOs for claim endpoints
type createClaimRequest struct {
	FoodListingID string `json:"foodListingId"`
}

// CreateClaim handles the creation of a new claim
func (h *Handler) CreateClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRole := middleware.UserRole(ctx)
	if userRole != string(ActorNGO) {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}

	var req createClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FoodListingID == "" {
		writeError(w, http.StatusBadRequest, "foodListingId is required")
		return
	}

	claim, err := h.service.CreateClaim(ctx, req.FoodListingID, middleware.UserID(ctx))
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, claim)
}

// ApproveClaim approves a claim by the donor
func (h *Handler) ApproveClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRole := middleware.UserRole(ctx)
	if userRole != string(ActorDonor) {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}

	claimID := r.PathValue("id")
	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	if err := h.service.ApproveClaim(ctx, claimID, middleware.UserID(ctx)); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RejectClaim rejects a claim by the donor
func (h *Handler) RejectClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRole := middleware.UserRole(ctx)
	if userRole != string(ActorDonor) {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}

	claimID := r.PathValue("id")
	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	if err := h.service.RejectClaim(ctx, claimID, middleware.UserID(ctx)); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CancelClaim cancels a claim by the NGO
func (h *Handler) CancelClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRole := middleware.UserRole(ctx)
	if userRole != string(ActorNGO) {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}

	claimID := r.PathValue("id")
	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	if err := h.service.CancelByNGO(ctx, claimID, middleware.UserID(ctx)); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MarkPickedUp marks a claim as picked up by the NGO
func (h *Handler) MarkPickedUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRole := middleware.UserRole(ctx)
	if userRole != string(ActorNGO) {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}

	claimID := r.PathValue("id")
	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	if err := h.service.MarkPickedUp(ctx, claimID, middleware.UserID(ctx)); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MarkDelivered marks a claim as delivered by the NGO
func (h *Handler) MarkDelivered(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRole := middleware.UserRole(ctx)
	if userRole != string(ActorNGO) {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}

	claimID := r.PathValue("id")
	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	if err := h.service.MarkDelivered(ctx, claimID, middleware.UserID(ctx)); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetClaim retrieves a claim by ID
func (h *Handler) GetClaim(w http.ResponseWriter, r *http.Request) {
	claimID := r.PathValue("id")
	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	claim, err := h.service.GetClaimByID(r.Context(), claimID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, claim)
}

// GetClaimsByNGO retrieves all claims for the authenticated NGO
func (h *Handler) GetClaimsByNGO(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRole := middleware.UserRole(ctx)
	if userRole != string(ActorNGO) {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}

	claims, err := h.service.GetClaimsByNGO(ctx, middleware.UserID(ctx))
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, claims)
}

// GetClaimsByDonor retrieves all claims for the authenticated donor
func (h *Handler) GetClaimsByDonor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userRole := middleware.UserRole(ctx)
	if userRole != string(ActorDonor) {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}

	claims, err := h.service.GetClaimsByDonor(ctx, middleware.UserID(ctx))
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, claims)
}

// GetClaimsByFood retrieves all claims for a specific food listing
func (h *Handler) GetClaimsByFood(w http.ResponseWriter, r *http.Request) {
	foodListingID := r.PathValue("foodId")
	if foodListingID == "" {
		writeError(w, http.StatusBadRequest, "missing food listing id")
		return
	}

	claims, err := h.service.GetClaimsByFood(r.Context(), foodListingID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, claims)
}

// HealthCheck provides a lightweight endpoint for container orchestration readiness probes
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "claim_service"})
}

// handleServiceError maps service errors to appropriate HTTP status codes
func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrInvalidState):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrActiveClaimExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrNGONotVerified):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrFoodNotOpen):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrSelfClaim):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrClaimNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

// Response helpers for JSON responses
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}