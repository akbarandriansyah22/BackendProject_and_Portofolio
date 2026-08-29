// product dan product item itu di gabung jadi 1 di product_repository.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go.uber.org/zap"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/querybuilder"
)

type ProductRepository struct {
	db *sql.DB
	logger *zap.Logger
}

// Create reusable order builder for products (update)
var productOrderBuilder = querybuilder.NewSQLBuilder().
	AllowColumns("created_at", "updated_at", "price", "name", "stock").
	SetDefault("created_at", "DESC")

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB, logger *zap.Logger) *ProductRepository {
	return &ProductRepository{
		db: db,
	logger: logger,
	}
}

// Create creates a new product
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (name, slug, description, price, stock, sku, image_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.Stock,
		product.SKU,
		product.ImageURL,
		product.IsActive,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return err
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.SKU,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}

	return product, err
}

// GetBySlug retrieves a product by slug
func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		WHERE slug = $1
	`

	product := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.SKU,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}

	return product, err
}

// GetAll retrieves all products with pagination (including inactive)
func (r *ProductRepository) GetAll(limit, offset int) ([]models.Product, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM products`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get products
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
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

// GetActive retrieves only active products with pagination
// UPDATE METHOD!
func (r *ProductRepository) GetActive(limit, offset int) ([]models.Product, int64, error) {
	// Get total count (only active)
	var total int64
	countQuery := `SELECT COUNT(*) FROM products WHERE is_active = true`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get products (only active)
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		WHERE is_active = true
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

// List retrieves products with filter and pagination
func (r *ProductRepository) List(ctx context.Context, filter *models.ProductFilter) ([]*models.Product, int, error) {
	offset := (filter.Page - 1) * filter.Limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM products WHERE is_active = true`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get products
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, filter.Limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	products := []*models.Product{}
	for rows.Next() {
		p := &models.Product{}
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

// GetByCategory retrieves products by category ID
func (r *ProductRepository) GetByCategory(categoryID int, limit, offset int) ([]models.Product, int64, error) {
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
		log.Printf("failed to close rows: %v", err)
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

// Update updates a product
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, slug = $2, description = $3, price = $4, stock = $5, 
		    sku = $6, image_url = $7, is_active = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.Stock,
		product.SKU,
		product.ImageURL,
		product.IsActive,
		product.ID,
	).Scan(&product.UpdatedAt)
}

// Delete deletes a product (soft delete - set is_active to false)
func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE products
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// HardDelete deletes a product permanently (hard delete)
func (r *ProductRepository) HardDelete(id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// UpdateStock updates product stock (direct set)
func (r *ProductRepository) UpdateStock(ctx context.Context, productID, quantity int) error {
	query := `
		UPDATE products
		SET stock = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, quantity, productID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// DecrementStock decrements product stock (for checkout)
func (r *ProductRepository) DecrementStock(ctx context.Context, productID, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock - $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND stock >= $1
	`

	result, err := r.db.ExecContext(ctx, query, quantity, productID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("insufficient stock or product not found")
	}

	return nil
}

// IncrementStock increments product stock (for cancel/return)
func (r *ProductRepository) IncrementStock(ctx context.Context, productID, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, quantity, productID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// Search searches products by name or description
func (r *ProductRepository) Search(ctx context.Context, keyword string, page, limit int) ([]*models.Product, int, error) {
	searchPattern := "%" + keyword + "%"
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE (name ILIKE $1 OR description ILIKE $1) AND is_active = true
	`
	if err := r.db.QueryRowContext(ctx, countQuery, searchPattern).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get products
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		WHERE (name ILIKE $1 OR description ILIKE $1) AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	products := []*models.Product{}
	for rows.Next() {
		p := &models.Product{}
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

// GetByPriceRange retrieves products within a price range
func (r *ProductRepository) GetByPriceRange(minPrice, maxPrice float64, limit, offset int) ([]models.Product, int64, error) {
	// Get total count
	var total int64
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE price BETWEEN $1 AND $2 AND is_active = true
	`
	if err := r.db.QueryRow(countQuery, minPrice, maxPrice).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get products
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		WHERE price BETWEEN $1 AND $2 AND is_active = true
		ORDER BY price ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, minPrice, maxPrice, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
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

// GetLowStock retrieves products with stock below threshold
func (r *ProductRepository) GetLowStock(threshold int, limit int) ([]models.Product, error) {
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		WHERE stock <= $1 AND is_active = true
		ORDER BY stock ASC
		LIMIT $2
	`

	rows, err := r.db.Query(query, threshold, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
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
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// GetOutOfStock retrieves out of stock products
func (r *ProductRepository) GetOutOfStock(limit, offset int) ([]models.Product, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM products WHERE stock = 0 AND is_active = true`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get products
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		WHERE stock = 0 AND is_active = true
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

// SlugExists checks if a slug already exists
func (r *ProductRepository) SlugExists(slug string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE slug = $1)`
	var exists bool
	err := r.db.QueryRow(query, slug).Scan(&exists)
	return exists, err
}

// SlugExistsExcludingProduct checks if slug exists for another product
func (r *ProductRepository) SlugExistsExcludingProduct(slug string, productID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE slug = $1 AND id != $2)`
	var exists bool
	err := r.db.QueryRow(query, slug, productID).Scan(&exists)
	return exists, err
}

// SKUExists checks if a SKU already exists
func (r *ProductRepository) SKUExists(sku string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1)`
	var exists bool
	err := r.db.QueryRow(query, sku).Scan(&exists)
	return exists, err
}

// SKUExistsExcludingProduct checks if SKU exists for another product
func (r *ProductRepository) SKUExistsExcludingProduct(sku string, productID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1 AND id != $2)`
	var exists bool
	err := r.db.QueryRow(query, sku, productID).Scan(&exists)
	return exists, err
}

// CountTotal counts total products
func (r *ProductRepository) CountTotal() (int64, error) {
	query := `SELECT COUNT(*) FROM products WHERE is_active = true`
	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// CountByCategory counts products in a category
func (r *ProductRepository) CountByCategory(categoryID int) (int64, error) {
	query := `
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		INNER JOIN product_categories pc ON p.id = pc.product_id
		WHERE pc.category_id = $1 AND p.is_active = true
	`
	var count int64
	err := r.db.QueryRow(query, categoryID).Scan(&count)
	return count, err
}

// GetAllWithFilters retrieves products with dynamic filters
// UPDATED: Now uses query builder to prevent SQL injection
func (r *ProductRepository) GetAllWithFilters(filters map[string]interface{}, limit, offset int) ([]models.Product, int64, error) {
	// Use WHERE clause builder for safe parameterized queries
	where := querybuilder.NewWhereClause()

	// Search filter (safe - uses parameterized query)
	if search, ok := filters["search"].(string); ok && search != "" {
		argNum := where.ArgCount()
		where.AddCondition(
			fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argNum, argNum),
			"%"+search+"%",
		)
	}

	// Category filter (safe - uses parameterized query)
	if categoryID, ok := filters["category_id"].(int); ok && categoryID > 0 {
		argNum := where.ArgCount()
		where.AddCondition(
			fmt.Sprintf("id IN (SELECT product_id FROM product_categories WHERE category_id = $%d)", argNum),
			categoryID,
		)
	}

	// Price range filters (safe - uses parameterized query)
	if minPrice, ok := filters["min_price"].(float64); ok && minPrice > 0 {
		argNum := where.ArgCount()
		where.AddCondition(fmt.Sprintf("price >= $%d", argNum), minPrice)
	}

	if maxPrice, ok := filters["max_price"].(float64); ok && maxPrice > 0 {
		argNum := where.ArgCount()
		where.AddCondition(fmt.Sprintf("price <= $%d", argNum), maxPrice)
	}

	// Only active products (default)
	where.AddCondition("is_active = true")

	//  Build WHERE clause safely
	whereClause, args := where.Build()

	// Get total count
	var total int64
	countQuery := "SELECT COUNT(*) FROM products " + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	//  Use query builder for ORDER BY (prevents SQL injection)
	sortBy := ""
	sortOrder := ""

	if sort, ok := filters["sort_by"].(string); ok {
		sortBy = sort
	}
	if order, ok := filters["sort_order"].(string); ok {
		sortOrder = order
	}

	orderClause := productOrderBuilder.BuildOrderClause(sortBy, sortOrder)

	//  Safe pagination
	args = append(args, limit, offset)
	argCountForLimit := len(args) - 1

	// Build final query
	query := fmt.Sprintf(`
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderClause, argCountForLimit, argCountForLimit+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
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

// GetFeatured retrieves featured products (most recent active products)
func (r *ProductRepository) GetFeatured(limit int) ([]models.Product, error) {
	query := `
		SELECT id, name, slug, description, price, stock, sku, image_url, 
		       is_active, created_at, updated_at
		FROM products
		WHERE is_active = true AND stock > 0
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
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
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// UpdateStatus updates product active status
func (r *ProductRepository) UpdateStatus(id int, isActive bool) error {
	query := `
		UPDATE products
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
		return fmt.Errorf("product not found")
	}

	return nil
}

// CountByStatus counts products by active status
func (r *ProductRepository) CountByStatus(isActive bool) (int64, error) {
	query := `SELECT COUNT(*) FROM products WHERE is_active = $1`
	var count int64
	err := r.db.QueryRow(query, isActive).Scan(&count)
	return count, err
}

// CountLowStock counts products with stock below or equal to threshold
func (r *ProductRepository) CountLowStock(threshold int) (int64, error) {
	query := `SELECT COUNT(*) FROM products WHERE stock <= $1 AND stock > 0 AND is_active = true`
	var count int64
	err := r.db.QueryRow(query, threshold).Scan(&count)
	return count, err
}

// CountOutOfStock counts products with zero stock
func (r *ProductRepository) CountOutOfStock() (int64, error) {
	query := `SELECT COUNT(*) FROM products WHERE stock = 0 AND is_active = true`
	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// GetTotalStockValue calculates total value of all stock (price * stock)
func (r *ProductRepository) GetTotalStockValue() (float64, error) {
	query := `SELECT COALESCE(SUM(price * stock), 0) FROM products WHERE is_active = true`
	var total float64
	err := r.db.QueryRow(query).Scan(&total)
	return total, err
}


// CATEGORY RELATIONSHIP METHODS


// AddToCategory adds a product to a category
func (r *ProductRepository) AddToCategory(productID, categoryID int) error {
	query := `
		INSERT INTO product_categories (product_id, category_id, created_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (product_id, category_id) DO NOTHING
	`

	_, err := r.db.Exec(query, productID, categoryID)
	return err
}

// RemoveFromCategory removes a product from a specific category
func (r *ProductRepository) RemoveFromCategory(productID, categoryID int) error {
	query := `DELETE FROM product_categories WHERE product_id = $1 AND category_id = $2`
	_, err := r.db.Exec(query, productID, categoryID)
	return err
}

// RemoveFromAllCategories removes a product from all categories
func (r *ProductRepository) RemoveFromAllCategories(productID int) error {
	query := `DELETE FROM product_categories WHERE product_id = $1`
	_, err := r.db.Exec(query, productID)
	return err
}

// GetCategories retrieves all categories for a product
func (r *ProductRepository) GetCategories(productID int) ([]models.Category, error) {
	query := `
		SELECT c.id, c.name, c.slug, c.description, c.parent_id, 
		       c.is_active, c.created_at, c.updated_at
		FROM categories c
		INNER JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = $1
		ORDER BY c.name
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
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

// GetProductsByMultipleCategories retrieves products that belong to any of the given categories
func (r *ProductRepository) GetProductsByMultipleCategories(categoryIDs []int, limit, offset int) ([]models.Product, int64, error) {
	if len(categoryIDs) == 0 {
		return []models.Product{}, 0, nil
	}

	// Build placeholder string for IN clause
	placeholders := make([]string, len(categoryIDs))
	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	inClause := "(" + placeholders[0]
	for i := 1; i < len(placeholders); i++ {
		inClause += ", " + placeholders[i]
	}
	inClause += ")"

	// Get total count
	var total int64
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		INNER JOIN product_categories pc ON p.id = pc.product_id
		WHERE pc.category_id IN %s AND p.is_active = true
	`, inClause)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get products
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT DISTINCT p.id, p.name, p.slug, p.description, p.price, p.stock, 
		       p.sku, p.image_url, p.is_active, p.created_at, p.updated_at
		FROM products p
		INNER JOIN product_categories pc ON p.id = pc.product_id
		WHERE pc.category_id IN %s AND p.is_active = true
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, inClause, len(categoryIDs)+1, len(categoryIDs)+2)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
    defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
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

// GetProductsByCategory retrieves products by category ID with context
func (r *ProductRepository) GetProductsByCategory(ctx context.Context, categoryID int, page, limit int) ([]*models.Product, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		INNER JOIN product_categories pc ON p.id = pc.product_id
		WHERE pc.category_id = $1 AND p.is_active = true
	`
	if err := r.db.QueryRowContext(ctx, countQuery, categoryID).Scan(&total); err != nil {
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

	rows, err := r.db.QueryContext(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	products := []*models.Product{}
	for rows.Next() {
		p := &models.Product{}
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

// SetActive sets the active status of a product
func (r *ProductRepository) SetActive(ctx context.Context, productID int, isActive bool) error {
	query := `
		UPDATE products
		SET is_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, isActive, productID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *ProductRepository) CheckStock(ctx context.Context, productID int, qty int) (bool, error) {
	var stock int
	err := r.db.QueryRowContext(
		ctx,
		"SELECT stock FROM products WHERE id = $1",
		productID,
	).Scan(&stock)

	if err != nil {
		return false, err
	}

	if stock < qty {
		return false, nil
	}

	return true, nil
}
