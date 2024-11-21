package models

import "github.com/google/uuid"

type Category struct {
	Id          string `jsonapi:"primary,categories"`
	Name        string `jsonapi:"attr,name"`
	Description string `jsonapi:"attr,description"`
}

func NewCategory(name, description string) *Category {
	return &Category{
		Id: uuid.NewString(),
		Name:        name,
		Description: description,
	}
}