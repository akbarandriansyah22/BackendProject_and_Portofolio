package service

import (
	"context"
	"errors"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/ports"
)
type ProductService struct {
	productRepo  ports.ProductRepository
	categoryRepo ports.CategoryRepository
	logger       observability.Logger
}


func NewProductService(
	productRepo ports.ProductRepository,
	categoryRepo ports.CategoryRepository,
	logger observability.Logger,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

// =======================
// CREATE
// =======================
func (s *ProductService) CreateProduct(
	ctx context.Context,
	req *models.CreateProductRequest,
) (*models.ProductResponse, error) {

	product := &models.Product{
		Name:     req.Name,
		Slug:     req.Slug,
		Price:    req.Price,
		Stock:    req.Stock,
		IsActive: true,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		s.logger.Error("failed to create product", err)
		return nil, err
	}

	return &models.ProductResponse{
		ID:    product.ID,
		Name:  product.Name,
		Slug:  product.Slug,
		Price: product.Price,
		Stock: product.Stock,
	}, nil
}

// =======================
// READ
// =======================
func (s *ProductService) GetProductByID(
	ctx context.Context,
	id int,
) (*models.ProductDetailResponse, error) {

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	return mapToDetailResponse(product), nil
}

func (s *ProductService) GetProductBySlug(
	ctx context.Context,
	slug string,
) (*models.ProductDetailResponse, error) {

	product, err := s.productRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	return mapToDetailResponse(product), nil
}

// =======================
// UPDATE
// =======================
func (s *ProductService) UpdateProduct(
	ctx context.Context,
	id int,
	req *models.UpdateProductRequest,
) error {

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	product.Name = req.Name
	product.Price = req.Price
	product.Stock = req.Stock

	return s.productRepo.Update(ctx, product)
}

// =======================
// DELETE
// =======================
func (s *ProductService) DeleteProduct(
	ctx context.Context,
	id int,
) error {
	return s.productRepo.Delete(ctx, id)
}

// =======================
// LIST & SEARCH
// =======================
func (s *ProductService) ListProducts(
	ctx context.Context,
	filter *models.ProductFilter,
) (*models.PaginatedResponse, error) {

	products, total, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return buildPaginatedResponse(products, total, filter.Page, filter.Limit), nil
}

func (s *ProductService) SearchProducts(
	ctx context.Context,
	keyword string,
	page, limit int,
) (*models.PaginatedResponse, error) {

	products, total, err := s.productRepo.Search(ctx, keyword, page, limit)
	if err != nil {
		return nil, err
	}

	return buildPaginatedResponse(products, total, page, limit), nil
}

// =======================
// STOCK
// =======================
func (s *ProductService) UpdateStock(
	ctx context.Context,
	productID int,
	quantity int,
) error {
	return s.productRepo.UpdateStock(ctx, productID, quantity)
}

// =======================
// STATUS
// =======================
func (s *ProductService) ActivateProduct(
	ctx context.Context,
	productID int,
) error {
	return s.productRepo.SetActive(ctx, productID, true)
}

func (s *ProductService) DeactivateProduct(
	ctx context.Context,
	productID int,
) error {
	return s.productRepo.SetActive(ctx, productID, false)
}

// =======================
// CATEGORY
// =======================
func (s *ProductService) GetProductsByCategory(
	ctx context.Context,
	categoryID, page, limit int,
) (*models.PaginatedResponse, error) {

	products, total, err := s.productRepo.GetProductsByCategory(ctx, categoryID, page, limit)
	if err != nil {
		return nil, err
	}

	return buildPaginatedResponse(products, total, page, limit), nil
}

// =======================
// HELPERS (PRIVATE)
// =======================
func mapToDetailResponse(p *models.Product) *models.ProductDetailResponse {
	return &models.ProductDetailResponse{
		ID:       p.ID,
		Name:     p.Name,
		Slug:     p.Slug,
		Price:    p.Price,
		Stock:    p.Stock,
		IsActive: p.IsActive,
	}
}

func buildPaginatedResponse(
	products []*models.Product,
	total int,
	page, limit int,
) *models.PaginatedResponse {

	items := make([]interface{}, 0, len(products))
	for _, p := range products {
		items = append(items, &models.ProductResponse{
			ID:    p.ID,
			Name:  p.Name,
			Slug:  p.Slug,
			Price: p.Price,
			Stock: p.Stock,
		})
	}

	return &models.PaginatedResponse{
		Data:  items,
		Total: int64(total),
		Page:  page,
		Limit: limit,
	}
}
