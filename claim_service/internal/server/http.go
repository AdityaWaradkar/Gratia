package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/adityawaradkar/gratia/claim_service/internal/claim"
	authmw "github.com/adityawaradkar/gratia/claim_service/internal/middleware"
)

/*
Server
*/

type Server struct {
	router *mux.Router
}

func NewServer(
	claimHandler *claim.Handler,
	authMiddleware *authmw.Middleware,
) *Server {

	r := mux.NewRouter()

	// health check
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// protected routes
	claimsRouter := r.PathPrefix("/claims").Subrouter()
	claimsRouter.Use(authMiddleware.RequireAuth)

	claimsRouter.HandleFunc("", claimHandler.CreateClaim).
		Methods(http.MethodPost)

	claimByID := claimsRouter.PathPrefix("/{id}").Subrouter()
	claimByID.Use(claimIDMiddleware)

	claimByID.HandleFunc("/approve", claimHandler.ApproveClaim).
		Methods(http.MethodPost)

	claimByID.HandleFunc("/reject", claimHandler.RejectClaim).
		Methods(http.MethodPost)

	claimByID.HandleFunc("/cancel", claimHandler.CancelClaim).
		Methods(http.MethodPost)

	claimByID.HandleFunc("/pickup", claimHandler.MarkPickedUp).
		Methods(http.MethodPost)

	claimByID.HandleFunc("/deliver", claimHandler.MarkDelivered).
		Methods(http.MethodPost)

	return &Server{router: r}
}

func (s *Server) Handler() http.Handler {
	return s.router
}

/*
Context middleware
*/

type claimIDContextKey struct{}

func claimIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		claimID := vars["id"]

		if claimID == "" {
			http.Error(w, "missing claim id", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), claimIDContextKey{}, claimID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
