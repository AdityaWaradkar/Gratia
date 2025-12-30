package config

import "os"

type Config struct {
	Env string

	ServerPort string

	JWTSecret string

	UserServiceURL string
	FoodServiceURL string

	LogLevel string
}

func Load() (*Config, error) {
	cfg := &Config{
		Env: getEnv("ENV", "development"),

		ServerPort: getEnv("SERVER_PORT", "8080"),

		JWTSecret: getEnv("JWT_SECRET", ""),

		UserServiceURL: getEnv("USER_SERVICE_URL", ""),
		FoodServiceURL: getEnv("FOOD_SERVICE_URL", ""),

		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
