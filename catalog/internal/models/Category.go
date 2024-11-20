package models

type Category struct {
	Id          string `jsonapi:"primary,categories"`
	Name        string `jsonapi:"attr,name"`
	Description string `jsonapi:"attr,description"`
}
