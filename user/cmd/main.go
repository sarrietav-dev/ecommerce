package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/sarrietav-dev/ecommerce/user/internal/config"
	"github.com/sarrietav-dev/ecommerce/user/internal/database"
	"github.com/sarrietav-dev/ecommerce/user/internal/handlers"
	"github.com/sarrietav-dev/ecommerce/user/internal/logger"
	"github.com/sarrietav-dev/ecommerce/user/internal/middleware"
)

func main() {
	logger.InitLogger()

	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		logger.Logger.Error("Failed to connect to the database", slog.String("error", err.Error()))
		log.Fatal(err)
	}

	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	middlewareMux := middleware.RecoveryMiddleware(middleware.LoggingMiddleware(mux))

	port := cfg.Port
	logger.Logger.Info("Server starting", slog.String("port", port))
	if err := http.ListenAndServe(":"+port, middlewareMux); err != nil {
		logger.Logger.Error("Server failed to start", slog.String("error", err.Error()))
		log.Fatal(err)
	}
}
