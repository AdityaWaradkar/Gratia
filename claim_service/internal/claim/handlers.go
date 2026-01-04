package claim

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/adityawaradkar/gratia/claim_service/internal/middleware"
)

/*
Handler
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

// CreateClaim allows NGO to create a claim
func (h *Handler) CreateClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if roleFromContext(ctx) != string(ActorNGO) {
		writeError(w, http.StatusForbidden, "only NGO can create claims")
		return
	}

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

// ApproveClaim allows donor to approve a claim
func (h *Handler) ApproveClaim(w http.ResponseWriter, r *http.Request) {
	if roleFromContext(r.Context()) != string(ActorDonor) {
		writeError(w, http.StatusForbidden, "only donor can approve claims")
		return
	}
	h.handleDonorAction(w, r, h.service.ApproveClaim)
}

// RejectClaim allows donor to reject a claim
func (h *Handler) RejectClaim(w http.ResponseWriter, r *http.Request) {
	if roleFromContext(r.Context()) != string(ActorDonor) {
		writeError(w, http.StatusForbidden, "only donor can reject claims")
		return
	}
	h.handleDonorAction(w, r, h.service.RejectClaim)
}

// CancelClaim allows NGO to cancel a claim
func (h *Handler) CancelClaim(w http.ResponseWriter, r *http.Request) {
	if roleFromContext(r.Context()) != string(ActorNGO) {
		writeError(w, http.StatusForbidden, "only NGO can cancel claims")
		return
	}
	h.handleNGOAction(w, r, h.service.CancelByNGO)
}

// MarkPickedUp allows NGO to mark pickup
func (h *Handler) MarkPickedUp(w http.ResponseWriter, r *http.Request) {
	if roleFromContext(r.Context()) != string(ActorNGO) {
		writeError(w, http.StatusForbidden, "only NGO can mark pickup")
		return
	}
	h.handleNGOAction(w, r, h.service.MarkPickedUp)
}

// MarkDelivered allows NGO to mark delivery
func (h *Handler) MarkDelivered(w http.ResponseWriter, r *http.Request) {
	if roleFromContext(r.Context()) != string(ActorNGO) {
		writeError(w, http.StatusForbidden, "only NGO can mark delivery")
		return
	}
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

func userIDFromContext(ctx context.Context) string {
	if v := ctx.Value(middleware.UserIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

func roleFromContext(ctx context.Context) string {
	if v := ctx.Value(middleware.UserRoleKey); v != nil {
		if role, ok := v.(string); ok {
			return role
		}
	}
	return ""
}

func claimIDFromContext(ctx context.Context) string {
	if v := ctx.Value("claimID"); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
