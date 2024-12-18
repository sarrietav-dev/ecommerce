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

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	productRepository := repository.NewProductRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	productService := services.NewProductService(productRepository, categoryRepository)

	return &ProductHandler{
		productService: productService,
	}
}

func (ph *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var requestProduct models.Product

	if err := jsonapi.UnmarshalPayload(r.Body, &requestProduct); err != nil {
		logger.Logger.Warn("Invalid request payload", slog.String("error", err.Error()))
		internal.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	product := models.NewProduct(requestProduct.Title, requestProduct.Description, requestProduct.Image, requestProduct.Price)
	product.Categories = requestProduct.Categories
	newProduct, err := ph.productService.CreateProduct(product)
	if err != nil {
		logger.Logger.Error("Failed to create product", slog.String("error", err.Error()))
		internal.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Product created successfully", slog.String("product_id", newProduct.Id))
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)
	jsonapi.MarshalPayload(w, newProduct)
}

func (ph *ProductHandler) Index(w http.ResponseWriter, r *http.Request) {
	products, err := ph.productService.GetProducts(10, 0)
	if err != nil {
		logger.Logger.Error("Failed to get products", slog.String("error", err.Error()))
		internal.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	logger.Logger.Info("Products retrieved successfully")
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)
	jsonapi.MarshalPayload(w, products)
}

func (ph *ProductHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	product, err := ph.productService.GetProductByID(id)
	if err != nil {
		logger.Logger.Error("Failed to get product", slog.String("error", err.Error()))
		internal.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if product == nil {
		logger.Logger.Warn("Product not found", slog.String("product_id", id))
		internal.WriteErrorResponse(w, err, http.StatusNotFound)
		return
	}
	logger.Logger.Info("Product retrieved successfully", slog.String("product_id", product.Id))
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)
	jsonapi.MarshalPayload(w, product)
}
