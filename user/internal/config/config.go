package config

import (
	"os"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Database  string
	JWTSecret string
	Port      string
}

func LoadConfig() *Config {
	dbCfg := mysql.Config{
		Net:    "tcp",
		User:   getEnv("DB_USER", "root"),
		Passwd: getEnv("DB_PASS", "password"),
		Addr:   getEnv("DB_ADDR", "localhost:3306"),
		DBName: getEnv("DB_NAME", "ecommerce"),
	}

	cfg := &Config{
		Database:  dbCfg.FormatDSN(),
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
