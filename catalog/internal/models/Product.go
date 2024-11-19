package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
}

func NewProduct(title, description, image string, price float64) *Product {
	return &Product{
		Id:          uuid.New().String(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       title,
		Description: description,
		Image:       image,
		Price:       price,
	}
}
