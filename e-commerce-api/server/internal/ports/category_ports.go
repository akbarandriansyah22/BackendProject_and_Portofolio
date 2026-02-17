package ports

import (
	"context"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

// CategoryRepository mendefinisikan kontrak untuk data access layer Category
type CategoryRepository interface {
	// Create membuat category baru
	Create(ctx context.Context, category *models.Category) error

	// GetByID mengambil category berdasarkan ID
	GetByID(ctx context.Context, id int) (*models.Category, error)

	// GetBySlug mengambil category berdasarkan slug
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)

	// Update memperbarui data category
	Update(ctx context.Context, category *models.Category) error

	// Delete menghapus category (soft delete)
	Delete(ctx context.Context, id int) error

	// List mengambil semua category
	List(ctx context.Context) ([]*models.Category, error)

	// GetChildren mengambil child categories dari parent
	GetChildren(ctx context.Context, parentID int) ([]*models.Category, error)

	// GetRootCategories mengambil kategori root (tanpa parent)
	GetRootCategories(ctx context.Context) ([]*models.Category, error)

	// HasProducts memeriksa apakah category memiliki product
	HasProducts(ctx context.Context, categoryID int) (bool, error)

	// SetActive mengatur status aktif category
	SetActive(ctx context.Context, categoryID int, isActive bool) error

	// AssignProductToCategory menghubungkan product dengan category
	AssignProductToCategory(ctx context.Context, productID, categoryID int) error

	// RemoveProductFromCategory memutus hubungan product dengan category
	RemoveProductFromCategory(ctx context.Context, productID, categoryID int) error

	// GetProductCategories mengambil categories dari product
	GetProductCategories(ctx context.Context, productID int) ([]*models.Category, error)
}

// CategoryService mendefinisikan kontrak untuk business logic layer Category
type CategoryService interface {
	// CreateCategory membuat category baru (admin only)
	CreateCategory(ctx context.Context, req *models.CreateCategoryRequest) (*models.CategoryResponse, error)

	// GetCategoryByID mengambil detail category
	GetCategoryByID(ctx context.Context, id int) (*models.CategoryDetailResponse, error)

	// GetCategoryBySlug mengambil detail category berdasarkan slug
	GetCategoryBySlug(ctx context.Context, slug string) (*models.CategoryDetailResponse, error)

	// UpdateCategory memperbarui category (admin only)
	UpdateCategory(ctx context.Context, id int, req *models.UpdateCategoryRequest) error

	// DeleteCategory menghapus category (admin only)
	DeleteCategory(ctx context.Context, id int) error

	// ListCategories mengambil semua category
	ListCategories(ctx context.Context) ([]*models.CategoryResponse, error)

	// GetCategoryTree mengambil category dalam bentuk tree/hierarchy
	GetCategoryTree(ctx context.Context) ([]*models.CategoryTreeResponse, error)

	// GetRootCategories mengambil kategori root
	GetRootCategories(ctx context.Context) ([]*models.CategoryResponse, error)

	// GetChildCategories mengambil child categories
	GetChildCategories(ctx context.Context, parentID int) ([]*models.CategoryResponse, error)

	// ActivateCategory mengaktifkan category (admin only)
	ActivateCategory(ctx context.Context, categoryID int) error

	// DeactivateCategory menonaktifkan category (admin only)
	DeactivateCategory(ctx context.Context, categoryID int) error

	// AssignProductToCategory menghubungkan product dengan category (admin only)
	AssignProductToCategory(ctx context.Context, productID, categoryID int) error

	// RemoveProductFromCategory memutus hubungan product dengan category (admin only)
	RemoveProductFromCategory(ctx context.Context, productID, categoryID int) error
}
