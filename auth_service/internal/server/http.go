package http

import (
	"net/http"

	"github.com/adityawaradkar/gratia/auth_service/internal/auth"
	"github.com/adityawaradkar/gratia/auth_service/internal/middleware" // assuming your middleware package
)

// RegisterRoutes sets up all routes for the auth service
func RegisterRoutes(handler *auth.Handler) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/auth/register", handler.RegisterUser)
	mux.HandleFunc("/auth/login", handler.LoginUser)
	mux.HandleFunc("/auth/refresh", handler.RefreshTokens)
	mux.HandleFunc("/auth/logout", handler.Logout)

	// Protected route: Get current user profile
	mux.Handle("/auth/me", middleware.AuthMiddleware(http.HandlerFunc(handler.GetCurrentUser)))

	return mux
}
