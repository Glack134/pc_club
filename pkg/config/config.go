package config

import "os"

type Config struct {
	ServerPort     string
	DBPath         string
	JWTUserSecret  string
	JWTAdminSecret string
}

func Load() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "50051"),
		DBPath:         getEnv("DB_PATH", "./data/sessions.db"),
		JWTUserSecret:  getEnv("JWT_USER_SECRET", "default_user_secret"),
		JWTAdminSecret: getEnv("JWT_ADMIN_SECRET", "default_admin_secret"),
	}
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
