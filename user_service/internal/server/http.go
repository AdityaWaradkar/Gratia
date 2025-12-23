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
		"/internal/users/{id}/donor",
		handler.GetDonorProfileInternal,
	)

	mux.HandleFunc(
		"/internal/users/{id}/ngo",
		handler.GetNGOProfileInternal,
	)

	/* ===================== DONOR PROFILE ===================== */

	mux.Handle(
		"/donors/profile",
		middleware.Auth(http.HandlerFunc(handler.CreateDonorProfile)),
	)

	mux.Handle(
		"/donors/profile/me",
		middleware.Auth(http.HandlerFunc(handler.GetMyDonorProfile)),
	)

	mux.Handle(
		"/donors/profile/me/update",
		middleware.Auth(http.HandlerFunc(handler.UpdateMyDonorProfile)),
	)

	/* ===================== NGO PROFILE ===================== */

	mux.Handle(
		"/ngos",
		middleware.Auth(
			http.HandlerFunc(handler.CreateNGOProfile),
		),
	)

	mux.Handle(
		"/ngos/me",
		middleware.Auth(
			http.HandlerFunc(handler.GetMyNGOProfile),
		),
	)

	/* ===================== ADMIN ===================== */

	mux.Handle(
		"/admin/ngos/verify",
		middleware.Auth(
			middleware.RequireRole("ADMIN")(
				http.HandlerFunc(handler.VerifyNGO),
			),
		),
	)

	/* ===================== HEALTH ===================== */

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user service healthy"))
	})

	return mux
}
