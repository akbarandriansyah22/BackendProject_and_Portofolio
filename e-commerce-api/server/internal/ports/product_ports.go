package ports

import (
	"context"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

// ProductRepository mendefinisikan kontrak untuk data access layer Product
type ProductRepository interface {
	// Create membuat product baru
	Create(ctx context.Context, product *models.Product) error

	// GetByID mengambil product berdasarkan ID
	GetByID(ctx context.Context, id int) (*models.Product, error)

	// GetBySlug mengambil product berdasarkan slug
	GetBySlug(ctx context.Context, slug string) (*models.Product, error)

	// Update memperbarui data product
	Update(ctx context.Context, product *models.Product) error

	// Delete menghapus product (soft delete)
	Delete(ctx context.Context, id int) error

	// List mengambil daftar product dengan filter dan pagination
	List(ctx context.Context, filter *models.ProductFilter) ([]*models.Product, int, error)

	// UpdateStock memperbarui stok product
	UpdateStock(ctx context.Context, productID, quantity int) error

	// DecrementStock mengurangi stok product
	DecrementStock(ctx context.Context, productID, quantity int) error

	// IncrementStock menambah stok product
	IncrementStock(ctx context.Context, productID, quantity int) error

	// CheckStock memeriksa ketersediaan stok
	CheckStock(ctx context.Context, productID, quantity int) (bool, error)

	// GetProductsByCategory mengambil product berdasarkan kategori
	GetProductsByCategory(ctx context.Context, categoryID int, page, limit int) ([]*models.Product, int, error)

	// Search mencari product berdasarkan keyword
	Search(ctx context.Context, keyword string, page, limit int) ([]*models.Product, int, error)

	// SetActive mengatur status aktif product
	SetActive(ctx context.Context, productID int, isActive bool) error
}

// ProductService mendefinisikan kontrak untuk business logic layer Product
type ProductService interface {
	// CreateProduct membuat product baru (admin only)
	CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.ProductResponse, error)

	// GetProductByID mengambil detail product
	GetProductByID(ctx context.Context, id int) (*models.ProductDetailResponse, error)

	// GetProductBySlug mengambil detail product berdasarkan slug
	GetProductBySlug(ctx context.Context, slug string) (*models.ProductDetailResponse, error)

	// UpdateProduct memperbarui product (admin only)
	UpdateProduct(ctx context.Context, id int, req *models.UpdateProductRequest) error

	// DeleteProduct menghapus product (admin only)
	DeleteProduct(ctx context.Context, id int) error

	// ListProducts mengambil daftar product dengan filter
	ListProducts(ctx context.Context, filter *models.ProductFilter) (*models.PaginatedResponse, error)

	// SearchProducts mencari product
	SearchProducts(ctx context.Context, keyword string, page, limit int) (*models.PaginatedResponse, error)

	// UpdateStock memperbarui stok product (admin only)
	UpdateStock(ctx context.Context, productID int, quantity int) error

	// ActivateProduct mengaktifkan product (admin only)
	ActivateProduct(ctx context.Context, productID int) error

	// DeactivateProduct menonaktifkan product (admin only)
	DeactivateProduct(ctx context.Context, productID int) error

	// GetProductsByCategory mengambil product berdasarkan kategori
	GetProductsByCategory(ctx context.Context, categoryID, page, limit int) (*models.PaginatedResponse, error)
}
