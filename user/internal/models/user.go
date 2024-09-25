package models

// User represents a user in the system
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ID       uint   `json:"id"`
}
