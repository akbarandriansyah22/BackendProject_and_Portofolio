package ports

import (
	"context"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

// RoleRepository mendefinisikan kontrak untuk data access layer Role
type RoleRepository interface {
	// Create membuat role baru
	Create(ctx context.Context, role *models.Role) error

	// GetByID mengambil role berdasarkan ID
	GetByID(ctx context.Context, id int) (*models.Role, error)

	// GetByName mengambil role berdasarkan nama
	GetByName(ctx context.Context, name string) (*models.Role, error)

	// Update memperbarui data role
	Update(ctx context.Context, role *models.Role) error

	// Delete menghapus role
	Delete(ctx context.Context, id int) error

	// List mengambil semua role
	List(ctx context.Context) ([]*models.Role, error)
}

// RoleService mendefinisikan kontrak untuk business logic layer Role
type RoleService interface {
	// CreateRole membuat role baru (admin only)
	CreateRole(ctx context.Context, req *models.CreateRoleRequest) (*models.RoleResponse, error)

	// GetRoleByID mengambil role berdasarkan ID
	GetRoleByID(ctx context.Context, id int) (*models.RoleResponse, error)

	// UpdateRole memperbarui role (admin only)
	UpdateRole(ctx context.Context, id int, req *models.UpdateRoleRequest) error

	// DeleteRole menghapus role (admin only)
	DeleteRole(ctx context.Context, id int) error

	// ListRoles mengambil semua role
	ListRoles(ctx context.Context) ([]*models.RoleResponse, error)
}
