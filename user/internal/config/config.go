package config

import "os"

type Config struct {
	Database  string
	JWTSecret string
	Port      string
}

func LoadConfig() *Config {
	cfg := &Config{
		Database:  getEnv("DATABASE", "./users.db"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
		Port:      getEnv("PORT", "8080"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
