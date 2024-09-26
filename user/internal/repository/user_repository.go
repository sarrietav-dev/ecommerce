package repository

import (
	"database/sql"
	"errors"

	"github.com/sarrietav-dev/ecommerce/user/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

var ErrUserAlreadyExists = errors.New("this email already exists")

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := repo.DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) Create(user *models.User) error {
	_, err := repo.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}
