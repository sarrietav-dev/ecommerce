package main

import (
	"log"
	"net/http"

	"github.com/sarrietav-dev/ecommerce/user/internal/config"
	"github.com/sarrietav-dev/ecommerce/user/internal/database"
	"github.com/sarrietav-dev/ecommerce/user/internal/handlers"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)

	http.HandleFunc("POST /register", authHandler.Register)
	http.HandleFunc("POST /login", authHandler.Login)

	port := cfg.Port
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
