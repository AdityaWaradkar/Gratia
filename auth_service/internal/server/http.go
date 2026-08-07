package server

import (
    "net/http"

    "github.com/adityawaradkar/gratia/auth_service/internal/auth"
    "github.com/adityawaradkar/gratia/auth_service/internal/middleware"
    "github.com/gorilla/mux"
)

// RegisterRoutes constructs the full routing tree using strict HTTP methods
func RegisterRoutes(handler *auth.Handler, jwtSecret string) http.Handler {
    router := mux.NewRouter()

    // Public authentication endpoints for user onboarding and session management
    router.HandleFunc("/register", handler.RegisterUser).Methods(http.MethodPost)
    router.HandleFunc("/login", handler.LoginUser).Methods(http.MethodPost)
    router.HandleFunc("/refresh", handler.RefreshTokens).Methods(http.MethodPost)
    router.HandleFunc("/logout", handler.Logout).Methods(http.MethodPost)
    
    // Account recovery endpoints
    router.HandleFunc("/forgot-password", handler.ForgotPassword).Methods(http.MethodPost)
    router.HandleFunc("/reset-password", handler.ResetPassword).Methods(http.MethodPost)
    
    // Unauthenticated infrastructure health check for Docker/Kubernetes
    router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

    // Protected endpoints requiring a valid JSON Web Token
    protected := router.PathPrefix("/").Subrouter()
    protected.Use(middleware.Auth(jwtSecret))
    protected.HandleFunc("/me", handler.GetCurrentUser).Methods(http.MethodGet)

    return router
}