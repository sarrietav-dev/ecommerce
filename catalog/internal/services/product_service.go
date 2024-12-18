package services

import (
	"github.com/sarrietav-dev/ecommerce/catalog/internal/models"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/repository"
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
}

func NewProductService(productRepo *repository.ProductRepository, categoryRepo *repository.CategoryRepository) *ProductService {
	return &ProductService{productRepo: productRepo, categoryRepo: categoryRepo}
}

func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	product.Categories, err = s.categoryRepo.GetCategoriesByProductID(product.Id)

	return product, err
}

func (s *ProductService) GetProducts(limit uint, offset uint) ([]*models.Product, error) {
	return s.productRepo.GetProducts(limit, offset)
}

func (s *ProductService) CreateProduct(product *models.Product) (*models.Product, error) {
	product, err := s.productRepo.CreateProduct(product)

	if err != nil {
		return nil, err
	}

	if len(product.Categories) > 0 {
		categoryIds := make([]string, 0, len(product.Categories))
		for _, category := range product.Categories {
			categoryIds = append(categoryIds, category.Id)
		}

		err = s.categoryRepo.LinkProductWithCategories(product.Id, categoryIds)
		if err != nil {
			return nil, err
		}
	}

	product.Categories, err = s.categoryRepo.GetCategoriesByProductID(product.Id)

	return product, err
}

func (s *ProductService) UpdateProduct() {
	// Update a product
}
func (s *ProductService) DeleteProduct() {
	// Delete a product
}
