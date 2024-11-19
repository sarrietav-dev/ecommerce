package repository

import (
	"database/sql"

	"github.com/sarrietav-dev/ecommerce/catalog/internal/models"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) GetCategoryByID(id string) (*models.Category, error) {
	var category models.Category

	err := r.DB.QueryRow("SELECT id, name, description FROM categories WHERE id = ?", id).Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetCategories() ([]*models.Category, error) {
	rows, err := r.DB.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}
