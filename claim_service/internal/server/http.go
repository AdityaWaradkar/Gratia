package server

import (
	"net/http"

	"github.com/adityawaradkar/gratia/claim_service/internal/claim"
	"github.com/adityawaradkar/gratia/claim_service/internal/middleware"
)

// RegisterRoutes wires all HTTP routes for the claim service using the JWT secret for token verification
func RegisterRoutes(handler *claim.Handler, jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", handler.HealthCheck)

	// Protected routes requiring authentication
	// Claim creation
	mux.Handle("POST /claims", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.CreateClaim)))

	// Claim retrieval
	mux.Handle("GET /claims/{id}", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.GetClaim)))

	// Claim actions (approve, reject, cancel, pickup, deliver)
	mux.Handle("POST /claims/{id}/approve", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.ApproveClaim)))
	mux.Handle("POST /claims/{id}/reject", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.RejectClaim)))
	mux.Handle("POST /claims/{id}/cancel", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.CancelClaim)))
	mux.Handle("POST /claims/{id}/pickup", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.MarkPickedUp)))
	mux.Handle("POST /claims/{id}/deliver", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.MarkDelivered)))

	// Get claims by user role (NGO or Donor)
	mux.Handle("GET /claims/ngo", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.GetClaimsByNGO)))
	mux.Handle("GET /claims/donor", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.GetClaimsByDonor)))

	// Get claims by food listing
	mux.Handle("GET /foods/{foodId}/claims", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.GetClaimsByFood)))

	return mux
}