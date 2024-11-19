package services

import (
	"github.com/sarrietav-dev/ecommerce/catalog/internal/models"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
	return s.productRepo.GetProductByID(id)
}

func (s *ProductService) GetProducts(limit uint, offset uint) ([]*models.Product, error) {
	return s.productRepo.GetProducts(limit, offset)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.productRepo.CreateProduct(product)
}

func (s *ProductService) UpdateProduct() {
	// Update a product
}

func (s *ProductService) DeleteProduct() {
	// Delete a product
}
