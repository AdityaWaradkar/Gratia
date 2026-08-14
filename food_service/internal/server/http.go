package server

import (
	"net/http"

	"github.com/adityawaradkar/gratia/food_service/internal/food"
	"github.com/adityawaradkar/gratia/food_service/internal/middleware"
)

// RegisterRoutes wires all HTTP routes for the food service using the JWT secret for token verification
func RegisterRoutes(handler *food.Handler, jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	// Internal service-to-service endpoints (Used by Claim Service, no external Auth required)
	mux.HandleFunc("GET /internal/foods/{id}/validate", handler.ValidateFoodClaim)
	mux.HandleFunc("PATCH /internal/foods/{id}/claim", handler.MarkFoodClaimed)

	// Food listing endpoints protected by JWT authentication
	mux.Handle("POST /foods", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.CreateFoodListing)))
	mux.Handle("GET /foods", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.ListAvailableFoodListings)))
	mux.Handle("GET /foods/{id}", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.GetFoodListing)))
	mux.Handle("PUT /foods/{id}", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.UpdateFoodListing)))
	mux.Handle("DELETE /foods/{id}", middleware.Auth(jwtSecret)(http.HandlerFunc(handler.CancelFoodListing)))

	// Health check endpoint for Docker/Kubernetes container orchestration
	mux.HandleFunc("GET /health", handler.HealthCheck)

	return mux
}