package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/logger"
)

func Connect(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		logger.Logger.Error("Error connecting to database")
		return nil, err
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS products (
      id UUID PRIMARY KEY,
      title VARCHAR(255) NOT NULL,
      description TEXT,
      price NUMERIC NOT NULL,
      image VARCHAR(255),
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
    `)
	if err != nil {
		logger.Logger.Error("Error creating products table")
		return nil, err
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS categories (
      id UUID PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      description TEXT
    )
  `)
	if err != nil {
		logger.Logger.Error("Error creating category table")
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS product_categories (
		product_id UUID,
		category_id UUID,
		PRIMARY KEY (product_id, category_id),
		FOREIGN KEY (product_id) REFERENCES products(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	)
	`)
	if err != nil {
		logger.Logger.Error("Error creating product_categories table")
		return nil, err
	}

	logger.Logger.Info("Database connected and initialized")

	return db, nil
}
