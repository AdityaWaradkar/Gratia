package server

import (
	"net/http"

	"github.com/adityawaradkar/gratia/food_service/internal/food"
	"github.com/adityawaradkar/gratia/food_service/internal/middleware"
)

// RegisterRoutes wires all HTTP routes
func RegisterRoutes(handler *food.Handler) http.Handler {
	mux := http.NewServeMux()

	/* ===================== FOOD LISTINGS ===================== */

	// Create food listing (DONOR only - validated in service)
	mux.Handle(
		"/foods",
		middleware.Auth(
			http.HandlerFunc(handler.CreateFoodListing),
		),
	)

	// List open food listings (public to authenticated users)
	mux.Handle(
		"/foods/open",
		middleware.Auth(
			http.HandlerFunc(handler.ListOpenFoodListings),
		),
	)

	// Get single food listing
	mux.Handle(
		"/foods/{id}",
		middleware.Auth(
			http.HandlerFunc(handler.GetFoodListing),
		),
	)

	// Update food listing (only owner donor)
	mux.Handle(
		"/foods/{id}/update",
		middleware.Auth(
			http.HandlerFunc(handler.UpdateFoodListing),
		),
	)

	/* ===================== HEALTH ===================== */

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("food service healthy"))
	})

	return mux
}
