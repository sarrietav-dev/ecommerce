package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
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
	q, args, err := sq.Select("id", "name", "description").From("categories").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(q, args...).Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetCategories() ([]*models.Category, error) {
	q, _, err := sq.Select("id", "name", "description").From("categories").ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Query(q)
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

func (r *CategoryRepository) CreateCategory(category *models.Category) (*models.Category, error) {
	q, args, err := sq.Insert("categories").Columns("id", "name", "description").Values(category.Id, category.Name, category.Description).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	_, err = r.DB.Exec(q, args...)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) LinkProductWithCategories(productID string, categoryIds []string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	for _, categoryID := range categoryIds {
		q, args, err := sq.Insert("product_categories").Columns("product_id", "category_id").Values(productID, categoryID).PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return err
		}
		_, err = tx.Exec(q, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}