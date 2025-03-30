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
		JWTSecret:  getEnv("JWT_SECRET", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDM0NTE2NDMsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjoiYWRtaW4ifQ.zj_ZukYHGf7v58LtOGx7qlGRhEpepAxrKXc6s4saPFc"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
