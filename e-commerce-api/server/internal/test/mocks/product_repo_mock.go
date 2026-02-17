package mocks

import (
	"context"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func NewMockProductRepository() *ProductRepositoryMock {
	return &ProductRepositoryMock{}
}

func (m *ProductRepositoryMock) Create(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) GetByID(ctx context.Context, id int) (*models.Product, error) {
	args := m.Called(ctx, id)
	if p := args.Get(0); p != nil {
		return p.(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	args := m.Called(ctx, slug)
	if p := args.Get(0); p != nil {
		return p.(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) Update(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) List(ctx context.Context, filter *models.ProductFilter) ([]*models.Product, int, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*models.Product), args.Int(1), args.Error(2)
}

func (m *ProductRepositoryMock) UpdateStock(ctx context.Context, productID, quantity int) error {
	return m.Called(ctx, productID, quantity).Error(0)
}

func (m *ProductRepositoryMock) DecrementStock(ctx context.Context, productID, quantity int) error {
	return m.Called(ctx, productID, quantity).Error(0)
}

func (m *ProductRepositoryMock) IncrementStock(ctx context.Context, productID, quantity int) error {
	return m.Called(ctx, productID, quantity).Error(0)
}

func (m *ProductRepositoryMock) CheckStock(ctx context.Context, productID, quantity int) (bool, error) {
	args := m.Called(ctx, productID, quantity)
	return args.Bool(0), args.Error(1)
}

func (m *ProductRepositoryMock) GetProductsByCategory(ctx context.Context, categoryID, page, limit int) ([]*models.Product, int, error) {
	args := m.Called(ctx, categoryID, page, limit)
	return args.Get(0).([]*models.Product), args.Int(1), args.Error(2)
}

func (m *ProductRepositoryMock) Search(ctx context.Context, keyword string, page, limit int) ([]*models.Product, int, error) {
	args := m.Called(ctx, keyword, page, limit)
	return args.Get(0).([]*models.Product), args.Int(1), args.Error(2)
}

func (m *ProductRepositoryMock) SetActive(ctx context.Context, productID int, isActive bool) error {
	return m.Called(ctx, productID, isActive).Error(0)
}
