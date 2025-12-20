package server

import (
	"net/http"

	"github.com/adityawaradkar/gratia/user_service/internal/middleware"
	"github.com/adityawaradkar/gratia/user_service/internal/user"
)

// RegisterRoutes wires all HTTP routes
func RegisterRoutes(handler *user.Handler) http.Handler {
	mux := http.NewServeMux()

	/* ===================== DONOR PROFILE ===================== */

	mux.Handle(
		"/donors/profile",
		http.HandlerFunc(handler.CreateDonorProfile),
	)

	mux.Handle(
		"/donors/profile/me",
		http.HandlerFunc(handler.GetMyDonorProfile),
	)

	mux.Handle(
		"/donors/profile/me/update",
		http.HandlerFunc(handler.UpdateMyDonorProfile),
	)

	/* ===================== NGO PROFILE ===================== */

	mux.Handle(
		"/ngos",
		middleware.RequireRole("NGO")(http.HandlerFunc(handler.CreateNGOProfile)),
	)

	mux.Handle(
		"/ngos/me",
		middleware.RequireRole("NGO")(http.HandlerFunc(handler.GetMyNGOProfile)),
	)

	/* ===================== ADMIN ===================== */

	mux.Handle(
		"/admin/ngos/verify",
		middleware.RequireRole("ADMIN")(http.HandlerFunc(handler.VerifyNGO)),
	)

	/* ===================== HEALTH ===================== */

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user service healthy"))
	})

	return mux
}
