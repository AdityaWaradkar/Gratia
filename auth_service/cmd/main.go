// @title Gratia auth_service API
// @version 1.0
// @description Authentication and Authorization Service for Gratia Platform
// @host localhost:8080
// @BasePath /

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	"github.com/adityawaradkar/gratia/auth_service/internal/database"
	"github.com/adityawaradkar/gratia/auth_service/internal/handler"
	"github.com/adityawaradkar/gratia/auth_service/internal/logger"
	"github.com/adityawaradkar/gratia/auth_service/internal/middleware"
	"github.com/adityawaradkar/gratia/auth_service/internal/repository"
	"github.com/adityawaradkar/gratia/auth_service/internal/service"

	_ "github.com/adityawaradkar/gratia/auth_service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	// Initialize Logger
	logger.InitLogger()
	defer logger.SyncLogger()

	logger.Log.Info("Starting auth_service")

	// Load Configuration
	cfg := config.LoadConfig()

	logger.Log.Info(
		"Configuration loaded successfully",
	)

	// Initialize Database
	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	logger.Log.Info(
		"Database connection established",
	)

	// Initialize Repository Layer
	userRepository := repository.NewUserRepository(db)

	// Initialize Service Layer
	authService := service.NewAuthService(
		userRepository,
		cfg,
	)

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	healthHandler := handler.NewHealthHandler()

	// Initialize Router
	router := gin.Default()

	//Swagger Route
	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	// Health Routes
	router.GET(
		"/health",
		healthHandler.HealthCheck,
	)

	// Public Routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST(
			"/signup",
			authHandler.Signup,
		)

		authRoutes.POST(
			"/login",
			authHandler.Login,
		)

		authRoutes.GET(
			"/validate",
			authHandler.Validate,
		)
	}

	// Protected Routes
	protectedRoutes := router.Group("/protected")

	protectedRoutes.Use(
		middleware.AuthMiddleware(authService),
	)

	{
		protectedRoutes.GET(
			"/donor",
			middleware.RBACMiddleware("DONOR"),
			func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "DONOR access granted",
				})
			},
		)

		protectedRoutes.GET(
			"/ngo",
			middleware.RBACMiddleware("NGO"),
			func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "NGO access granted",
				})
			},
		)

		protectedRoutes.GET(
			"/admin",
			middleware.RBACMiddleware("ADMIN"),
			func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "ADMIN access granted",
				})
			},
		)
	}

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {

		logger.Log.Info(
			"auth_service started",
			zap.String("port", cfg.AppPort),
		)

		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {

			logger.Log.Fatal(
				"Failed to start server",
				zap.Error(err),
			)
		}
	}()

	// Graceful Shutdown Signals
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	logger.Log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {

		logger.Log.Fatal(
			"Server forced to shutdown",
			zap.Error(err),
		)
	}

	logger.Log.Info("auth_service stopped gracefully")
}