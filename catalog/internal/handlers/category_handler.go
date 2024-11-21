package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/sarrietav-dev/ecommerce/catalog/internal"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/logger"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/models"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/repository"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/services"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) Index(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		logger.Logger.Error("Failed to get categories", slog.String("error", err.Error()))
		internal.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Categories retrieved successfully")
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)
	jsonapi.MarshalPayload(w, categories)
}

func (h *CategoryHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	category, err := h.categoryService.GetCategoryByID(id)
	if err != nil {
		logger.Logger.Error("Failed to get category", slog.String("error", err.Error()))
		internal.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Category retrieved successfully", slog.String("category_id", id))
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)
	jsonapi.MarshalPayload(w, category)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var requestCategory models.Category

	if err := jsonapi.UnmarshalPayload(r.Body, &requestCategory); err != nil {
		logger.Logger.Warn("Invalid request payload", slog.String("error", err.Error()))
		internal.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	category := models.NewCategory(requestCategory.Name, requestCategory.Description)
	newCategory, err := h.categoryService.CreateCategory(category)
	if err != nil {
		logger.Logger.Error("Failed to create category", slog.String("error", err.Error()))
		internal.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Category created successfully", slog.String("category_id", newCategory.Id))
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)
	jsonapi.MarshalPayload(w, newCategory)
}
