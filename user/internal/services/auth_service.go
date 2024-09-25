package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sarrietav-dev/ecommerce/user/internal/models"
	"github.com/sarrietav-dev/ecommerce/user/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo      *repository.UserRepository
	JWTSecret string
}

func NewAuthService(repo *repository.UserRepository, secret string) *AuthService {
	return &AuthService{
		Repo:      repo,
		JWTSecret: secret,
	}
}

func (s *AuthService) Register(user *models.User) error {
	existingUser, _ := s.Repo.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return s.Repo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.Email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
