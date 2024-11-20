package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

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
	productService := services.NewProductService(productRepository)

	return &ProductHandler{
		productService: productService,
	}
}

func writeErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	response := map[string]interface{}{
		"message": err.Error(),
		"status":  statusCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (ph *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		logger.Logger.Warn("Invalid request payload", slog.String("error", err.Error()))
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err := ph.productService.CreateProduct(
		models.NewProduct(product.Title, product.Description, product.Image, product.Price),
	)
	if err != nil {
		logger.Logger.Error("Failed to create product", slog.String("error", err.Error()))
		writeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Product created successfully", slog.String("product_id", product.Id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Product created successfully",
		"status":  http.StatusCreated,
	})
}

func (ph *ProductHandler) Index(w http.ResponseWriter, r *http.Request) {
	products, err := ph.productService.GetProducts(10, 0)
	if err != nil {
		logger.Logger.Error("Failed to get products", slog.String("error", err.Error()))
		writeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	logger.Logger.Info("Products retrieved successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (ph *ProductHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	product, err := ph.productService.GetProductByID(id)
	if err != nil {
		logger.Logger.Error("Failed to get product", slog.String("error", err.Error()))
		writeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if product == nil {
		logger.Logger.Warn("Product not found", slog.String("product_id", id))
		writeErrorResponse(w, err, http.StatusNotFound)
		return
	}
	logger.Logger.Info("Product retrieved successfully", slog.String("product_id", product.Id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}