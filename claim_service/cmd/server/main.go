package main

import (
	"context"
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

	/*
		Load Configuration
	*/

	config.Load()
	cfg := config.AppConfig

	/*
		Logger
	*/

	logr := logger.New(logger.Config{
		Level: cfg.LogLevel,
	})

	logr.Info(
		"starting claim service",
		"env", cfg.Env,
		"port", cfg.Port,
	)

	/*
		Database
	*/

	dbConn, err := db.New()
	if err != nil {
		logr.Error(
			"failed to connect to database",
			"error", err,
		)
		os.Exit(1)
	}

	defer dbConn.Close()

	/*
		Repository
	*/

	repo := claim.NewClaimRepository(dbConn)

	/*
		External Clients
	*/

	userClient := client.NewUserClient(
		cfg.UserServiceURL,
	)

	foodClient := client.NewFoodClient(
		cfg.FoodServiceURL,
	)

	/*
		Service Layer
	*/

	claimService := claim.NewService(
		dbConn,
		repo,
		userClient,
		foodClient,
	)

	/*
		Handlers
	*/

	claimHandler := claim.NewHandler(
		claimService,
	)

	/*
		Middleware
	*/

	authMiddleware := authmw.NewMiddleware(
		cfg.JWTSecret,
	)

	/*
		Router
	*/

	srv := server.NewServer(
		claimHandler,
		authMiddleware,
	)

	/*
		HTTP Server
	*/

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	/*
		Start Server
	*/

	go func() {

		logr.Info(
			"http server started",
			"port", cfg.Port,
		)

		if err := httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logr.Error(
				"http server failed",
				"error", err,
			)

			os.Exit(1)
		}
	}()

	/*
		Graceful Shutdown
	*/

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	logr.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logr.Error(
			"server shutdown failed",
			"error", err,
		)
		os.Exit(1)
	}

	logr.Info("server exited cleanly")
}