package claim

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

/*
Handler struct
*/

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

/*
Request DTOs
*/

type createClaimRequest struct {
	FoodListingID string `json:"foodListingId"`
}

/*
Handlers
*/

// CreateClaim handles NGO claim creation
func (h *Handler) CreateClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)

	var req createClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.FoodListingID == "" {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claim, err := h.service.CreateClaim(ctx, req.FoodListingID, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, claim)
}

// ApproveClaim handles donor approval
func (h *Handler) ApproveClaim(w http.ResponseWriter, r *http.Request) {
	h.handleDonorAction(w, r, h.service.ApproveClaim)
}

// RejectClaim handles donor rejection
func (h *Handler) RejectClaim(w http.ResponseWriter, r *http.Request) {
	h.handleDonorAction(w, r, h.service.RejectClaim)
}

// CancelClaim handles NGO cancellation
func (h *Handler) CancelClaim(w http.ResponseWriter, r *http.Request) {
	h.handleNGOAction(w, r, h.service.CancelByNGO)
}

// MarkPickedUp handles NGO pickup confirmation
func (h *Handler) MarkPickedUp(w http.ResponseWriter, r *http.Request) {
	h.handleNGOAction(w, r, h.service.MarkPickedUp)
}

// MarkDelivered handles NGO delivery confirmation
func (h *Handler) MarkDelivered(w http.ResponseWriter, r *http.Request) {
	h.handleNGOAction(w, r, h.service.MarkDelivered)
}

/*
Shared helpers
*/

func (h *Handler) handleDonorAction(
	w http.ResponseWriter,
	r *http.Request,
	fn func(ctx context.Context, claimID, donorUserID string) error,
) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	claimID := claimIDFromContext(ctx)

	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	if err := fn(ctx, claimID, userID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleNGOAction(
	w http.ResponseWriter,
	r *http.Request,
	fn func(ctx context.Context, claimID, ngoUserID string) error,
) {
	ctx := r.Context()
	userID := userIDFromContext(ctx)
	claimID := claimIDFromContext(ctx)

	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	if err := fn(ctx, claimID, userID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
Error mapping
*/

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

	case errors.Is(err, ErrClaimNotFound):
		writeError(w, http.StatusNotFound, err.Error())

	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

/*
Response helpers
*/

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

/*
Context helpers
*/

// userIDFromContext extracts authenticated user id
func userIDFromContext(ctx context.Context) string {
	if v := ctx.Value("userID"); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// claimIDFromContext extracts claim id injected by router
func claimIDFromContext(ctx context.Context) string {
	if v := ctx.Value("claimID"); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
