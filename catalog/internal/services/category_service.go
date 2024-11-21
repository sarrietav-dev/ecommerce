package services

import (
	"github.com/sarrietav-dev/ecommerce/catalog/internal/models"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) GetCategoryByID(id string) (*models.Category, error) {
	return s.categoryRepo.GetCategoryByID(id)
}

func (s *CategoryService) GetCategories() ([]*models.Category, error) {
	return s.categoryRepo.GetCategories()
}

func (s *CategoryService) CreateCategory(category *models.Category) (*models.Category, error) {
	return s.categoryRepo.CreateCategory(category)
}
