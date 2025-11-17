package httpserver

import (
	"net/http"

	"github.com/adityawaradkar/gratia/auth_service/internal/auth"
	"github.com/adityawaradkar/gratia/auth_service/internal/middleware"
)

func RegisterRoutes(handler *auth.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", handler.RegisterUser)
	mux.HandleFunc("/auth/login", handler.LoginUser)
	mux.HandleFunc("/auth/refresh", handler.RefreshTokens)
	mux.HandleFunc("/auth/logout", handler.Logout)
	mux.HandleFunc("/auth/forgot-password", handler.ForgotPassword)
	mux.HandleFunc("/auth/reset-password", handler.ResetPassword)
	mux.HandleFunc("/auth/verify-email", handler.VerifyEmail)
	mux.HandleFunc("/auth/validate", handler.ValidateTokenHandler)
	mux.HandleFunc("/auth/generate-email-verification", handler.GenerateEmailVerification)
	mux.HandleFunc("/auth/resend-email-verification", handler.ResendEmailVerification)
	mux.HandleFunc("/health", handler.HealthCheck)


	//protected routes
	mux.Handle("/auth/me", middleware.Auth(http.HandlerFunc(handler.GetCurrentUser)))
	mux.Handle("/auth/sessions", middleware.Auth(http.HandlerFunc(handler.GetSessions)))
	mux.Handle("/auth/sessions/", middleware.Auth(http.HandlerFunc(handler.DeleteSession)))

	return mux
}
