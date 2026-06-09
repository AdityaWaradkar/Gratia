package server

import (
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

	/*
		Health Check
	*/

	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	/*
		Protected Routes
	*/

	claimsRouter := r.PathPrefix("/claims").Subrouter()
	claimsRouter.Use(authMiddleware.RequireAuth)

	/*
		Claim Creation
	*/

	claimsRouter.HandleFunc(
		"",
		claimHandler.CreateClaim,
	).Methods(http.MethodPost)

	/*
		Claim Actions
	*/

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

func (s *Server) Handler() http.Handler {
	return s.router
}