// cart dan cart item itu satu code saja soalnya fungsinya saling berkaitan

package repository

import (
	"database/sql"
	"fmt"
	"log"

	"go.uber.org/zap"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

type CartRepository struct {
	db *sql.DB
	logger *zap.Logger
}


// NewCartRepository creates a new cart repository
func NewCartRepository(db *sql.DB, logger *zap.Logger) *CartRepository {
	return &CartRepository{
		db: db,
		logger: logger,
	}
}

// GetOrCreateCart gets existing cart or creates new one for user
func (r *CartRepository) GetOrCreateCart(userID int) (*models.Cart, error) {
	// Try to get existing cart
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`
	cart := &models.Cart{}
	err := r.db.QueryRow(query, userID).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)

	if err == sql.ErrNoRows {
		// Create new cart
		createQuery := `
			INSERT INTO carts (user_id)
			VALUES ($1)
			RETURNING id, user_id, created_at, updated_at
		`
		err = r.db.QueryRow(createQuery, userID).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return cart, nil
}

// GetByID retrieves a cart by ID
func (r *CartRepository) GetByID(id int) (*models.Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE id = $1`

	cart := &models.Cart{}
	err := r.db.QueryRow(query, id).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart not found")
	}

	return cart, err
}

// GetByUserID retrieves a cart by user ID
func (r *CartRepository) GetByUserID(userID int) (*models.Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`

	cart := &models.Cart{}
	err := r.db.QueryRow(query, userID).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart not found")
	}

	return cart, err
}

// AddItem adds an item to cart
func (r *CartRepository) AddItem(cartID, productID int, quantity int, price float64) error {
	// Check if item already exists in cart
	var existingID int
	var existingQty int
	checkQuery := `SELECT id, quantity FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	err := r.db.QueryRow(checkQuery, cartID, productID).Scan(&existingID, &existingQty)

	if err == sql.ErrNoRows {
		// Insert new item
		insertQuery := `
			INSERT INTO cart_items (cart_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`
		_, err = r.db.Exec(insertQuery, cartID, productID, quantity, price)
		return err
	} else if err != nil {
		return err
	}

	// Update existing item quantity
	updateQuery := `
		UPDATE cart_items
		SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err = r.db.Exec(updateQuery, quantity, existingID)
	return err
}

// GetCartItems retrieves all items in a cart
func (r *CartRepository) GetCartItems(cartID int) ([]models.CartItemWithProduct, error) {
	query := `
		SELECT 
			ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.price, 
			ci.created_at, ci.updated_at,
			p.id, p.name, p.slug, p.description, p.price, p.stock, 
			p.sku, p.image_url, p.is_active, p.created_at, p.updated_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1
	`

	rows, err := r.db.Query(query, cartID)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()
	items := []models.CartItemWithProduct{}
	for rows.Next() {
		var item models.CartItemWithProduct
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Product.ID,
			&item.Product.Name,
			&item.Product.Slug,
			&item.Product.Description,
			&item.Product.Price,
			&item.Product.Stock,
			&item.Product.SKU,
			&item.Product.ImageURL,
			&item.Product.IsActive,
			&item.Product.CreatedAt,
			&item.Product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetCartItemByID retrieves a cart item by ID
func (r *CartRepository) GetCartItemByID(id int) (*models.CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, quantity, price, created_at, updated_at
		FROM cart_items
		WHERE id = $1
	`

	item := &models.CartItem{}
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Quantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart item not found")
	}

	return item, err
}

// GetCartItemByCartAndProduct retrieves a cart item by cart ID and product ID
func (r *CartRepository) GetCartItemByCartAndProduct(cartID, productID int) (*models.CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, quantity, price, created_at, updated_at
		FROM cart_items
		WHERE cart_id = $1 AND product_id = $2
	`

	item := &models.CartItem{}
	err := r.db.QueryRow(query, cartID, productID).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Quantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart item not found")
	}

	return item, err
}

// UpdateItemQuantity updates cart item quantity
func (r *CartRepository) UpdateItemQuantity(cartItemID, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	result, err := r.db.Exec(query, quantity, cartItemID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// UpdateItemPrice updates cart item price (when product price changes)
func (r *CartRepository) UpdateItemPrice(cartItemID int, price float64) error {
	query := `
		UPDATE cart_items
		SET price = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	result, err := r.db.Exec(query, price, cartItemID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// IncrementItemQuantity increments cart item quantity
func (r *CartRepository) IncrementItemQuantity(cartItemID, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	result, err := r.db.Exec(query, quantity, cartItemID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// DecrementItemQuantity decrements cart item quantity
func (r *CartRepository) DecrementItemQuantity(cartItemID, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND quantity >= $1
	`
	result, err := r.db.Exec(query, quantity, cartItemID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("insufficient quantity or cart item not found")
	}

	return nil
}

// RemoveItem removes an item from cart
func (r *CartRepository) RemoveItem(cartItemID int) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	result, err := r.db.Exec(query, cartItemID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// ClearCart removes all items from cart
func (r *CartRepository) ClearCart(cartID int) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.Exec(query, cartID)
	return err
}

// GetCartTotal calculates total amount of cart
func (r *CartRepository) GetCartTotal(cartID int) (float64, error) {
	query := `
		SELECT COALESCE(SUM(price * quantity), 0)
		FROM cart_items
		WHERE cart_id = $1
	`
	var total float64
	err := r.db.QueryRow(query, cartID).Scan(&total)
	return total, err
}

// GetCartItemCount counts total items in cart
func (r *CartRepository) GetCartItemCount(cartID int) (int, error) {
	query := `SELECT COUNT(*) FROM cart_items WHERE cart_id = $1`
	var count int
	err := r.db.QueryRow(query, cartID).Scan(&count)
	return count, err
}

// GetCartTotalQuantity counts total quantity of all items in cart
func (r *CartRepository) GetCartTotalQuantity(cartID int) (int, error) {
	query := `SELECT COALESCE(SUM(quantity), 0) FROM cart_items WHERE cart_id = $1`
	var total int
	err := r.db.QueryRow(query, cartID).Scan(&total)
	return total, err
}

// ItemExistsInCart checks if a product already exists in cart
func (r *CartRepository) ItemExistsInCart(cartID, productID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM cart_items WHERE cart_id = $1 AND product_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, cartID, productID).Scan(&exists)
	return exists, err
}

// DeleteCart deletes a cart (hard delete)
func (r *CartRepository) DeleteCart(id int) error {
	query := `DELETE FROM carts WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("cart not found")
	}

	return nil
}

// SyncCartItemPrices syncs all cart item prices with current product prices
func (r *CartRepository) SyncCartItemPrices(cartID int) error {
	query := `
		UPDATE cart_items ci
		SET price = p.price, updated_at = CURRENT_TIMESTAMP
		FROM products p
		WHERE ci.product_id = p.id AND ci.cart_id = $1
	`
	_, err := r.db.Exec(query, cartID)
	return err
}

// RemoveInactiveProducts removes products that are no longer active from cart
func (r *CartRepository) RemoveInactiveProducts(cartID int) error {
	query := `
		DELETE FROM cart_items ci
		USING products p
		WHERE ci.product_id = p.id 
		AND ci.cart_id = $1 
		AND p.is_active = false
	`
	_, err := r.db.Exec(query, cartID)
	return err
}

// RemoveOutOfStockProducts removes products with zero stock from cart
func (r *CartRepository) RemoveOutOfStockProducts(cartID int) error {
	query := `
		DELETE FROM cart_items ci
		USING products p
		WHERE ci.product_id = p.id 
		AND ci.cart_id = $1 
		AND p.stock = 0
	`
	_, err := r.db.Exec(query, cartID)
	return err
}

// AdjustQuantityToStock adjusts cart item quantities to available stock
func (r *CartRepository) AdjustQuantityToStock(cartID int) error {
	query := `
		UPDATE cart_items ci
		SET quantity = LEAST(ci.quantity, p.stock), updated_at = CURRENT_TIMESTAMP
		FROM products p
		WHERE ci.product_id = p.id 
		AND ci.cart_id = $1 
		AND ci.quantity > p.stock
	`
	_, err := r.db.Exec(query, cartID)
	return err
}
