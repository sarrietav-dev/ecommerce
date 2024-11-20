package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/sarrietav-dev/ecommerce/catalog/internal/config"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/database"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/handlers"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/logger"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/middleware"
)

func main() {
	logger.InitLogger()

	config := config.LoadConfig()

	db, err := database.Connect(config.Database)
	if err != nil {
		logger.Logger.Error("Failed to connect to the database", slog.String("error", err.Error()))
		log.Fatal(err)
	}

	productHandler := handlers.NewProductHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products/{id}", productHandler.Show)
	mux.HandleFunc("GET /products", productHandler.Index)
	mux.HandleFunc("POST /products", productHandler.Create)

	middlewareMux := middleware.RecoveryMiddleware(middleware.LoggingMiddleware(mux))

	port := config.Port
	logger.Logger.Info("Server starting", slog.String("port", port))
	if err := http.ListenAndServe(":"+port, middlewareMux); err != nil {
		logger.Logger.Error("Server failed to start", slog.String("error", err.Error()))
		log.Fatal(err)
	}
}
