package claim

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/adityawaradkar/gratia/claim_service/internal/middleware"
)

/*
Handler
*/

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

/*
Request DTOs
*/

type createClaimRequest struct {
	FoodListingID string `json:"foodListingId"`
}

/*
Create Claim
*/

func (h *Handler) CreateClaim(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	if !requireRole(w, ctx, ActorNGO) {
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

	claim, err := h.service.CreateClaim(
		ctx,
		req.FoodListingID,
		userIDFromContext(ctx),
	)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, claim)
}

/*
Donor Actions
*/

func (h *Handler) ApproveClaim(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !requireRole(w, r.Context(), ActorDonor) {
		return
	}

	h.handleDonorAction(
		w,
		r,
		h.service.ApproveClaim,
	)
}

func (h *Handler) RejectClaim(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !requireRole(w, r.Context(), ActorDonor) {
		return
	}

	h.handleDonorAction(
		w,
		r,
		h.service.RejectClaim,
	)
}

/*
NGO Actions
*/

func (h *Handler) CancelClaim(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !requireRole(w, r.Context(), ActorNGO) {
		return
	}

	h.handleNGOAction(
		w,
		r,
		h.service.CancelByNGO,
	)
}

func (h *Handler) MarkPickedUp(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !requireRole(w, r.Context(), ActorNGO) {
		return
	}

	h.handleNGOAction(
		w,
		r,
		h.service.MarkPickedUp,
	)
}

func (h *Handler) MarkDelivered(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !requireRole(w, r.Context(), ActorNGO) {
		return
	}

	h.handleNGOAction(
		w,
		r,
		h.service.MarkDelivered,
	)
}

/*
Shared Action Helpers
*/

func (h *Handler) handleDonorAction(
	w http.ResponseWriter,
	r *http.Request,
	fn func(context.Context, string, string) error,
) {
	ctx := r.Context()

	vars := mux.Vars(r)
	claimID := vars["id"]
	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	if err := fn(
		ctx,
		claimID,
		userIDFromContext(ctx),
	); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleNGOAction(
	w http.ResponseWriter,
	r *http.Request,
	fn func(context.Context, string, string) error,
) {
	ctx := r.Context()

	vars := mux.Vars(r)
	claimID := vars["id"]
	if claimID == "" {
		writeError(w, http.StatusBadRequest, "missing claim id")
		return
	}

	if err := fn(
		ctx,
		claimID,
		userIDFromContext(ctx),
	); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
Role Helper
*/

func requireRole(
	w http.ResponseWriter,
	ctx context.Context,
	role ActorRole,
) bool {

	if roleFromContext(ctx) != string(role) {
		writeError(
			w,
			http.StatusForbidden,
			"insufficient permissions",
		)
		return false
	}

	return true
}

/*
Error Mapping
*/

func handleServiceError(
	w http.ResponseWriter,
	err error,
) {
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
		writeError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}

/*
Response Helpers
*/

func writeJSON(
	w http.ResponseWriter,
	status int,
	v any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(v)
}

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

/*
Context Helpers
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

