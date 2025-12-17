package server

import (
	"database/sql"
	"log"
	"net/http"

	"user_service/internal/config"
	"user_service/internal/middleware"
	"user_service/internal/user"
)

func Start(cfg *config.Config, db *sql.DB) {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user service healthy"))
	})

	// User module wiring
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetMe(w, r)
		case http.MethodPut:
			userHandler.UpdateMe(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Apply auth middleware
	protected := middleware.AuthMiddleware(cfg)(mux)

	log.Println("User Service running on port", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, protected))
}
