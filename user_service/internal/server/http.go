package server

import (
	"net/http"

	"github.com/adityawaradkar/gratia/user_service/internal/middleware"
	"github.com/adityawaradkar/gratia/user_service/internal/user"
)

// RegisterRoutes wires all HTTP routes
func RegisterRoutes(handler *user.Handler) http.Handler {

	mux := http.NewServeMux()

	/* ===================== INTERNAL (SERVICE-TO-SERVICE) ===================== */

	mux.HandleFunc(
		"GET /internal/users/{id}/donor",
		handler.GetDonorProfileInternal,
	)

	mux.HandleFunc(
		"GET /internal/users/{id}/ngo",
		handler.GetNGOProfileInternal,
	)

	/* ===================== DONOR PROFILE ===================== */

	mux.Handle(
		"POST /donors/profile",
		middleware.Auth(
			http.HandlerFunc(handler.CreateDonorProfile),
		),
	)

	mux.Handle(
		"GET /donors/profile/me",
		middleware.Auth(
			http.HandlerFunc(handler.GetMyDonorProfile),
		),
	)

	mux.Handle(
		"PUT /donors/profile/me",
		middleware.Auth(
			http.HandlerFunc(handler.UpdateMyDonorProfile),
		),
	)

	/* ===================== NGO PROFILE ===================== */

	mux.Handle(
		"POST /ngos/profile",
		middleware.Auth(
			http.HandlerFunc(handler.CreateNGOProfile),
		),
	)

	mux.Handle(
		"GET /ngos/profile/me",
		middleware.Auth(
			http.HandlerFunc(handler.GetMyNGOProfile),
		),
	)

	/* ===================== ADMIN ===================== */

	mux.Handle(
		"PUT /admin/ngos/verify",
		middleware.Auth(
			middleware.RequireRole(user.RoleAdmin)(
				http.HandlerFunc(handler.VerifyNGO),
			),
		),
	)

	/* ===================== HEALTH ===================== */

	mux.HandleFunc(
		"GET /health",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("user service healthy"))
		},
	)

	return mux
}