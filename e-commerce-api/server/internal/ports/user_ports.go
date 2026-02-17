package ports

import (
	"context"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

// UserRepository mendefinisikan kontrak untuk data access layer User
type UserRepository interface {
	// Create membuat user baru
	Create(ctx context.Context, user *models.User) error

	// GetByID mengambil user berdasarkan ID
	GetByID(ctx context.Context, id int) (*models.User, error)

	// GetByEmail mengambil user berdasarkan email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// Update memperbarui data user
	Update(ctx context.Context, user *models.User) error

	// Delete menghapus user (soft delete)
	Delete(ctx context.Context, id int) error

	// List mengambil daftar user dengan pagination
	List(ctx context.Context, page, limit int) ([]*models.User, int, error)

	// UpdatePassword memperbarui password user
	UpdatePassword(ctx context.Context, userID int, passwordHash string) error

	// VerifyEmail memverifikasi email user
	VerifyEmail(ctx context.Context, userID int) error

	// SetActive mengatur status aktif user
	SetActive(ctx context.Context, userID int, isActive bool) error
}

// UserService mendefinisikan kontrak untuk business logic layer User
type UserService interface {
	// GetProfile mengambil profile user yang sedang login
	GetProfile(ctx context.Context, userID int) (*models.UserProfileResponse, error)

	// UpdateProfile memperbarui profile user
	UpdateProfile(ctx context.Context, userID int, req *models.UpdateProfileRequest) error

	// ChangePassword mengubah password user
	ChangePassword(ctx context.Context, userID int, req *models.ChangePasswordRequest) error

	// GetUserByID mengambil user berdasarkan ID (admin only)
	GetUserByID(ctx context.Context, id int) (*models.UserResponse, error)

	// ListUsers mengambil daftar semua user (admin only)
	ListUsers(ctx context.Context, page, limit int) (*models.PaginatedResponse, error)

	// DeactivateUser menonaktifkan user (admin only)
	DeactivateUser(ctx context.Context, userID int) error

	// ActivateUser mengaktifkan user (admin only)
	ActivateUser(ctx context.Context, userID int) error
}
