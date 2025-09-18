package config

import (
	"log"
	"os"
	"strconv"

	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// AppConfig holds the global configuration instance
var AppConfig *Config

// LoadConfig initializes configuration from .env or environment variables
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	accessMinutes, err := strconv.Atoi(getEnv("ACCESS_TOKEN_MINUTES", "15"))
	if err != nil {
		log.Println("Invalid ACCESS_TOKEN_MINUTES, defaulting to 15")
		accessMinutes = 15
	}

	refreshDays, err := strconv.Atoi(getEnv("REFRESH_TOKEN_DAYS", "30"))
	if err != nil {
		log.Println("Invalid REFRESH_TOKEN_DAYS, defaulting to 30")
		refreshDays = 30
	}

	AppConfig = &Config{
		Port:            port,
		DatabaseURL:     dbURL,
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  time.Duration(accessMinutes) * time.Minute,
		RefreshTokenTTL: time.Duration(refreshDays*24) * time.Hour,
	}
}

// getEnv returns the environment variable value or a default
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

// Helper functions
func GetPort() string {
	return AppConfig.Port
}

func GetDatabaseURL() string {
	return AppConfig.DatabaseURL
}

func GetJWTSecret() string {
	return AppConfig.JWTSecret
}

func GetAccessTokenTTL() time.Duration {
	return AppConfig.AccessTokenTTL
}

func GetRefreshTokenTTL() time.Duration {
	return AppConfig.RefreshTokenTTL
}
