package httpserver

import (
	"net/http"

	"github.com/adityawaradkar/gratia/auth_service/internal/auth"
	"github.com/adityawaradkar/gratia/auth_service/internal/middleware"
)

func RegisterRoutes(handler *auth.Handler) http.Handler {
	mux := http.NewServeMux()

	// Public auth routes
	mux.HandleFunc("/auth/register", handler.RegisterUser)
	mux.HandleFunc("/auth/login", handler.LoginUser)
	mux.HandleFunc("/auth/refresh", handler.RefreshTokens)
	mux.HandleFunc("/auth/logout", handler.Logout)

	// Password recovery
	mux.HandleFunc("/auth/forgot-password", handler.ForgotPassword)
	mux.HandleFunc("/auth/reset-password", handler.ResetPassword)

	// Health
	mux.HandleFunc("/health", handler.HealthCheck)

	// Protected routes
	mux.Handle(
		"/auth/me",
		middleware.Auth(http.HandlerFunc(handler.GetCurrentUser)),
	)

	return mux
}
