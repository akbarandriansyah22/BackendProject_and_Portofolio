package mocks

import (
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) GetByID(id int) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) SlugExists(slug string) (bool, error) {
	args := m.Called(slug)
	return args.Bool(0), args.Error(1)
}

func (m *ProductRepositoryMock) SKUExists(sku string) (bool, error) {
	args := m.Called(sku)
	return args.Bool(0), args.Error(1)
}

func (m *ProductRepositoryMock) Update(product *models.Product) error {
	return m.Called(product).Error(0)
}

func (m *ProductRepositoryMock) Delete(id int) error {
	return m.Called(id).Error(0)
}

func (m *ProductRepositoryMock) IncrementStock(id, qty int) error {
	return m.Called(id, qty).Error(0)
}

func (m *ProductRepositoryMock) DecrementStock(id, qty int) error {
	return m.Called(id, qty).Error(0)
}

func (m *ProductRepositoryMock) UpdateStock(id, stock int) error {
	return m.Called(id, stock).Error(0)
}
