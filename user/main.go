package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

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

		_, err := getUserByEmail(db, user.Email)
		if err == nil {
			http.Error(w, formatErrorToJson(ErrUserAlreadyExist), http.StatusConflict)
			return
		}

		encryptedPassword, err := encryptPassword(user.Password)
		if err != nil {
			http.Error(w, formatErrorToJson(err), http.StatusInternalServerError)
			return
		}

		if err := createUser(db, User{
			Name:     user.Name,
			Email:    user.Email,
			Password: encryptedPassword,
		}); err != nil {
			http.Error(w, formatErrorToJson(err), http.StatusInternalServerError)
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
			http.Error(w, formatErrorToJson(err), http.StatusBadRequest)
			return
		}

		existingUser, err := getUserByEmail(db, user.Email)
		if err != nil {
			http.Error(w, formatErrorToJson(ErrUserNotFound), http.StatusNotFound)
		}

		if err := decryptPassword(existingUser.Password, user.Password); err != nil {
			http.Error(w, formatErrorToJson(ErrInvalidPassword), http.StatusUnauthorized)
			return
		}

		token, err := generateToken(existingUser)
		if err != nil {
			http.Error(w, formatErrorToJson(err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token": token,
		})
	})

	log.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
