package service

import (
	"database/sql"
	"fmt"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/repository"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/utils"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	productRepo  *repository.ProductRepository // ✅ ADD: Needed for GetProductsByCategory
}

// NewCategoryService creates a new category service
// FIX: Add productRepo parameter
func NewCategoryService(categoryRepo *repository.CategoryRepository, productRepo *repository.ProductRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}

// GetAll gets all categories
func (s *CategoryService) GetAll() ([]models.Category, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		utils.Error("CategoryService.GetAll: Failed - Error=%v", err)
		return nil, fmt.Errorf("failed to get categories")
	}
	return categories, nil
}

// GetByID gets category by ID
func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

// GetBySlug gets category by slug
// ADD: Missing method
func (s *CategoryService) GetBySlug(slug string) (*models.Category, error) {
	category, err := s.categoryRepo.GetBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

// Create creates a new category
func (s *CategoryService) Create(req *models.CreateCategoryRequest) (*models.Category, error) {
	// Validate
	if req.Name == "" {
		return nil, fmt.Errorf("category name is required")
	}
	if req.Slug == "" {
		return nil, fmt.Errorf("category slug is required")
	}

	// Check if slug exists
	exists, err := s.categoryRepo.SlugExists(req.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to create category")
	}
	if exists {
		return nil, fmt.Errorf("category slug already exists")
	}

	// Create category
	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		IsActive:    true,
	}

	if req.ParentID != nil {
		category.ParentID = sql.NullInt32{Int32: int32(*req.ParentID), Valid: true}
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category")
	}

	utils.Info("Category created: ID=%d, Name=%s", category.ID, category.Name)
	return category, nil
}

// Update updates a category
func (s *CategoryService) Update(id int, req *models.UpdateCategoryRequest) (*models.Category, error) {
	// Get existing category
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found")
	}

	// Update fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		// Check if new slug exists
		if req.Slug != category.Slug {
			exists, err := s.categoryRepo.SlugExistsExcludingCategory(req.Slug, id)
			if err != nil {
				return nil, fmt.Errorf("failed to update category")
			}
			if exists {
				return nil, fmt.Errorf("category slug already exists")
			}
		}
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = sql.NullString{String: req.Description, Valid: true}
	}
	if req.ParentID != nil {
		category.ParentID = sql.NullInt32{Int32: int32(*req.ParentID), Valid: true}
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, fmt.Errorf("failed to update category")
	}

	utils.Info("Category updated: ID=%d, Name=%s", category.ID, category.Name)
	return category, nil
}

// Delete deletes a category (soft delete)
func (s *CategoryService) Delete(id int) error {
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("category not found")
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete category")
	}

	utils.Info("Category deleted: ID=%d", id)
	return nil
}

// ToggleStatus activates or deactivates a category
// FIX: Use Update() instead of non-existent UpdateStatus()
func (s *CategoryService) ToggleStatus(id int, isActive bool) error {
	// Get category
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("category not found")
	}

	// Update is_active field
	category.IsActive = isActive
	
	// Use existing Update method
	if err := s.categoryRepo.Update(category); err != nil {
		return fmt.Errorf("failed to update category status")
	}

	utils.Info("Category status toggled: ID=%d, IsActive=%v", id, isActive)
	return nil
}

// GetProductsByCategory gets products for a category
// FIX: Use ProductRepository instead of non-existent CategoryRepository.GetProducts()
func (s *CategoryService) GetProductsByCategory(categoryID, page, limit int) ([]models.Product, int64, error) {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("category not found")
	}

	// Use ProductRepository to get products by category
	offset := (page - 1) * limit
	products, total, err := s.productRepo.GetByCategory(categoryID, limit, offset)
	if err != nil {
		utils.Error("CategoryService.GetProductsByCategory: Failed - CategoryID=%d, Error=%v", categoryID, err)
		return nil, 0, fmt.Errorf("failed to get products")
	}

	return products, total, nil
}

// GetSubCategories gets sub-categories of a category
// FIX: Use correct method name GetSubCategories() instead of GetByParentID()
func (s *CategoryService) GetSubCategories(parentID int) ([]models.Category, error) {
	categories, err := s.categoryRepo.GetSubCategories(parentID)
	if err != nil {
		utils.Error("CategoryService.GetSubCategories: Failed - ParentID=%d, Error=%v", parentID, err)
		return nil, fmt.Errorf("failed to get sub-categories")
	}
	return categories, nil
}

// GetParentCategories gets all parent categories (categories without parent)
// ADD: Missing method
func (s *CategoryService) GetParentCategories() ([]models.Category, error) {
	categories, err := s.categoryRepo.GetParentCategories()
	if err != nil {
		utils.Error("CategoryService.GetParentCategories: Failed - Error=%v", err)
		return nil, fmt.Errorf("failed to get parent categories")
	}
	return categories, nil
}

// GetCategoryHierarchy gets categories in hierarchical structure
// ADD: Missing method
func (s *CategoryService) GetCategoryHierarchy() ([]models.Category, error) {
	categories, err := s.categoryRepo.GetCategoryHierarchy()
	if err != nil {
		utils.Error("CategoryService.GetCategoryHierarchy: Failed - Error=%v", err)
		return nil, fmt.Errorf("failed to get category hierarchy")
	}
	return categories, nil
}

// Search searches categories by name
// ADD: Missing method
func (s *CategoryService) Search(keyword string) ([]models.Category, error) {
	if keyword == "" {
		return nil, fmt.Errorf("search keyword is required")
	}

	categories, err := s.categoryRepo.Search(keyword)
	if err != nil {
		utils.Error("CategoryService.Search: Failed - Keyword=%s, Error=%v", keyword, err)
		return nil, fmt.Errorf("failed to search categories")
	}
	return categories, nil
}

// GetAllWithPagination gets all categories with pagination
// ADD: Missing method
func (s *CategoryService) GetAllWithPagination(page, limit int) ([]models.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	categories, total, err := s.categoryRepo.GetAllWithPagination(limit, offset)
	if err != nil {
		utils.Error("CategoryService.GetAllWithPagination: Failed - Error=%v", err)
		return nil, 0, fmt.Errorf("failed to get categories")
	}

	return categories, total, nil
}

// GetCategoryStats gets category statistics
// FIX: Simple implementation without repository methods (calculate from GetAll)
func (s *CategoryService) GetCategoryStats() (map[string]interface{}, error) {
	// Get all categories (including inactive)
	allCategories, err := s.categoryRepo.GetAllIncludeInactive()
	if err != nil {
		utils.Error("CategoryService.GetCategoryStats: Failed - Error=%v", err)
		return nil, fmt.Errorf("failed to get category statistics")
	}

	// Calculate stats manually
	var totalCategories int64
	var activeCategories int64
	var inactiveCategories int64

	for _, category := range allCategories {
		totalCategories++
		if category.IsActive {
			activeCategories++
		} else {
			inactiveCategories++
		}
	}

	stats := map[string]interface{}{
		"total_categories":    totalCategories,
		"active_categories":   activeCategories,
		"inactive_categories": inactiveCategories,
	}

	return stats, nil
}