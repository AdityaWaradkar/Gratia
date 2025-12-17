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
	UserServiceURL string
}

var AppConfig *Config

func LoadConfig() {
	godotenv.Load()

	port := getEnv("PORT", "8080")

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	secret := getEnv("JWT_SECRET", "")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	userServiceURL := getEnv("USER_SERVICE_URL", "")
	if userServiceURL == "" {
		log.Fatal("USER_SERVICE_URL is required")
	}

	accessMinutes, err := strconv.Atoi(getEnv("ACCESS_TOKEN_MINUTES", "15"))
	if err != nil {
		accessMinutes = 15
	}

	refreshDays, err := strconv.Atoi(getEnv("REFRESH_TOKEN_DAYS", "30"))
	if err != nil {
		refreshDays = 30
	}

	AppConfig = &Config{
		Port:            port,
		DatabaseURL:     databaseURL,
		JWTSecret:       secret,
		AccessTokenTTL:  time.Duration(accessMinutes) * time.Minute,
		RefreshTokenTTL: time.Duration(refreshDays*24) * time.Hour,
		UserServiceURL:  userServiceURL,
	}
}


func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok && value != "" {
		return value
	}
	return defaultValue
}

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
