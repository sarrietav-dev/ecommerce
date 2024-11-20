package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	CreatedAt   time.Time `json:"created_at" jsonapi:"attr,created_at"`
	UpdatedAt   time.Time `json:"updated_at" jsonapi:"attr,updated_at"`
	Id          string    `json:"id" jsonapi:"primary,products"`
	Title       string    `json:"title" jsonapi:"attr,title"`
	Description string    `json:"description" jsonapi:"attr,description"`
	Image       string    `json:"image" jsonapi:"attr,image"`
	Price       float64   `json:"price" jsonapi:"attr,price"`
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
