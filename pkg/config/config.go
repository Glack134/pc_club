package config

import "os"

type Config struct {
	ServerPort string
	DBPath     string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "50051"),
		DBPath:     getEnv("DB_PATH", "./data/sessions.db"),
		JWTSecret:  getEnv("JWT_SECRET", "default-secret-key-change-me"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
