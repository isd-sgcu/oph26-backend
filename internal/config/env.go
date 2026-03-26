package config

import (
	"os"
)

type Config struct {
	DataBaseURL          string
	JWTSecret            string
	GoogleClientID       string
	AppEnv               string
	Port                 string
	AllowOrigins         string
	MetricsBasicAuthUser string
	MetricsBasicAuthPass string
}

func LoadEnv() *Config {
	return &Config{
		AppEnv:       getEnv("APP_ENV", "development"),
		Port:         getEnv("PORT", "8080"),
		AllowOrigins: getEnv("ALLOW_ORIGINS", "http://localhost:3000"),

		DataBaseURL:          getEnv("DATABASE_URL", ""),
		JWTSecret:            getEnv("JWT_SECRET", "secret"),
		GoogleClientID:       getEnv("GOOGLE_CLIENT_ID", ""),
		MetricsBasicAuthUser: getEnv("METRICS_BASIC_AUTH_USER", "metrics"),
		MetricsBasicAuthPass: getEnv("METRICS_BASIC_AUTH_PASS", "metrics"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
