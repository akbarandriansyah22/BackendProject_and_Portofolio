package repository

import (
	"database/sql"
	"fmt"
	"log"

	"go.uber.org/zap"

	models "github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

type RoleRepository struct {
	db *sql.DB
	logger *zap.Logger
}

func NewRoleRepository(db *sql.DB, logger *zap.Logger) *RoleRepository {
	return &RoleRepository{
		db: db,
	logger: logger,}
}

// GetAll retrieves all roles
func (r *RoleRepository) GetAll() ([]models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	roles := []models.Role{}
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(id int) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	role := &models.Role{}
	err := r.db.QueryRow(query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("role not found")
	}

	return role, err
}

// GetByName retrieves a role by name
func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	role := &models.Role{}
	err := r.db.QueryRow(query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("role not found")
	}

	return role, err
}

// Create creates a new role
func (r *RoleRepository) Create(role *models.Role) error {
	query := `
		INSERT INTO roles (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, role.Name, role.Description).Scan(
		&role.ID,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	return err
}

// Update updates an existing role
func (r *RoleRepository) Update(role *models.Role) error {
	query := `
		UPDATE roles
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING updated_at
	`

	return r.db.QueryRow(query, role.Name, role.Description, role.ID).Scan(&role.UpdatedAt)
}

// Delete deletes a role by ID
func (r *RoleRepository) Delete(id int) error {
	query := `DELETE FROM roles WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// NameExists checks if a role name already exists
func (r *RoleRepository) NameExists(name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)`
	var exists bool
	err := r.db.QueryRow(query, name).Scan(&exists)
	return exists, err
}
