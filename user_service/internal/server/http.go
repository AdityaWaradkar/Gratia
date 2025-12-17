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
	rootMux := http.NewServeMux()

	// --------------------
	// Health (public)
	// --------------------
	rootMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user service healthy"))
	})

	// --------------------
	// User module wiring
	// --------------------
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// --------------------
	// INTERNAL (NO JWT)
	// --------------------
	rootMux.HandleFunc("/internal/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		userHandler.CreateProfile(w, r)
	})

	// --------------------
	// PROTECTED (JWT)
	// --------------------
	protectedMux := http.NewServeMux()

	protectedMux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetMe(w, r)
		case http.MethodPut:
			userHandler.UpdateMe(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	secured := middleware.AuthMiddleware(cfg)(protectedMux)

	// Mount secured routes under /users
	rootMux.Handle("/users/", secured)

	log.Println("User Service running on port", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, rootMux))
}
