package config

import (
	"os"
)

type Config struct {
	DataBaseURL    string
	JWTSecret      string
	GoogleClientID string
	AppEnv         string
}

func LoadEnv() *Config {
	return &Config{
		DataBaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		AppEnv:         getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
