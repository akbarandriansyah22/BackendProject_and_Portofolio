package service

import (
	"context"
	"fmt"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/ports"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo ports.CategoryRepository
	productRepo  ports.ProductRepository
	logger       observability.Logger
}

// NewCategoryService creates a new category service
func NewCategoryService(
	categoryRepo ports.CategoryRepository,
	productRepo ports.ProductRepository,
	logger observability.Logger,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
		logger:       logger,
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(ctx context.Context, req *models.CreateCategoryRequest) (*models.CategoryResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("category name is required")
	}
	if req.Slug == "" {
		return nil, fmt.Errorf("category slug is required")
	}

	category := &models.Category{
		Name:     req.Name,
		Slug:     req.Slug,
		IsActive: true,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		s.logger.Error("CategoryService.CreateCategory failed", err)
		return nil, fmt.Errorf("failed to create category")
	}

	s.logger.Info("Category created: ID=%d, Name=%s", category.ID, category.Name)

	return &models.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}, nil
}

// GetCategoryByID gets a category by ID
func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*models.CategoryDetailResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category")
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	return &models.CategoryDetailResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}, nil
}

// GetCategoryBySlug gets a category by slug
func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*models.CategoryDetailResponse, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get category")
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	return &models.CategoryDetailResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}, nil
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(ctx context.Context, id int, req *models.UpdateCategoryRequest) error {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil || category == nil {
		return fmt.Errorf("category not found")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		s.logger.Error("CategoryService.UpdateCategory failed", err)
		return fmt.Errorf("failed to update category")
	}

	s.logger.Info("Category updated: ID=%d, Name=%s", id, category.Name)
	return nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil || category == nil {
		return fmt.Errorf("category not found")
	}

	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		s.logger.Error("CategoryService.DeleteCategory failed", err)
		return fmt.Errorf("failed to delete category")
	}

	s.logger.Info("Category deleted: ID=%d", id)
	return nil
}

// ListCategories lists all categories
func (s *CategoryService) ListCategories(ctx context.Context) ([]*models.CategoryResponse, error) {
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		s.logger.Error("CategoryService.ListCategories failed", err)
		return nil, fmt.Errorf("failed to list categories")
	}

	response := make([]*models.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		response = append(response, &models.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		})
	}

	return response, nil
}

// GetCategoryTree gets categories in tree structure
func (s *CategoryService) GetCategoryTree(ctx context.Context) ([]*models.CategoryTreeResponse, error) {
	rootCategories, err := s.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		s.logger.Error("CategoryService.GetCategoryTree failed", err)
		return nil, fmt.Errorf("failed to get category tree")
	}

	response := make([]*models.CategoryTreeResponse, 0, len(rootCategories))
	for _, category := range rootCategories {
		children, _ := s.categoryRepo.GetChildren(ctx, category.ID)
		treeNode := &models.CategoryTreeResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		}

		if len(children) > 0 {
			childResponses := make([]*models.CategoryTreeResponse, 0, len(children))
			for _, child := range children {
				childResponses = append(childResponses, &models.CategoryTreeResponse{
					ID:   child.ID,
					Name: child.Name,
					Slug: child.Slug,
				})
			}
			treeNode.Children = childResponses
		}

		response = append(response, treeNode)
	}

	return response, nil
}

// GetRootCategories gets root categories
func (s *CategoryService) GetRootCategories(ctx context.Context) ([]*models.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		s.logger.Error("CategoryService.GetRootCategories failed", err)
		return nil, fmt.Errorf("failed to get root categories")
	}

	response := make([]*models.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		response = append(response, &models.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		})
	}

	return response, nil
}

// GetChildCategories gets child categories of a parent
func (s *CategoryService) GetChildCategories(ctx context.Context, parentID int) ([]*models.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetChildren(ctx, parentID)
	if err != nil {
		s.logger.Error("CategoryService.GetChildCategories failed", err)
		return nil, fmt.Errorf("failed to get child categories")
	}

	response := make([]*models.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		response = append(response, &models.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		})
	}

	return response, nil
}

// ActivateCategory activates a category
func (s *CategoryService) ActivateCategory(ctx context.Context, categoryID int) error {
	return s.categoryRepo.SetActive(ctx, categoryID, true)
}

// DeactivateCategory deactivates a category
func (s *CategoryService) DeactivateCategory(ctx context.Context, categoryID int) error {
	return s.categoryRepo.SetActive(ctx, categoryID, false)
}

// AssignProductToCategory assigns a product to a category
func (s *CategoryService) AssignProductToCategory(ctx context.Context, productID, categoryID int) error {
	return s.categoryRepo.AssignProductToCategory(ctx, productID, categoryID)
}

// RemoveProductFromCategory removes a product from a category
func (s *CategoryService) RemoveProductFromCategory(ctx context.Context, productID, categoryID int) error {
	return s.categoryRepo.RemoveProductFromCategory(ctx, productID, categoryID)
}

// ============================================
// Handler-level methods (non-context versions)
// ============================================

// GetAll gets all categories (handler wrapper)
func (s *CategoryService) GetAll() ([]*models.CategoryResponse, error) {
	return s.ListCategories(context.Background())
}

// GetByID gets category by ID (handler wrapper)
func (s *CategoryService) GetByID(id int) (*models.CategoryDetailResponse, error) {
	return s.GetCategoryByID(context.Background(), id)
}

// GetProductsByCategory gets products for a category with pagination
func (s *CategoryService) GetProductsByCategory(categoryID, page, limit int) ([]*models.Product, int, error) {
	ctx := context.Background()

	// Verify category exists
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil || category == nil {
		return nil, 0, fmt.Errorf("category not found")
	}

	// Get products
	products, total, err := s.productRepo.GetProductsByCategory(ctx, categoryID, page, limit)
	if err != nil {
		s.logger.Error("CategoryService.GetProductsByCategory failed", err)
		return nil, 0, fmt.Errorf("failed to get products")
	}

	return products, total, nil
}

// GetSubCategories gets subcategories of a category
func (s *CategoryService) GetSubCategories(parentID int) ([]*models.CategoryResponse, error) {
	return s.GetChildCategories(context.Background(), parentID)
}

// Create creates a new category (handler wrapper)
func (s *CategoryService) Create(req *models.CreateCategoryRequest) (*models.CategoryResponse, error) {
	return s.CreateCategory(context.Background(), req)
}

// Update updates a category (handler wrapper)
func (s *CategoryService) Update(id int, req *models.UpdateCategoryRequest) (*models.CategoryDetailResponse, error) {
	ctx := context.Background()
	if err := s.UpdateCategory(ctx, id, req); err != nil {
		return nil, err
	}
	return s.GetCategoryByID(ctx, id)
}

// Delete deletes a category (handler wrapper)
func (s *CategoryService) Delete(id int) error {
	return s.DeleteCategory(context.Background(), id)
}

// ToggleStatus toggles category status (handler wrapper)
func (s *CategoryService) ToggleStatus(id int, isActive bool) error {
	ctx := context.Background()
	if isActive {
		return s.ActivateCategory(ctx, id)
	}
	return s.DeactivateCategory(ctx, id)
}

// GetCategoryStats gets category statistics
func (s *CategoryService) GetCategoryStats() (interface{}, error) {
	ctx := context.Background()
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		s.logger.Error("CategoryService.GetCategoryStats failed", err)
		return nil, fmt.Errorf("failed to get category statistics")
	}

	stats := map[string]interface{}{
		"total_categories": len(categories),
		"active":           0,
		"inactive":         0,
	}

	for _, cat := range categories {
		if cat.IsActive {
			stats["active"] = stats["active"].(int) + 1
		} else {
			stats["inactive"] = stats["inactive"].(int) + 1
		}
	}

	return stats, nil
}
