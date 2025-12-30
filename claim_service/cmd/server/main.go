package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityawaradkar/gratia/claim_service/internal/claim"
	"github.com/adityawaradkar/gratia/claim_service/internal/client"
	"github.com/adityawaradkar/gratia/claim_service/internal/config"
	"github.com/adityawaradkar/gratia/claim_service/internal/db"
	"github.com/adityawaradkar/gratia/claim_service/internal/logger"
	authmw "github.com/adityawaradkar/gratia/claim_service/internal/middleware"
	"github.com/adityawaradkar/gratia/claim_service/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logr := logger.New(logger.Config{
		Level: cfg.LogLevel,
	})

	logr.Info("starting claim service", "env", cfg.Env)

	dbConn, err := db.New()
	if err != nil {
		logr.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	repo := claim.NewClaimRepository(dbConn)

	userClient := client.NewUserClient(cfg.UserServiceURL)
	foodClient := client.NewFoodClient(cfg.FoodServiceURL)

	claimService := claim.NewService(
		dbConn,
		repo,
		userClient,
		foodClient,
	)

	claimHandler := claim.NewHandler(claimService)

	authMiddleware := authmw.NewMiddleware(cfg.JWTSecret)

	srv := server.NewServer(
		claimHandler,
		authMiddleware,
	)

	httpServer := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: srv.Handler(),
	}

	go func() {
		logr.Info("http server started", "port", cfg.ServerPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logr.Error("http server failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logr.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logr.Error("server shutdown failed", "err", err)
	}

	logr.Info("server exited cleanly")
}
