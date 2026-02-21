package repository

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
	logger *zap.Logger
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB, logger *zap.Logger) *CategoryRepository {
	return &CategoryRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	query := `
		INSERT INTO categories (name, slug, description, parent_id, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.IsActive,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	return err
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	category := &models.Category{}
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ParentID,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}

	return category, err
}

// GetBySlug retrieves a category by slug
func (r *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, is_active, created_at, updated_at
		FROM categories
		WHERE slug = $1
	`

	category := &models.Category{}
	err := r.db.QueryRow(query, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ParentID,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}

	return category, err
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		r.logger.Warn("rows close error", zap.Error(err))

    }
}()

	category := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&c.Description,
			&c.ParentID,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		category = append(category, c)
	}

	return category, nil
}

// GetAllWithPagination retrieves all categories with pagination
func (r *CategoryRepository) GetAllWithPagination(limit, offset int) ([]models.Category, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM categories WHERE is_active = true`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get categories
	query := `
		SELECT id, name, slug, description, parent_id, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
    if err := rows.Close(); err != nil {
		r.logger.Warn("rows close error", zap.Error(err))
    }
}()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&c.Description,
			&c.ParentID,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		categories = append(categories, c)
	}

	return categories, total, nil
}

// GetAllIncludeInactive retrieves all categories including inactive (for admin)
func (r *CategoryRepository) GetAllIncludeInactive() ([]models.Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, is_active, created_at, updated_at
		FROM categories
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		r.logger.Warn("rows close error", zap.Error(err))
    }
}()

	category := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&c.Description,
			&c.ParentID,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		category = append(category, c)
	}

	return category, nil
}

// GetParentCategories retrieves all parent categories (no parent_id)
func (r *CategoryRepository) GetParentCategories() ([]models.Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, is_active, created_at, updated_at
		FROM categories
		WHERE parent_id IS NULL AND is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		r.logger.Warn("rows close error", zap.Error(err))
    }
}()

	category := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&c.Description,
			&c.ParentID,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		category = append(category, c)
	}

	return category, nil
}

// GetSubCategories retrieves subcategories of a parent category
func (r *CategoryRepository) GetSubCategories(parentID int) ([]models.Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, is_active, created_at, updated_at
		FROM categories
		WHERE parent_id = $1 AND is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		r.logger.Warn("rows close error", zap.Error(err))
    }
}()

	category := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&c.Description,
			&c.ParentID,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		category = append(category, c)
	}

	return category, nil
}

// ADDED: GetByParentID - alias for GetSubCategories (dipanggil di category_service.go)
func (r *CategoryRepository) GetByParentID(parentID int) ([]models.Category, error) {
	return r.GetSubCategories(parentID)
}

// GetCategoryHierarchy retrieves categories in hierarchical structure
func (r *CategoryRepository) GetCategoryHierarchy() ([]models.Category, error) {
	// First get all parent categories
	parents, err := r.GetParentCategories()
	if err != nil {
		return nil, err
	}

	// Then get subcategories for each parent
	// Note: This returns flat list, frontend should organize into tree
	allCategories := parents
	for _, parent := range parents {
		subs, err := r.GetSubCategories(parent.ID)
		if err != nil {
			continue
		}
		allCategories = append(allCategories, subs...)
	}

	return allCategories, nil
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	query := `
		UPDATE categories
		SET name = $1, slug = $2, description = $3, parent_id = $4, 
		    is_active = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at
	`

	return r.db.QueryRow(
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.IsActive,
		category.ID,
	).Scan(&category.UpdatedAt)
}

// ADDED: UpdateStatus - update category active status (dipanggil di category_service.go)
func (r *CategoryRepository) UpdateStatus(id int, isActive bool) error {
	query := `
		UPDATE categories
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
		return fmt.Errorf("category not found")
	}

	return nil
}

// Delete deletes a category (hard delete)
func (r *CategoryRepository) Delete(id int) error {
	// First check if category has subcategories
	var subCount int
	checkQuery := `SELECT COUNT(*) FROM categories WHERE parent_id = $1`
	if err := r.db.QueryRow(checkQuery, id).Scan(&subCount); err != nil {
		return err
	}

	if subCount > 0 {
		return fmt.Errorf("cannot delete category with subcategories")
	}

	// Check if category has products
	var prodCount int
	checkProdQuery := `SELECT COUNT(*) FROM product_categories WHERE category_id = $1`
	if err := r.db.QueryRow(checkProdQuery, id).Scan(&prodCount); err != nil {
		return err
	}

	if prodCount > 0 {
		return fmt.Errorf("cannot delete category with products")
	}

	// Delete category
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// SoftDelete soft deletes a category (set is_active to false)
func (r *CategoryRepository) SoftDelete(id int) error {
	query := `
		UPDATE categories
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// SlugExists checks if a slug already exists
func (r *CategoryRepository) SlugExists(slug string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE slug = $1)`
	var exists bool
	err := r.db.QueryRow(query, slug).Scan(&exists)
	return exists, err
}

// SlugExistsExcludingCategory checks if slug exists for another category
func (r *CategoryRepository) SlugExistsExcludingCategory(slug string, categoryID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE slug = $1 AND id != $2)`
	var exists bool
	err := r.db.QueryRow(query, slug, categoryID).Scan(&exists)
	return exists, err
}

// HasSubCategories checks if a category has subcategories
func (r *CategoryRepository) HasSubCategories(categoryID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE parent_id = $1)`
	var exists bool
	err := r.db.QueryRow(query, categoryID).Scan(&exists)
	return exists, err
}

// HasProducts checks if a category has products
func (r *CategoryRepository) HasProducts(categoryID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM product_categories WHERE category_id = $1)`
	var exists bool
	err := r.db.QueryRow(query, categoryID).Scan(&exists)
	return exists, err
}

// CountTotal counts total categories
func (r *CategoryRepository) CountTotal() (int64, error) {
	query := `SELECT COUNT(*) FROM categories WHERE is_active = true`
	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// ADDED: CountByStatus - count categories by active status (dipanggil di category_service.go)
func (r *CategoryRepository) CountByStatus(isActive bool) (int64, error) {
	query := `SELECT COUNT(*) FROM categories WHERE is_active = $1`
	var count int64
	err := r.db.QueryRow(query, isActive).Scan(&count)
	return count, err
}

// CountSubCategories counts subcategories of a parent
func (r *CategoryRepository) CountSubCategories(parentID int) (int64, error) {
	query := `SELECT COUNT(*) FROM categories WHERE parent_id = $1 AND is_active = true`
	var count int64
	err := r.db.QueryRow(query, parentID).Scan(&count)
	return count, err
}

// Search searches categories by name
func (r *CategoryRepository) Search(keyword string) ([]models.Category, error) {
	searchPattern := "%" + keyword + "%"

	query := `
		SELECT id, name, slug, description, parent_id, is_active, created_at, updated_at
		FROM categories
		WHERE name ILIKE $1 AND is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query, searchPattern)
	if err != nil {
		return nil, err
	}

	defer func() {
    if err := rows.Close(); err != nil {
		r.logger.Warn("rows close error", zap.Error(err))
    }
}()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&c.Description,
			&c.ParentID,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

// ADDED: GetProducts - get products for a category (dipanggil di category_service.go)
func (r *CategoryRepository) GetProducts(categoryID, limit, offset int) ([]models.Product, int64, error) {
	// Get total count
	var total int64
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		INNER JOIN product_categories pc ON p.id = pc.product_id
		WHERE pc.category_id = $1 AND p.is_active = true
	`
	if err := r.db.QueryRow(countQuery, categoryID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get products
	query := `
		SELECT DISTINCT p.id, p.name, p.slug, p.description, p.price, p.stock, 
		       p.sku, p.image_url, p.is_active, p.created_at, p.updated_at
		FROM products p
		INNER JOIN product_categories pc ON p.id = pc.product_id
		WHERE pc.category_id = $1 AND p.is_active = true
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, categoryID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		r.logger.Warn("rows close error", zap.Error(err))
    }
}()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Slug,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.SKU,
			&p.ImageURL,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}

	return products, total, nil
}
