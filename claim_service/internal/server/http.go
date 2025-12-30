package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/adityawaradkar/gratia/claim_service/internal/claim"
	authmw "github.com/adityawaradkar/gratia/claim_service/internal/middleware"
)

type Server struct {
	router http.Handler
}

func NewServer(
	claimHandler *claim.Handler,
	authMiddleware *authmw.Middleware,
) *Server {

	r := chi.NewRouter()

	// global middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// health
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// protected routes
	r.Route("/claims", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)

		r.Post("/", claimHandler.CreateClaim)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(claimIDMiddleware)

			r.Post("/approve", claimHandler.ApproveClaim)
			r.Post("/reject", claimHandler.RejectClaim)
			r.Post("/cancel", claimHandler.CancelClaim)
			r.Post("/pickup", claimHandler.MarkPickedUp)
			r.Post("/deliver", claimHandler.MarkDelivered)
		})
	})

	return &Server{router: r}
}

func (s *Server) Handler() http.Handler {
	return s.router
}

/*
Context middleware
*/

func claimIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claimID := chi.URLParam(r, "id")

		ctx := context.WithValue(r.Context(), "claimID", claimID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
