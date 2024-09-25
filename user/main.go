package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyExist = errors.New("this email already exists")

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ID       uint   `json:"id"`
}

func main() {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT NOT NULL,
      email TEXT NOT NULL,
      password TEXT NOT NULL
    );
    `)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var existingId string
		err := db.QueryRow("SELECT id FROM users WHERE email = :email", sql.Named("email", user.Email)).Scan(&existingId)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if existingId != "" {
			http.Error(w, ErrUserAlreadyExist.Error(), http.StatusBadRequest)
			return
		}

		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		insertStmt, err := db.Prepare(`INSERT INTO users (name, email, password) VALUES (:name, :email, :password);`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = insertStmt.Exec(sql.Named("name", user.Name), sql.Named("email", user.Email), sql.Named("password", encryptedPassword))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User created successfully",
		})
	})

	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var existingUser User
		err := db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", user.Email).Scan(&existingUser.ID, &existingUser.Email, &existingUser.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   user.Email,
			Issuer:    "auth.service",
			Audience:  jwt.ClaimStrings{"user"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:        strconv.FormatUint(uint64(existingUser.ID), 10),
		})

		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token": tokenString,
		})
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
