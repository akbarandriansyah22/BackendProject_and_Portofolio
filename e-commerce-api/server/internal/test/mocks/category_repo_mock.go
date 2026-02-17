package mocks

import (
	"context"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

type MockCategoryRepository struct {
	Categories map[int]*models.Category
}

func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		Categories: make(map[int]*models.Category),
	}
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *models.Category) error {
	m.Categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	if c, ok := m.Categories[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockCategoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	for _, c := range m.Categories {
		if c.Slug == slug {
			return c, nil
		}
	}
	return nil, nil
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	m.Categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id int) error {
	delete(m.Categories, id)
	return nil
}

func (m *MockCategoryRepository) List(ctx context.Context) ([]*models.Category, error) {
	categories := make([]*models.Category, 0, len(m.Categories))
	for _, c := range m.Categories {
		categories = append(categories, c)
	}
	return categories, nil
}

func (m *MockCategoryRepository) GetChildren(ctx context.Context, parentID int) ([]*models.Category, error) {
	return make([]*models.Category, 0), nil
}

func (m *MockCategoryRepository) GetRootCategories(ctx context.Context) ([]*models.Category, error) {
	return make([]*models.Category, 0), nil
}

func (m *MockCategoryRepository) HasProducts(ctx context.Context, categoryID int) (bool, error) {
	return false, nil
}

func (m *MockCategoryRepository) SetActive(ctx context.Context, categoryID int, isActive bool) error {
	return nil
}

func (m *MockCategoryRepository) AssignProductToCategory(ctx context.Context, productID, categoryID int) error {
	return nil
}

func (m *MockCategoryRepository) RemoveProductFromCategory(ctx context.Context, productID, categoryID int) error {
	return nil
}

func (m *MockCategoryRepository) GetProductCategories(ctx context.Context, productID int) ([]*models.Category, error) {
	return make([]*models.Category, 0), nil
}
