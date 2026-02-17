package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/service"
	testmocks "github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/test/mocks"
)


func TestGetAllProducts(t *testing.T) {
	productRepo := testmocks.NewMockProductRepository()
	categoryRepo := testmocks.NewMockCategoryRepository()
	logger := observability.NewLogger()

	svc := service.NewProductService(
		productRepo,
		categoryRepo,
		logger,
	)

	// Verify service was created successfully
	assert.NotNil(t, svc)
}

func TestProductService_ListProducts_Success(t *testing.T) {
	// arrange
	ctx := context.Background()

	productRepoMock := testmocks.NewMockProductRepository()
	categoryRepoMock := testmocks.NewMockCategoryRepository()
	logger := observability.NewLogger()

	productService := service.NewProductService(productRepoMock, categoryRepoMock, logger)

	filter := &models.ProductFilter{
		Page:  1,
		Limit: 10,
	}

	products := []*models.Product{
		{ID: 1, Name: "Product A", Price: 10000, Stock: 5},
		{ID: 2, Name: "Product B", Price: 20000, Stock: 10},
	}

	productRepoMock.
		On("List", ctx, filter).
		Return(products, 2, nil)

	// act
	resp, err := productService.ListProducts(ctx, filter)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.Limit)

	productRepoMock.AssertExpectations(t)
}

