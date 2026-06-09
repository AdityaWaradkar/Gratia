package server

import (
	"net/http"

	"github.com/adityawaradkar/gratia/food_service/internal/food"
	"github.com/adityawaradkar/gratia/food_service/internal/middleware"
)

// RegisterRoutes wires all HTTP routes
func RegisterRoutes(handler *food.Handler) http.Handler {

	mux := http.NewServeMux()

	/* ===================== INTERNAL (SERVICE-TO-SERVICE) ===================== */

	mux.HandleFunc(
		"GET /internal/foods/{id}/validate",
		handler.ValidateFoodClaim,
	)

	mux.HandleFunc(
		"PATCH /internal/foods/{id}/claim",
		handler.MarkFoodClaimed,
	)

	/* ===================== FOOD LISTINGS ===================== */

	mux.Handle(
		"POST /foods",
		middleware.Auth(
			http.HandlerFunc(handler.CreateFoodListing),
		),
	)

	mux.Handle(
		"GET /foods",
		middleware.Auth(
			http.HandlerFunc(handler.ListAvailableFoodListings),
		),
	)

	mux.Handle(
		"GET /foods/{id}",
		middleware.Auth(
			http.HandlerFunc(handler.GetFoodListing),
		),
	)

	mux.Handle(
		"PUT /foods/{id}",
		middleware.Auth(
			http.HandlerFunc(handler.UpdateFoodListing),
		),
	)

	mux.Handle(
		"DELETE /foods/{id}",
		middleware.Auth(
			http.HandlerFunc(handler.CancelFoodListing),
		),
	)

	/* ===================== HEALTH ===================== */

	mux.HandleFunc(
		"GET /health",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("food service healthy"))
		},
	)

	return mux
}