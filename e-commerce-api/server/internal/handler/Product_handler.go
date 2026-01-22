package handler

import (
	"strconv"

	"github.com/akbarandriansyah22/Devops_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/Devops_Portofolio/e-commerce-api/server/internal/service"
	"github.com/akbarandriansyah22/Devops_Portofolio/e-commerce-api/server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// ProductHandler handles product-related HTTP requests
// ✅ FIXED: Menggunakan FIBER (bukan Gin)
// ✅ CLEAN: Hanya handle HTTP, semua logic di service
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// ============================================
// PUBLIC ENDPOINTS (No Auth Required)
// ============================================

// GetAllProducts handles getting all products with pagination and filters
// GET /api/products
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	// 1. Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Build filters from query params
	filters := make(map[string]interface{})
	
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}
	
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if catID, err := strconv.Atoi(categoryIDStr); err == nil {
			filters["category_id"] = catID
		}
	}
	
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters["min_price"] = minPrice
		}
	}
	
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters["max_price"] = maxPrice
		}
	}
	
	filters["sort_by"] = c.Query("sort_by", "created_at")
	filters["sort_order"] = c.Query("sort_order", "desc")

	// 2. Call service
	products, total, err := h.productService.GetAllProducts(page, limit, filters)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 3. Return paginated response
	return utils.PaginatedSuccessResponse(c, "Products retrieved successfully", 
		products, page, limit, total)
}

// GetProductByID handles getting a product by ID
// GET /api/products/:id
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	// 1. Parse ID parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid product ID")
	}

	// 2. Call service
	product, err := h.productService.GetProductByID(id)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 3. Return response
	return utils.SuccessResponse(c, "Product retrieved successfully", product)
}

// GetProductBySlug handles getting a product by slug
// GET /api/products/slug/:slug
func (h *ProductHandler) GetProductBySlug(c *fiber.Ctx) error {
	// 1. Get slug parameter
	slug := c.Params("slug")
	if slug == "" {
		return utils.BadRequestResponse(c, "Product slug is required")
	}

	// 2. Call service
	product, err := h.productService.GetProductBySlug(slug)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 3. Return response
	return utils.SuccessResponse(c, "Product retrieved successfully", product)
}

// SearchProducts handles searching products
// GET /api/products/search
func (h *ProductHandler) SearchProducts(c *fiber.Ctx) error {
	// 1. Get search query
	query := c.Query("q")
	if query == "" {
		return utils.BadRequestResponse(c, "Search query is required")
	}

	// 2. Parse pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 3. Call service
	products, total, err := h.productService.SearchProducts(query, page, limit)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 4. Return response
	return utils.PaginatedSuccessResponse(c, "Products found", 
		products, page, limit, total)
}

// GetFeaturedProducts handles getting featured/popular products
// GET /api/products/featured
func (h *ProductHandler) GetFeaturedProducts(c *fiber.Ctx) error {
	// 1. Parse limit
	limit := c.QueryInt("limit", 10)
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// 2. Call service
	products, err := h.productService.GetFeaturedProducts(limit)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 3. Return response
	return utils.SuccessResponse(c, "Featured products retrieved successfully", products)
}

// GetProductsByCategory handles getting products by category
// GET /api/products/category/:category_id
func (h *ProductHandler) GetProductsByCategory(c *fiber.Ctx) error {
	// 1. Parse category ID
	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid category ID")
	}

	// 2. Parse pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 3. Call service
	products, total, err := h.productService.GetProductsByCategory(categoryID, page, limit)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 4. Return response
	return utils.PaginatedSuccessResponse(c, "Products retrieved successfully", 
		products, page, limit, total)
}

// GetProductWithCategories handles getting a product with its categories
// GET /api/products/:id/categories
func (h *ProductHandler) GetProductWithCategories(c *fiber.Ctx) error {
	// 1. Parse ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid product ID")
	}

	// 2. Call service
	productWithCategories, err := h.productService.GetProductWithCategories(id)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 3. Return response
	return utils.SuccessResponse(c, "Product retrieved successfully", productWithCategories)
}

// ============================================
// ADMIN ENDPOINTS (Requires Admin Role)
// ============================================

// CreateProduct handles creating a new product (Admin only)
// POST /api/admin/products
// Protected: Requires admin role
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	// 1. Parse request
	var req models.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Warn("ProductHandler.CreateProduct: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// 2. Call service
	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 3. Return response
	return utils.CreatedResponse(c, "Product created successfully", product)
}

// UpdateProduct handles updating a product (Admin only)
// PUT /api/admin/products/:id
// Protected: Requires admin role
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	// 1. Parse ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid product ID")
	}

	// 2. Parse request
	var req models.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Warn("ProductHandler.UpdateProduct: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// 3. Call service
	product, err := h.productService.UpdateProduct(id, &req)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 4. Return response
	return utils.SuccessResponse(c, "Product updated successfully", product)
}

// DeleteProduct handles deleting a product (Admin only)
// DELETE /api/admin/products/:id
// Protected: Requires admin role
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	// 1. Parse ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid product ID")
	}

	// 2. Call service
	if err := h.productService.DeleteProduct(id); err != nil {
		return h.handleProductError(c, err)
	}

	// 3. Return response
	return utils.SuccessMessage(c, "Product deleted successfully")
}

// UpdateProductStock handles updating product stock (Admin only)
// PATCH /api/admin/products/:id/stock
// Protected: Requires admin role
func (h *ProductHandler) UpdateProductStock(c *fiber.Ctx) error {
	// 1. Parse ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid product ID")
	}

	// 2. Parse request
	var req struct {
		Stock int `json:"stock"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// 3. Call service
	if err := h.productService.UpdateStock(id, req.Stock); err != nil {
		return h.handleProductError(c, err)
	}

	// 4. Return response
	return utils.SuccessMessage(c, "Product stock updated successfully")
}

// ToggleProductStatus handles activating/deactivating a product (Admin only)
// PATCH /api/admin/products/:id/status
// Protected: Requires admin role
func (h *ProductHandler) ToggleProductStatus(c *fiber.Ctx) error {
	// 1. Parse ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid product ID")
	}

	// 2. Parse request
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// 3. Call service
	if err := h.productService.ToggleStatus(id, req.IsActive); err != nil {
		return h.handleProductError(c, err)
	}

	// 4. Determine status message
	status := "deactivated"
	if req.IsActive {
		status = "activated"
	}

	// 5. Return response
	return utils.SuccessMessage(c, "Product "+status+" successfully")
}

// GetLowStockProducts handles getting products with low stock (Admin only)
// GET /api/admin/products/low-stock
// Protected: Requires admin role
func (h *ProductHandler) GetLowStockProducts(c *fiber.Ctx) error {
	// 1. Parse query params
	threshold := c.QueryInt("threshold", 10)
	limit := c.QueryInt("limit", 20)

	if threshold < 1 {
		threshold = 10
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// 2. Call service
	products, err := h.productService.GetLowStockProducts(threshold, limit)
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 3. Return response
	return utils.SuccessResponse(c, "Low stock products retrieved successfully", products)
}

// GetProductStats handles getting product statistics (Admin only)
// GET /api/admin/products/stats
// Protected: Requires admin role
func (h *ProductHandler) GetProductStats(c *fiber.Ctx) error {
	// 1. Call service
	stats, err := h.productService.GetProductStats()
	if err != nil {
		return h.handleProductError(c, err)
	}

	// 2. Return response
	return utils.SuccessResponse(c, "Product statistics retrieved successfully", stats)
}

// BulkUpdateStock handles bulk stock update (Admin only)
// POST /api/admin/products/bulk-stock
// Protected: Requires admin role
func (h *ProductHandler) BulkUpdateStock(c *fiber.Ctx) error {
	// 1. Parse request
	var req struct {
		Updates []struct {
			ProductID int `json:"product_id"`
			Stock     int `json:"stock"`
		} `json:"updates"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// 2. Validate
	if len(req.Updates) == 0 {
		return utils.BadRequestResponse(c, "No updates provided")
	}

	// 3. Process bulk update
	successCount := 0
	errors := []string{}

	for _, update := range req.Updates {
		if err := h.productService.UpdateStock(update.ProductID, update.Stock); err != nil {
			errors = append(errors, err.Error())
		} else {
			successCount++
		}
	}

	// 4. Build response
	response := fiber.Map{
		"success_count": successCount,
		"total_count":   len(req.Updates),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	// 5. Return response
	return utils.SuccessResponse(c, "Bulk stock update completed", response)
}

// ============================================
// ERROR HANDLING HELPER
// ============================================

// handleProductError maps service errors to HTTP responses
func (h *ProductHandler) handleProductError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// Map specific errors to HTTP status codes
	switch errMsg {
	// Not found errors → 404
	case "product not found",
		"category not found":
		return utils.NotFoundResponse(c, errMsg)

	// Validation errors → 400
	case "product name is required",
		"product slug is required",
		"product price must be greater than 0",
		"product stock cannot be negative",
		"quantity must be greater than 0",
		"search keyword is required":
		return utils.BadRequestResponse(c, errMsg)

	// Conflict errors → 409
	case "product slug already exists",
		"product SKU already exists":
		return utils.ConflictResponse(c, errMsg)

	// Stock errors → 400
	default:
		if len(errMsg) > 20 && errMsg[:20] == "insufficient stock" {
			return utils.BadRequestResponse(c, errMsg)
		}

		// Generic error → 500
		utils.Error("ProductHandler: Unhandled error - %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to process product request")
	}
}