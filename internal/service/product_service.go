package service

import (
	"fmt"
	"database/sql"

	"e-commerce-api/server/internal/models"
	"e-commerce-api/server/internal/repository"
	"e-commerce-api/server/internal/utils"
)

// ProductService handles product business logic
type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo *repository.ProductRepository, categoryRepo *repository.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// GetAll gets all products with pagination
func (s *ProductService) GetAll(page, limit int) ([]models.Product, int64, error) {
	offset := (page - 1) * limit

	products, total, err := s.productRepo.GetAll(limit, offset)
	if err != nil {
		utils.Error("ProductService.GetAll: Failed - Error=%v", err)
		return nil, 0, fmt.Errorf("failed to get products")
	}

	return products, total, nil
}

// GetActive gets all active products with pagination
func (s *ProductService) GetActive(page, limit int) ([]models.Product, int64, error) {
	offset := (page - 1) * limit

	products, total, err := s.productRepo.GetActive(limit, offset)
	if err != nil {
		utils.Error("ProductService.GetActive: Failed - Error=%v", err)
		return nil, 0, fmt.Errorf("failed to get products")
	}

	return products, total, nil
}

// GetByID gets product by ID
func (s *ProductService) GetByID(id int) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}

// GetBySlug gets product by slug
func (s *ProductService) GetBySlug(slug string) (*models.Product, error) {
	product, err := s.productRepo.GetBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}

// GetByCategory gets products by category
func (s *ProductService) GetByCategory(categoryID, page, limit int) ([]models.Product, int64, error) {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("category not found")
	}

	offset := (page - 1) * limit

	products, total, err := s.productRepo.GetByCategory(categoryID, limit, offset)
	if err != nil {
		utils.Error("ProductService.GetByCategory: Failed - CategoryID=%d, Error=%v", categoryID, err)
		return nil, 0, fmt.Errorf("failed to get products")
	}

	return products, total, nil
}

// Search searches products by name or description
func (s *ProductService) Search(keyword string, page, limit int) ([]models.Product, int64, error) {
	if keyword == "" {
		return nil, 0, fmt.Errorf("search keyword is required")
	}

	offset := (page - 1) * limit

	products, total, err := s.productRepo.Search(keyword, limit, offset)
	if err != nil {
		utils.Error("ProductService.Search: Failed - Keyword=%s, Error=%v", keyword, err)
		return nil, 0, fmt.Errorf("failed to search products")
	}

	return products, total, nil
}

// Create creates a new product
func (s *ProductService) Create(product *models.Product) error {
	// Validate product
	if err := s.validateProduct(product); err != nil {
		return err
	}

	// Generate slug if not provided
	if product.Slug == "" {
		product.Slug = generateSlug(product.Name)
	}

	// Check if slug already exists
	exists, err := s.productRepo.SlugExists(product.Slug)
	if err != nil {
		utils.Error("ProductService.Create: Failed to check slug - Error=%v", err)
		return fmt.Errorf("failed to create product")
	}
	if exists {
		return fmt.Errorf("product slug already exists")
	}

	// Check if SKU already exists (if provided)
	// ✅ FIXED: product.SKU (uppercase)
	if product.SKU.Valid && product.SKU.String != "" {
		exists, err := s.productRepo.SKUExists(product.SKU.String)
		if err != nil {
			utils.Error("ProductService.Create: Failed to check SKU - Error=%v", err)
			return fmt.Errorf("failed to create product")
		}
		if exists {
			return fmt.Errorf("product SKU already exists")
		}
	}

	// Create product
	if err := s.productRepo.Create(product); err != nil {
		utils.Error("ProductService.Create: Failed - Error=%v", err)
		return fmt.Errorf("failed to create product")
	}

	utils.Info("Product created: ID=%d, Name=%s", product.ID, product.Name)
	return nil
}

// Update updates a product
func (s *ProductService) Update(product *models.Product) error {
	// Check if product exists
	existing, err := s.productRepo.GetByID(product.ID)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	// Validate product
	if err := s.validateProduct(product); err != nil {
		return err
	}

	// Check if slug changed and is unique
	if product.Slug != existing.Slug {
		exists, err := s.productRepo.SlugExistsExcludingProduct(product.Slug, product.ID)
		if err != nil {
			utils.Error("ProductService.Update: Failed to check slug - Error=%v", err)
			return fmt.Errorf("failed to update product")
		}
		if exists {
			return fmt.Errorf("product slug already exists")
		}
	}

	// Check if SKU changed and is unique
	if product.SKU.Valid && product.SKU.String != "" && product.SKU.String != existing.SKU.String {
		exists, err := s.productRepo.SKUExistsExcludingProduct(product.SKU.String, product.ID)
		if err != nil {
			utils.Error("ProductService.Update: Failed to check SKU - Error=%v", err)
			return fmt.Errorf("failed to update product")
		}
		if exists {
			return fmt.Errorf("product SKU already exists")
		}
	}

	// Update product
	if err := s.productRepo.Update(product); err != nil {
		utils.Error("ProductService.Update: Failed - ProductID=%d, Error=%v", product.ID, err)
		return fmt.Errorf("failed to update product")
	}

	utils.Info("Product updated: ID=%d, Name=%s", product.ID, product.Name)
	return nil
}

// Delete deletes a product (soft delete)
func (s *ProductService) Delete(id int) error {
	// Check if product exists
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	// Soft delete (set is_active = false)
	if err := s.productRepo.Delete(id); err != nil {
		utils.Error("ProductService.Delete: Failed - ProductID=%d, Error=%v", id, err)
		return fmt.Errorf("failed to delete product")
	}

	utils.Info("Product deleted: ID=%d", id)
	return nil
}

// UpdateStock updates product stock
func (s *ProductService) UpdateStock(id, quantity int) error {
	// Check if product exists
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	// Validate quantity
	if quantity < 0 {
		return fmt.Errorf("stock quantity cannot be negative")
	}

	// Update stock
	if err := s.productRepo.UpdateStock(id, quantity); err != nil {
		utils.Error("ProductService.UpdateStock: Failed - ProductID=%d, Error=%v", id, err)
		return fmt.Errorf("failed to update stock")
	}

	utils.Info("Product stock updated: ID=%d, OldStock=%d, NewStock=%d", id, product.Stock, quantity)
	return nil
}

// IncrementStock increments product stock
func (s *ProductService) IncrementStock(id, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	if err := s.productRepo.IncrementStock(id, quantity); err != nil {
		utils.Error("ProductService.IncrementStock: Failed - ProductID=%d, Error=%v", id, err)
		return fmt.Errorf("failed to increment stock")
	}

	utils.Info("Product stock incremented: ID=%d, Quantity=%d", id, quantity)
	return nil
}

// DecrementStock decrements product stock
func (s *ProductService) DecrementStock(id, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	// Check current stock
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	if product.Stock < quantity {
		return fmt.Errorf("insufficient stock. Available: %d", product.Stock)
	}

	if err := s.productRepo.DecrementStock(id, quantity); err != nil {
		utils.Error("ProductService.DecrementStock: Failed - ProductID=%d, Error=%v", id, err)
		return fmt.Errorf("failed to decrement stock")
	}

	utils.Info("Product stock decremented: ID=%d, Quantity=%d", id, quantity)
	return nil
}
// CreateProduct creates a new product (wrapper untuk handler)
func (s *ProductService) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if req.Slug == "" {
		return nil, fmt.Errorf("product slug is required")
	}
	if req.Price <= 0 {
		return nil, fmt.Errorf("product price must be greater than 0")
	}

	// Check if slug already exists
	exists, err := s.productRepo.SlugExists(req.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("product with slug '%s' already exists", req.Slug)
	}

	// Create product
	product := &models.Product{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         sql.NullString{String: req.SKU, Valid: req.SKU != ""},
		ImageURL:    sql.NullString{String: req.ImageURL, Valid: req.ImageURL != ""},
		IsActive:    true,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	// Add product to categories if provided
	if len(req.CategoryIDs) > 0 {
		for _, categoryID := range req.CategoryIDs {
			if err := s.productRepo.AddToCategory(product.ID, categoryID); err != nil {
				// Log error but don't fail the entire operation
				continue
			}
		}
	}

	utils.Info("Product created: ID=%d, Name=%s", product.ID, product.Name)
	return product, nil
}

// GetAllProducts retrieves all products with filters
func (s *ProductService) GetAllProducts(page, limit int, filters map[string]interface{}) ([]models.Product, int64, error) {
	offset := (page - 1) * limit
	return s.productRepo.GetAllWithFilters(filters, limit, offset)
}

// GetProductByID retrieves a product by ID (alias untuk GetByID)
func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.GetByID(id)
}

// GetProductBySlug retrieves a product by slug (alias untuk GetBySlug)
func (s *ProductService) GetProductBySlug(slug string) (*models.Product, error) {
	return s.GetBySlug(slug)
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(id int, req *models.UpdateProductRequest) (*models.Product, error) {
	// Get existing product
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	// Update fields if provided
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = sql.NullString{String: req.Description, Valid: true}
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.SKU != "" {
		product.SKU = sql.NullString{String: req.SKU, Valid: true}
	}
	if req.ImageURL != "" {
		product.ImageURL = sql.NullString{String: req.ImageURL, Valid: true}
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	// Update product
	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	// Update categories if provided
	if len(req.CategoryIDs) > 0 {
		// Remove existing categories
		if err := s.productRepo.RemoveFromAllCategories(product.ID); err != nil {
			return nil, err
		}

		// Add new categories
		for _, categoryID := range req.CategoryIDs {
			if err := s.productRepo.AddToCategory(product.ID, categoryID); err != nil {
				continue
			}
		}
	}

	utils.Info("Product updated: ID=%d, Name=%s", product.ID, product.Name)
	return product, nil
}

// DeleteProduct deletes a product (alias untuk Delete)
func (s *ProductService) DeleteProduct(id int) error {
	return s.Delete(id)
}

// ToggleStatus activates or deactivates a product
func (s *ProductService) ToggleStatus(id int, isActive bool) error {
	// Check if product exists
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	if err := s.productRepo.UpdateStatus(id, isActive); err != nil {
		utils.Error("ProductService.ToggleStatus: Failed - ProductID=%d, Error=%v", id, err)
		return fmt.Errorf("failed to update product status")
	}

	utils.Info("Product status toggled: ID=%d, IsActive=%v", id, isActive)
	return nil
}

// GetProductsByCategory retrieves products by category (alias untuk GetByCategory)
func (s *ProductService) GetProductsByCategory(categoryID, page, limit int) ([]models.Product, int64, error) {
	return s.GetByCategory(categoryID, page, limit)
}

// SearchProducts searches products by name or description (alias untuk Search)
func (s *ProductService) SearchProducts(query string, page, limit int) ([]models.Product, int64, error) {
	return s.Search(query, page, limit)
}

// GetFeaturedProducts retrieves featured/popular products
func (s *ProductService) GetFeaturedProducts(limit int) ([]models.Product, error) {
	products, err := s.productRepo.GetFeatured(limit)
	if err != nil {
		utils.Error("ProductService.GetFeaturedProducts: Failed - Error=%v", err)
		return nil, fmt.Errorf("failed to get featured products")
	}

	return products, nil
}

// GetLowStockProducts retrieves products with low stock
func (s *ProductService) GetLowStockProducts(threshold, limit int) ([]models.Product, error) {
	products, err := s.productRepo.GetLowStock(threshold, limit)
	if err != nil {
		utils.Error("ProductService.GetLowStockProducts: Failed - Error=%v", err)
		return nil, fmt.Errorf("failed to get low stock products")
	}

	return products, nil
}

// GetProductStats retrieves product statistics
func (s *ProductService) GetProductStats() (map[string]interface{}, error) {
	totalProducts, err := s.productRepo.CountTotal()
	if err != nil {
		return nil, err
	}

	activeProducts, err := s.productRepo.CountByStatus(true)
	if err != nil {
		return nil, err
	}

	inactiveProducts, err := s.productRepo.CountByStatus(false)
	if err != nil {
		return nil, err
	}

	lowStockProducts, err := s.productRepo.CountLowStock(10)
	if err != nil {
		return nil, err
	}

	outOfStockProducts, err := s.productRepo.CountOutOfStock()
	if err != nil {
		return nil, err
	}

	totalStockValue, err := s.productRepo.GetTotalStockValue()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_products":        totalProducts,
		"active_products":       activeProducts,
		"inactive_products":     inactiveProducts,
		"low_stock_products":    lowStockProducts,
		"out_of_stock_products": outOfStockProducts,
		"total_stock_value":     totalStockValue,
	}

	return stats, nil
}

// GetProductWithCategories retrieves a product with its categories
func (s *ProductService) GetProductWithCategories(id int) (*models.ProductWithCategories, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	categories, err := s.productRepo.GetCategories(id)
	if err != nil {
		return nil, err
	}

	return &models.ProductWithCategories{
		Product:    *product,
		Categories: categories,
	}, nil
}

// GetActiveProducts retrieves only active products
func (s *ProductService) GetActiveProducts(page, limit int) ([]models.Product, int64, error) {
	return s.GetActive(page, limit)
}

// CheckStock checks if product has enough stock
func (s *ProductService) CheckStock(productID, quantity int) (bool, error) {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return false, err
	}

	return product.Stock >= quantity, nil
}

// DeductStock deducts stock from a product (alias untuk DecrementStock)
func (s *ProductService) DeductStock(productID, quantity int) error {
	return s.DecrementStock(productID, quantity)
}

// RestoreStock restores stock to a product (alias untuk IncrementStock)
func (s *ProductService) RestoreStock(productID, quantity int) error {
	return s.IncrementStock(productID, quantity)
}

// BulkUpdateStock updates stock for multiple products
func (s *ProductService) BulkUpdateStock(updates []struct {
	ProductID int
	Stock     int
}) (int, []error) {
	successCount := 0
	errors := []error{}

	for _, update := range updates {
		if err := s.UpdateStock(update.ProductID, update.Stock); err != nil {
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	return successCount, errors
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// validateProduct validates product data
func (s *ProductService) validateProduct(product *models.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}

	if product.Price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}

	if product.Stock < 0 {
		return fmt.Errorf("product stock cannot be negative")
	}

	return nil
}

// generateSlug generates a URL-friendly slug from name
func generateSlug(name string) string {
	slug := ""
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			slug += string(char)
		} else if char == ' ' {
			slug += "-"
		}
	}

	// Convert to lowercase
	lowerSlug := ""
	for _, char := range slug {
		if char >= 'A' && char <= 'Z' {
			lowerSlug += string(char + 32)
		} else {
			lowerSlug += string(char)
		}
	}

	return lowerSlug
}