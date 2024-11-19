package config

import "os"

type Config struct {
	Database string
	Port     string
}

func LoadConfig() *Config {
	return &Config{
		Database: getEnv("DATABASE_URI", "postgres://postgres:postgres@localhost:5432/catalog?sslmode=disable"),
		Port:     getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
