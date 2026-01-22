package handler

import (
	"strconv"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/service"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CategoryHandler handles category-related requests
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// ============================================
// PUBLIC ENDPOINTS (No Auth Required)
// ============================================

// GetAllCategories handles getting all categories
// GET /api/categories
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	// Get all categories
	categories, err := h.categoryService.GetAll()
	if err != nil {
		utils.Error("CategoryHandler.GetAllCategories: Failed - Error=%v", err)
		return utils.InternalServerErrorResponse(c, "Failed to get categories")
	}

	return utils.SuccessResponse(c, "Categories retrieved successfully", categories)
}

// GetCategoryByID handles getting a category by ID
// GET /api/categories/:id
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	// Parse ID parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid category ID")
	}

	// Get category
	category, err := h.categoryService.GetByID(id)
	if err != nil {
		if err.Error() == "category not found" {
			return utils.NotFoundResponse(c, "Category not found")
		}
		utils.Error("CategoryHandler.GetCategoryByID: Failed - CategoryID=%d, Error=%v", id, err)
		return utils.InternalServerErrorResponse(c, "Failed to get category")
	}

	return utils.SuccessResponse(c, "Category retrieved successfully", category)
}

// GetProductsByCategory handles getting products for a category
// GET /api/categories/:id/products
func (h *CategoryHandler) GetProductsByCategory(c *fiber.Ctx) error {
	// Parse ID parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid category ID")
	}

	// Get pagination parameters
	page, limit := utils.GetPaginationParams(c)

	// Get products
	products, total, err := h.categoryService.GetProductsByCategory(id, page, limit)
	if err != nil {
		if err.Error() == "category not found" {
			return utils.NotFoundResponse(c, "Category not found")
		}
		utils.Error("CategoryHandler.GetProductsByCategory: Failed - CategoryID=%d, Error=%v", id, err)
		return utils.InternalServerErrorResponse(c, "Failed to get products")
	}

	return utils.PaginatedSuccessResponse(c, "Products retrieved successfully", products, page, limit, total)
}

// GetSubCategories handles getting subcategories of a parent category
// GET /api/categories/:id/subcategories
func (h *CategoryHandler) GetSubCategories(c *fiber.Ctx) error {
	// Parse ID parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid category ID")
	}

	// Get subcategories
	categories, err := h.categoryService.GetSubCategories(id)
	if err != nil {
		utils.Error("CategoryHandler.GetSubCategories: Failed - ParentID=%d, Error=%v", id, err)
		return utils.InternalServerErrorResponse(c, "Failed to get subcategories")
	}

	return utils.SuccessResponse(c, "Subcategories retrieved successfully", categories)
}

// ============================================
// ADMIN ENDPOINTS (Requires Admin Role)
// ============================================

// CreateCategory handles creating a new category (Admin only)
// POST /api/admin/categories
// Protected: Requires admin role
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	// Parse request
	var req models.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Warn("CategoryHandler.CreateCategory: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate input
	var validationErrors []utils.ValidationError

	if req.Name == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "name",
			Message: "Category name is required",
		})
	}

	if req.Slug == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "slug",
			Message: "Category slug is required",
		})
	}

	if len(validationErrors) > 0 {
		return utils.ValidationErrorsResponse(c, "Validation failed", validationErrors)
	}

	// Create category
	category, err := h.categoryService.Create(&req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "category slug already exists" {
			return utils.ConflictResponse(c, errMsg)
		}
		utils.Error("CategoryHandler.CreateCategory: Failed - Error=%v", err)
		return utils.InternalServerErrorResponse(c, "Failed to create category")
	}

	utils.Info("Category created: ID=%d, Name=%s", category.ID, category.Name)

	return utils.CreatedResponse(c, "Category created successfully", category)
}

// UpdateCategory handles updating a category (Admin only)
// PUT /api/admin/categories/:id
// Protected: Requires admin role
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	// Parse ID parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid category ID")
	}

	// Parse request body
	var req models.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Warn("CategoryHandler.UpdateCategory: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Update category
	category, err := h.categoryService.Update(id, &req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "category not found" {
			return utils.NotFoundResponse(c, "Category not found")
		}
		if errMsg == "category slug already exists" {
			return utils.ConflictResponse(c, errMsg)
		}
		utils.Error("CategoryHandler.UpdateCategory: Failed - CategoryID=%d, Error=%v", id, err)
		return utils.InternalServerErrorResponse(c, "Failed to update category")
	}

	utils.Info("Category updated: ID=%d, Name=%s", category.ID, category.Name)

	return utils.SuccessResponse(c, "Category updated successfully", category)
}

// DeleteCategory handles deleting a category (Admin only)
// DELETE /api/admin/categories/:id
// Protected: Requires admin role
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	// Parse ID parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid category ID")
	}

	// Delete category
	if err := h.categoryService.Delete(id); err != nil {
		errMsg := err.Error()
		if errMsg == "category not found" {
			return utils.NotFoundResponse(c, "Category not found")
		}
		if errMsg == "cannot delete category with subcategories" || 
		   errMsg == "cannot delete category with products" {
			return utils.BadRequestResponse(c, errMsg)
		}
		utils.Error("CategoryHandler.DeleteCategory: Failed - CategoryID=%d, Error=%v", id, err)
		return utils.InternalServerErrorResponse(c, "Failed to delete category")
	}

	utils.Info("Category deleted: ID=%d", id)

	return utils.SuccessMessage(c, "Category deleted successfully")
}

// ToggleCategoryStatus handles activating/deactivating a category (Admin only)
// PATCH /api/admin/categories/:id/status
// Protected: Requires admin role
func (h *CategoryHandler) ToggleCategoryStatus(c *fiber.Ctx) error {
	// Parse ID parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid category ID")
	}

	// Parse request body
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Toggle status
	if err := h.categoryService.ToggleStatus(id, req.IsActive); err != nil {
		errMsg := err.Error()
		if errMsg == "category not found" {
			return utils.NotFoundResponse(c, "Category not found")
		}
		utils.Error("CategoryHandler.ToggleCategoryStatus: Failed - CategoryID=%d, Error=%v", id, err)
		return utils.InternalServerErrorResponse(c, "Failed to update category status")
	}

	// Determine status message
	status := "deactivated"
	if req.IsActive {
		status = "activated"
	}

	utils.Info("Category status toggled: ID=%d, IsActive=%v", id, req.IsActive)

	return utils.SuccessMessage(c, "Category "+status+" successfully")
}

// GetCategoryStats handles getting category statistics (Admin only)
// GET /api/admin/categories/stats
// Protected: Requires admin role
func (h *CategoryHandler) GetCategoryStats(c *fiber.Ctx) error {
	// Get statistics
	stats, err := h.categoryService.GetCategoryStats()
	if err != nil {
		utils.Error("CategoryHandler.GetCategoryStats: Failed - Error=%v", err)
		return utils.InternalServerErrorResponse(c, "Failed to get category statistics")
	}

	return utils.SuccessResponse(c, "Category statistics retrieved successfully", stats)
}