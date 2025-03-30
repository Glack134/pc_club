package config

import (
	"os"
	"time"
)

type Config struct {
	ServerPort    string
	JWTSecret     string
	TokenDuration time.Duration
}

func Load() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "50051"),
		JWTSecret:     getEnv("JWT_SECRET", "default_secret_key"),
		TokenDuration: 24 * time.Hour,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
