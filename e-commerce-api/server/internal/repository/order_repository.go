// order dan order item itu satu code saja soalnya fungsinya saling berkaitan

package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/querybuilder"
)

type OrderRepository struct {
	db *sql.DB
	logger *zap.Logger
	
}

// Newquerybuilder it is save backend
var orderOrderBuilder = querybuilder.NewSQLBuilder().
	AllowColumns("created_at", "updated_at", "total_amount", "status", "order_number").
	SetDefault("created_at", "DESC")

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sql.DB, logger *zap.Logger) *OrderRepository {
	return &OrderRepository{
		db: db,
		logger: logger,
	}
}

// Create creates a new order
func (r *OrderRepository) Create(order *models.Order) error {
	query := `
		INSERT INTO orders (user_id, order_number, status, total_amount, 
		                    payment_method, shipping_address, shipping_phone, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		order.UserID,
		order.OrderNumber,
		order.Status,
		order.TotalAmount,
		order.PaymentMethod, // update
		order.ShippingAddress,
		order.ShippingPhone,
		order.Notes,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	return err
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(id int) (*models.Order, error) {
	query := `
		SELECT id, user_id, order_number, status, total_amount, 
		       payment_method, shipping_address, shipping_phone, notes, 
		       created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	order := &models.Order{}
	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.TotalAmount,
		&order.PaymentMethod, // TAMBAHKAN INI
		&order.ShippingAddress,
		&order.ShippingPhone,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}

	return order, err
}

// GetByOrderNumber retrieves an order by order number
func (r *OrderRepository) GetByOrderNumber(orderNumber string) (*models.Order, error) {
	query := `
		SELECT id, user_id, order_number, status, total_amount, 
		       payment_method, shipping_address, shipping_phone, notes, 
		       created_at, updated_at
		FROM orders
		WHERE order_number = $1
	`

	order := &models.Order{}
	err := r.db.QueryRow(query, orderNumber).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.TotalAmount,
		&order.PaymentMethod, // update
		&order.ShippingAddress,
		&order.ShippingPhone,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}

	return order, err
}

// GetByUserID retrieves all orders for a specific user
func (r *OrderRepository) GetByUserID(userID int, limit, offset int) ([]models.Order, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	if err := r.db.QueryRow(countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get orders
	query := `
		SELECT id, user_id, order_number, status, total_amount, 
		       payment_method, shipping_address, shipping_phone, notes, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.OrderNumber,
			&o.Status,
			&o.TotalAmount,
			&o.PaymentMethod, // update
			&o.ShippingAddress,
			&o.ShippingPhone,
			&o.Notes,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

// GetAll retrieves all orders with pagination (for admin)
func (r *OrderRepository) GetAll(limit, offset int) ([]models.Order, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM orders`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get orders
	query := `
		SELECT id, user_id, order_number, status, total_amount, 
		       payment_method, shipping_address, shipping_phone, notes, created_at, updated_at
		FROM orders
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

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.OrderNumber,
			&o.Status,
			&o.TotalAmount,
			&o.PaymentMethod, // update
			&o.ShippingAddress,
			&o.ShippingPhone,
			&o.Notes,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

// GetByStatus retrieves orders by status
func (r *OrderRepository) GetByStatus(status string, limit, offset int) ([]models.Order, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM orders WHERE status = $1`
	if err := r.db.QueryRow(countQuery, status).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get orders
	query := `
		SELECT id, user_id, order_number, status, total_amount, 
		       payment_method, shipping_address, shipping_phone, notes, created_at, updated_at
		FROM orders
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.OrderNumber,
			&o.Status,
			&o.TotalAmount,
			&o.PaymentMethod, // update
			&o.ShippingAddress,
			&o.ShippingPhone,
			&o.Notes,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

// GetByUserIDAndStatus retrieves orders by user ID and status
func (r *OrderRepository) GetByUserIDAndStatus(userID int, status string, limit, offset int) ([]models.Order, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1 AND status = $2`
	if err := r.db.QueryRow(countQuery, userID, status).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get orders
	query := `
		SELECT id, user_id, order_number, status, total_amount, 
		       payment_method, shipping_address, shipping_phone, notes, created_at, updated_at
		FROM orders
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, userID, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.OrderNumber,
			&o.Status,
			&o.TotalAmount,
			&o.PaymentMethod, // update
			&o.ShippingAddress,
			&o.ShippingPhone,
			&o.Notes,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

// GetByDateRange retrieves orders within a date range
func (r *OrderRepository) GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]models.Order, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM orders WHERE created_at BETWEEN $1 AND $2`
	if err := r.db.QueryRow(countQuery, startDate, endDate).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get orders
	query := `
		SELECT id, user_id, order_number, status, total_amount, 
		       payment_method, shipping_address, shipping_phone, notes, created_at, updated_at
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.OrderNumber,
			&o.Status,
			&o.TotalAmount,
			&o.PaymentMethod, // update
			&o.ShippingAddress,
			&o.ShippingPhone,
			&o.Notes,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

// Update updates an order
func (r *OrderRepository) Update(order *models.Order) error {
	query := `
		UPDATE orders
		SET shipping_address = $1, shipping_phone = $2, notes = $3, 
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	return r.db.QueryRow(
		query,
		order.ShippingAddress,
		order.ShippingPhone,
		order.Notes,
		order.ID,
	).Scan(&order.UpdatedAt)
}

// UpdateStatus updates order status
func (r *OrderRepository) UpdateStatus(id int, status string) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// UpdateTotalAmount updates order total amount
func (r *OrderRepository) UpdateTotalAmount(id int, totalAmount float64) error {
	query := `
		UPDATE orders
		SET total_amount = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, totalAmount, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// Delete deletes an order (hard delete)
func (r *OrderRepository) Delete(id int) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// Cancel cancels an order (set status to cancelled)
func (r *OrderRepository) Cancel(id int) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND status IN ($3, $4)
	`

	result, err := r.db.Exec(query, models.OrderStatusCancelled, id, models.OrderStatusPending, models.OrderStatusProcessing)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("order cannot be cancelled (already shipped/delivered or not found)")
	}

	return nil
}

// OrderNumberExists checks if an order number already exists
func (r *OrderRepository) OrderNumberExists(orderNumber string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE order_number = $1)`
	var exists bool
	err := r.db.QueryRow(query, orderNumber).Scan(&exists)
	return exists, err
}

// CountByStatus counts orders by status
func (r *OrderRepository) CountByStatus(status string) (int64, error) {
	query := `SELECT COUNT(*) FROM orders WHERE status = $1`
	var count int64
	err := r.db.QueryRow(query, status).Scan(&count)
	return count, err
}

// CountByUserID counts orders for a specific user
func (r *OrderRepository) CountByUserID(userID int) (int64, error) {
	query := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	var count int64
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// CountTotalOrders counts total orders
func (r *OrderRepository) CountTotalOrders() (int64, error) {
	query := `SELECT COUNT(*) FROM orders`
	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// GetTotalSales calculates total sales amount
func (r *OrderRepository) GetTotalSales() (float64, error) {
	query := `SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status != $1`
	var total float64
	err := r.db.QueryRow(query, models.OrderStatusCancelled).Scan(&total)
	return total, err
}

// GetTotalSalesByDateRange calculates total sales within a date range
func (r *OrderRepository) GetTotalSalesByDateRange(startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM orders 
		WHERE status != $1 AND created_at BETWEEN $2 AND $3
	`
	var total float64
	err := r.db.QueryRow(query, models.OrderStatusCancelled, startDate, endDate).Scan(&total)
	return total, err
}

// GetTotalSalesByStatus calculates total sales by status
func (r *OrderRepository) GetTotalSalesByStatus(status string) (float64, error) {
	query := `SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status = $1`
	var total float64
	err := r.db.QueryRow(query, status).Scan(&total)
	return total, err
}

// GetRecentOrders retrieves recent orders (for admin dashboard)
func (r *OrderRepository) GetRecentOrders(limit int) ([]models.Order, error) {
	query := `
		SELECT id, user_id, order_number, status, total_amount, 
		       payment_method, shipping_address, shipping_phone, notes, created_at, updated_at
		FROM orders
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

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.OrderNumber,
			&o.Status,
			&o.TotalAmount,
			&o.PaymentMethod, // update
			&o.ShippingAddress,
			&o.ShippingPhone,
			&o.Notes,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

// GetPendingOrders retrieves pending orders (need action)
func (r *OrderRepository) GetPendingOrders(limit, offset int) ([]models.Order, int64, error) {
	return r.GetByStatus(models.OrderStatusPending, limit, offset)
}

// SearchByOrderNumber searches orders by order number (partial match)
func (r *OrderRepository) SearchByOrderNumber(searchTerm string, limit, offset int) ([]models.Order, int64, error) {
	searchPattern := "%" + searchTerm + "%"

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM orders WHERE order_number ILIKE $1`
	if err := r.db.QueryRow(countQuery, searchPattern).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get orders
	query := `
		SELECT id, user_id, order_number, status, total_amount, payment_method,
		       shipping_address, shipping_phone, notes, created_at, updated_at
		FROM orders
		WHERE order_number ILIKE $1
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

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.OrderNumber,
			&o.Status,
			&o.TotalAmount,
			&o.PaymentMethod,
			&o.ShippingAddress,
			&o.ShippingPhone,
			&o.Notes,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}


// ORDER ITEMS METHODS


// AddOrderItem adds an item to an order
func (r *OrderRepository) AddOrderItem(orderItem *models.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, price, subtotal)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		orderItem.OrderID,
		orderItem.ProductID,
		orderItem.Quantity,
		orderItem.Price,
		orderItem.Subtotal,
	).Scan(&orderItem.ID, &orderItem.CreatedAt)

	return err
}

// GetOrderItems retrieves all items for a specific order
func (r *OrderRepository) GetOrderItems(orderID int) ([]models.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity, price, subtotal, created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	items := []models.OrderItem{}
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}


// STATISTICS & ANALYTICS METHODS

// CountTotal counts total orders (alias untuk CountTotalOrders)
func (r *OrderRepository) CountTotal() (int64, error) {
	return r.CountTotalOrders()
}

// GetTotalRevenue calculates total revenue from all completed orders
func (r *OrderRepository) GetTotalRevenue() (float64, error) {
	query := `
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM orders 
		WHERE status IN ($1, $2)
	`
	var total float64
	err := r.db.QueryRow(query, models.OrderStatusDelivered, models.OrderStatusShipped).Scan(&total)
	return total, err
}

// CountTodayOrders counts orders created today
func (r *OrderRepository) CountTodayOrders() (int64, error) {
	query := `
		SELECT COUNT(*) 
		FROM orders 
		WHERE DATE(created_at) = CURRENT_DATE
	`
	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// GetTodayRevenue calculates revenue from orders created today
func (r *OrderRepository) GetTodayRevenue() (float64, error) {
	query := `
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM orders 
		WHERE DATE(created_at) = CURRENT_DATE 
		  AND status IN ($1, $2)
	`
	var total float64
	err := r.db.QueryRow(query, models.OrderStatusDelivered, models.OrderStatusShipped).Scan(&total)
	return total, err
}

// GetAllWithFilters retrieves orders with dynamic filters
// UPDATED: Now uses query builder to prevent SQL injection
func (r *OrderRepository) GetAllWithFilters(filters map[string]interface{}, limit, offset int) ([]models.Order, int64, error) {
	// NEW: Use WHERE clause builder for safe parameterized queries
	where := querybuilder.NewWhereClause()

	// Status filter
	if status, ok := filters["status"].(string); ok && status != "" {
		argNum := where.ArgCount()
		where.AddCondition(fmt.Sprintf("status = $%d", argNum), status)
	}

	// User ID filter
	if userID, ok := filters["user_id"].(int); ok && userID > 0 {
		argNum := where.ArgCount()
		where.AddCondition(fmt.Sprintf("user_id = $%d", argNum), userID)
	}

	// Date range filters
	if startDate, ok := filters["start_date"].(time.Time); ok && !startDate.IsZero() {
		argNum := where.ArgCount()
		where.AddCondition(fmt.Sprintf("created_at >= $%d", argNum), startDate)
	}

	if endDate, ok := filters["end_date"].(time.Time); ok && !endDate.IsZero() {
		argNum := where.ArgCount()
		where.AddCondition(fmt.Sprintf("created_at <= $%d", argNum), endDate)
	}

	// Build WHERE clause
	whereClause, args := where.Build()
	if whereClause == "" {
		whereClause = "WHERE 1=1" // Ensure valid SQL if no filters
	}

	// Get total count
	var total int64
	countQuery := "SELECT COUNT(*) FROM orders " + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	//  NEW: Use query builder for ORDER BY (prevents SQL injection)
	sortBy := ""
	sortOrder := ""

	if sort, ok := filters["sort_by"].(string); ok {
		sortBy = sort
	}
	if order, ok := filters["sort_order"].(string); ok {
		sortOrder = order
	}

	orderClause := orderOrderBuilder.BuildOrderClause(sortBy, sortOrder)

	// NEW: Safe pagination
	args = append(args, limit, offset)
	argCountForLimit := len(args) - 1

	// Build final query
	query := fmt.Sprintf(`
		SELECT id, user_id, order_number, status, total_amount, 
		       payment_method, shipping_address, shipping_phone, notes, 
		       created_at, updated_at
		FROM orders
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

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.OrderNumber,
			&o.Status,
			&o.TotalAmount,
			&o.PaymentMethod,
			&o.ShippingAddress,
			&o.ShippingPhone,
			&o.Notes,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}
