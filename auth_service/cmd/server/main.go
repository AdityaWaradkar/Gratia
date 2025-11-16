package main

import (
	"context"
	"log"
	"net/http"

	"github.com/adityawaradkar/gratia/auth_service/internal/auth"
	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	httpServer "github.com/adityawaradkar/gratia/auth_service/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config.LoadConfig()

	dbPool, err := pgxpool.New(context.Background(), config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	repository := auth.NewRepository(dbPool)
	service := auth.NewService(repository, config.GetJWTSecret(), config.GetAccessTokenTTL(), config.GetRefreshTokenTTL())
	handler := auth.NewHandler(service)
	router := httpServer.RegisterRoutes(handler)

	address := ":" + config.GetPort()
	log.Printf("auth service running on %s", address)

	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
