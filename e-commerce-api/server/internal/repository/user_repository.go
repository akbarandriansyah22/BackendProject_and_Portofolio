package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	models "github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

type UserRepository struct {
	db *sql.DB
	logger *zap.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db: db,
	logger: logger,}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (role_id, email, password_hash, full_name, phone, address, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.RoleID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Phone,
		user.Address,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, role_id, email, password_hash, full_name, phone, address, 
		       is_active, email_verified_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Phone,
		&user.Address,
		&user.IsActive,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	return user, err
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, role_id, email, password_hash, full_name, phone, address, 
		       is_active, email_verified_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.RoleID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Phone,
		&user.Address,
		&user.IsActive,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	return user, err
}

// GetAll retrieves all users with pagination
func (r *UserRepository) GetAll(limit, offset int) ([]models.User, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM users`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get users
	query := `
		SELECT id, role_id, email, password_hash, full_name, phone, address, 
		       is_active, email_verified_at, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.RoleID,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&user.Phone,
			&user.Address,
			&user.IsActive,
			&user.EmailVerifiedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// GetByRole retrieves users by role ID
func (r *UserRepository) GetByRole(roleID int, limit, offset int) ([]models.User, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE role_id = $1`
	if err := r.db.QueryRow(countQuery, roleID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get users
	query := `
		SELECT id, role_id, email, password_hash, full_name, phone, address, 
		       is_active, email_verified_at, created_at, updated_at
		FROM users
		WHERE role_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, roleID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.RoleID,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&user.Phone,
			&user.Address,
			&user.IsActive,
			&user.EmailVerifiedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// Update updates user information
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET full_name = $1, phone = $2, address = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	return r.db.QueryRow(
		query,
		user.Name,
		user.Phone,
		user.Address,
		user.ID,
	).Scan(&user.UpdatedAt)
}

// UpdatePassword updates user password
func (r *UserRepository) UpdatePassword(id int, hashedPassword string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, hashedPassword, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// VerifyEmail marks user email as verified
func (r *UserRepository) VerifyEmail(id int) error {
	query := `
		UPDATE users
		SET email_verified_at = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// SetActive activates or deactivates a user
func (r *UserRepository) SetActive(id int, isActive bool) error {
	query := `
		UPDATE users
		SET is_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, isActive, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user by ID (hard delete)
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// EmailExistsExcludingUser checks if email exists for another user (for update validation)
func (r *UserRepository) EmailExistsExcludingUser(email string, userID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)`
	var exists bool
	err := r.db.QueryRow(query, email, userID).Scan(&exists)
	return exists, err
}

// CountByRole counts users by role
func (r *UserRepository) CountByRole(roleID int) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE role_id = $1`
	var count int64
	err := r.db.QueryRow(query, roleID).Scan(&count)
	return count, err
}

// CountActive counts active users
func (r *UserRepository) CountActive() (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE is_active = true`
	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// SearchByName searches users by name (for admin panel)
func (r *UserRepository) SearchByName(searchTerm string, limit, offset int) ([]models.User, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE full_name ILIKE $1`
	searchPattern := "%" + searchTerm + "%"
	if err := r.db.QueryRow(countQuery, searchPattern).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get users
	query := `
		SELECT id, role_id, email, password_hash, full_name, phone, address, 
		       is_active, email_verified_at, created_at, updated_at
		FROM users
		WHERE full_name ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.RoleID,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&user.Phone,
			&user.Address,
			&user.IsActive,
			&user.EmailVerifiedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}
