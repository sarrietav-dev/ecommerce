package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

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

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if err := h.AuthService.Register(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.Login(user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
