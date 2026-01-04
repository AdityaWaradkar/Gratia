package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds claim service configuration
type Config struct {
	Env string

	Port string

	JWTSecret string

	UserServiceURL string
	FoodServiceURL string

	LogLevel string
}

// AppConfig is the loaded configuration
var AppConfig *Config

// Load reads environment variables into Config
func Load() {
	_ = godotenv.Load()

	AppConfig = &Config{
		Env: getEnv("ENV", "development"),

		Port: getEnv("PORT", "8083"),

		JWTSecret: mustEnv("JWT_SECRET"),

		UserServiceURL: mustEnv("USER_SERVICE_URL"),
		FoodServiceURL: mustEnv("FOOD_SERVICE_URL"),

		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

/* ===================== HELPERS ===================== */

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}
	return v
}
