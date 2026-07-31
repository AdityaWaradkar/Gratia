package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/adityawaradkar/gratia/claim_service/internal/claim"
	authmw "github.com/adityawaradkar/gratia/claim_service/internal/middleware"
)

// Server holds the HTTP router configuration
type Server struct {
	router *mux.Router
}

// NewServer creates a new server with all routes configured
func NewServer(
	claimHandler *claim.Handler,
	authMiddleware *authmw.Middleware,
) *Server {

	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Protected routes requiring authentication
	claimsRouter := r.PathPrefix("/claims").Subrouter()
	claimsRouter.Use(authMiddleware.RequireAuth)

	// Claim creation endpoint
	claimsRouter.HandleFunc(
		"",
		claimHandler.CreateClaim,
	).Methods(http.MethodPost)

	// Claim action endpoints
	claimsRouter.HandleFunc(
		"/{id}/approve",
		claimHandler.ApproveClaim,
	).Methods(http.MethodPost)

	claimsRouter.HandleFunc(
		"/{id}/reject",
		claimHandler.RejectClaim,
	).Methods(http.MethodPost)

	claimsRouter.HandleFunc(
		"/{id}/cancel",
		claimHandler.CancelClaim,
	).Methods(http.MethodPost)

	claimsRouter.HandleFunc(
		"/{id}/pickup",
		claimHandler.MarkPickedUp,
	).Methods(http.MethodPost)

	claimsRouter.HandleFunc(
		"/{id}/deliver",
		claimHandler.MarkDelivered,
	).Methods(http.MethodPost)

	return &Server{
		router: r,
	}
}

// Handler returns the HTTP handler for the server
func (s *Server) Handler() http.Handler {
	return s.router
}