package service

import (
	"database/sql"
	"fmt"

	"e-commerce-api/server/internal/models"
	"e-commerce-api/server/internal/repository"
	"e-commerce-api/server/internal/utils"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
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
func (s *CategoryService) ToggleStatus(id int, isActive bool) error {
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("category not found")
	}

	if err := s.categoryRepo.UpdateStatus(id, isActive); err != nil {
		return fmt.Errorf("failed to update category status")
	}

	utils.Info("Category status toggled: ID=%d, IsActive=%v", id, isActive)
	return nil
}

// GetProductsByCategory gets products for a category
func (s *CategoryService) GetProductsByCategory(categoryID, page, limit int) ([]models.Product, int64, error) {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("category not found")
	}

	offset := (page - 1) * limit
	return s.categoryRepo.GetProducts(categoryID, limit, offset)
}

// GetSubCategories gets sub-categories of a category
func (s *CategoryService) GetSubCategories(parentID int) ([]models.Category, error) {
	categories, err := s.categoryRepo.GetByParentID(parentID)
	if err != nil {
		utils.Error("CategoryService.GetSubCategories: Failed - ParentID=%d, Error=%v", parentID, err)
		return nil, fmt.Errorf("failed to get sub-categories")
	}
	return categories, nil
}

// GetCategoryStats gets category statistics
func (s *CategoryService) GetCategoryStats() (map[string]interface{}, error) {
	totalCategories, err := s.categoryRepo.CountTotal()
	if err != nil {
		return nil, err
	}

	activeCategories, err := s.categoryRepo.CountByStatus(true)
	if err != nil {
		return nil, err
	}

	inactiveCategories, err := s.categoryRepo.CountByStatus(false)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_categories":   totalCategories,
		"active_categories":  activeCategories,
		"inactive_categories": inactiveCategories,
	}

	return stats, nil
}

