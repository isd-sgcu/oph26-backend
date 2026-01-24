package config

import (
	"os"
)

type Config struct {
	DataBaseURL string
}

func LoadEnv() *Config {
	return &Config {
		DataBaseURL: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}