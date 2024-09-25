package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/sarrietav-dev/ecommerce/user/internal/logger"
	"github.com/sarrietav-dev/ecommerce/user/internal/models"
	"github.com/sarrietav-dev/ecommerce/user/internal/repository"
	"github.com/sarrietav-dev/ecommerce/user/internal/services"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(db *sql.DB, secret string) *AuthHandler {
	repo := repository.NewUserRepository(db)
	service := services.NewAuthService(repo, secret)
	return &AuthHandler{AuthService: service}
}

func writeErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	response := map[string]interface{}{
		"message": err.Error(),
		"status":  statusCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Logger.Warn("Invalid request payload", slog.String("error", err.Error()))
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err := h.AuthService.Register(&user)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExist) {
			logger.Logger.Info("User registration failed: email already exists", slog.String("email", user.Email))
			writeErrorResponse(w, err, http.StatusConflict)
		} else {
			logger.Logger.Error("User registration failed", slog.String("error", err.Error()))
			writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	logger.Logger.Info("User registered successfully", slog.String("email", user.Email))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Logger.Warn("Invalid login payload", slog.String("error", err.Error()))
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.Login(user.Email, user.Password)
	if err != nil {
		logger.Logger.Error("Login failed", slog.String("email", user.Email), slog.String("error", err.Error()))
		writeErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	logger.Logger.Info("User logged in successfully", slog.String("email", user.Email))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
