package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetProductByID(id string) (*models.Product, error) {
	var product models.Product

	err := r.DB.QueryRow("SELECT id, title, description, image, price, created_at FROM products WHERE id = ?", id).Scan(&product.Id, &product.Title, &product.Description, &product.Image, &product.Price, &product.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetProducts(limit uint, offset uint) ([]*models.Product, error) {
	q, args, err := sq.Select("id", "title", "description", "image", "price", "created_at").From("products").Limit(uint64(limit)).Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.Id, &product.Title, &product.Description, &product.Image, &product.Price, &product.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	q, args, err := sq.Insert("products").Columns("id", "title", "description", "image", "price").
		Values(product.Id, product.Title, product.Description, product.Image, product.Price).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(q, args...)
	if err != nil {
		return err
	}
	return nil
}
