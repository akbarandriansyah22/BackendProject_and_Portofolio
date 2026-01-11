package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"e-commerce-api/server/internal/models"
	"e-commerce-api/server/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles creating a new product (Admin only)
// POST /api/admin/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

// GetAllProducts handles getting all products with pagination and filters
// GET /api/products
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	categoryID := c.Query("category_id")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Build filters
	filters := make(map[string]interface{})
	if search != "" {
		filters["search"] = search
	}
	if categoryID != "" {
		if catID, err := strconv.Atoi(categoryID); err == nil {
			filters["category_id"] = catID
		}
	}
	if minPrice != "" {
		if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters["min_price"] = min
		}
	}
	if maxPrice != "" {
		if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters["max_price"] = max
		}
	}
	filters["sort_by"] = sortBy
	filters["sort_order"] = sortOrder

	products, total, err := h.productService.GetAllProducts(page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    products,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}

// GetProductByID handles getting a product by ID
// GET /api/products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid product ID",
		})
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    product,
	})
}

// GetProductBySlug handles getting a product by slug
// GET /api/products/slug/:slug
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	product, err := h.productService.GetProductBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    product,
	})
}

// UpdateProduct handles updating a product (Admin only)
// PUT /api/admin/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid product ID",
		})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	product, err := h.productService.UpdateProduct(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DeleteProduct handles deleting a product (Admin only)
// DELETE /api/admin/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid product ID",
		})
		return
	}

	if err := h.productService.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}

// UpdateProductStock handles updating product stock (Admin only)
// PATCH /api/admin/products/:id/stock
func (h *ProductHandler) UpdateProductStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid product ID",
		})
		return
	}

	var req struct {
		Stock int `json:"stock" binding:"required,gte=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.productService.UpdateStock(id, req.Stock); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product stock updated successfully",
	})
}

// ToggleProductStatus handles activating/deactivating a product (Admin only)
// PATCH /api/admin/products/:id/status
func (h *ProductHandler) ToggleProductStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid product ID",
		})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.productService.ToggleStatus(id, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	status := "deactivated"
	if req.IsActive {
		status = "activated"
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product " + status + " successfully",
	})
}

// GetProductsByCategory handles getting products by category
// GET /api/products/category/:category_id
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid category ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	products, total, err := h.productService.GetProductsByCategory(categoryID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    products,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}

// SearchProducts handles searching products
// GET /api/products/search
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Search query is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	products, total, err := h.productService.SearchProducts(query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    products,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}

// GetFeaturedProducts handles getting featured/popular products
// GET /api/products/featured
func (h *ProductHandler) GetFeaturedProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit < 1 || limit > 50 {
		limit = 10
	}

	products, err := h.productService.GetFeaturedProducts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    products,
	})
}

// GetLowStockProducts handles getting products with low stock (Admin only)
// GET /api/admin/products/low-stock
func (h *ProductHandler) GetLowStockProducts(c *gin.Context) {
	threshold, _ := strconv.Atoi(c.DefaultQuery("threshold", "10"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if threshold < 1 {
		threshold = 10
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	products, err := h.productService.GetLowStockProducts(threshold, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    products,
	})
}

// GetProductStats handles getting product statistics (Admin only)
// GET /api/admin/products/stats
func (h *ProductHandler) GetProductStats(c *gin.Context) {
	stats, err := h.productService.GetProductStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// BulkUpdateStock handles bulk stock update (Admin only)
// POST /api/admin/products/bulk-stock
func (h *ProductHandler) BulkUpdateStock(c *gin.Context) {
	var req struct {
		Updates []struct {
			ProductID int `json:"product_id" binding:"required"`
			Stock     int `json:"stock" binding:"required,gte=0"`
		} `json:"updates" binding:"required,dive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	successCount := 0
	errors := []string{}

	for _, update := range req.Updates {
		if err := h.productService.UpdateStock(update.ProductID, update.Stock); err != nil {
			errors = append(errors, err.Error())
		} else {
			successCount++
		}
	}

	response := map[string]interface{}{
		"success_count": successCount,
		"total_count":   len(req.Updates),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Bulk stock update completed",
		Data:    response,
	})
}

// GetProductWithCategories handles getting a product with its categories
// GET /api/products/:id/categories
func (h *ProductHandler) GetProductWithCategories(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid product ID",
		})
		return
	}

	productWithCategories, err := h.productService.GetProductWithCategories(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    productWithCategories,
	})
}