package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
Config
*/

type Config struct {
	Env string

	Port string

	JWTSecret string

	DatabaseURL string

	UserServiceURL string
	FoodServiceURL string

	LogLevel string
}

/*
Global Config
*/

var AppConfig *Config

/*
Load

Loads application configuration from environment variables.
*/

func Load() {
	_ = godotenv.Load()

	AppConfig = &Config{
		Env: getEnv(
			"ENV",
			"development",
		),

		Port: getEnv(
			"PORT",
			"8083",
		),

		JWTSecret: mustEnv(
			"JWT_SECRET",
		),

		DatabaseURL: mustEnv(
			"DATABASE_URL",
		),

		UserServiceURL: mustEnv(
			"USER_SERVICE_URL",
		),

		FoodServiceURL: mustEnv(
			"FOOD_SERVICE_URL",
		),

		LogLevel: getEnv(
			"LOG_LEVEL",
			"info",
		),
	}
}

/*
Helpers
*/

func getEnv(
	key string,
	defaultValue string,
) string {

	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func mustEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s is required", key)
	}

	return value
}