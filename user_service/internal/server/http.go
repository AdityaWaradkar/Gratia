package server

import (
	"database/sql"
	"log"
	"net/http"

	"user_service/internal/config"
)

func Start(cfg *config.Config, db *sql.DB) {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("user service healthy"))
	})

	log.Println("User Service running on port", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, mux))
}
