package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

type PaymentRepository struct {
	db *sql.DB
	logger *zap.Logger
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *sql.DB, logger *zap.Logger) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	logger: logger,
}
}

// Create creates a new payment
func (r *PaymentRepository) Create(payment *models.Payment) error {
	query := `
		INSERT INTO payments (order_id, payment_method, amount, status, transaction_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		payment.OrderID,
		payment.PaymentMethod,
		payment.Amount,
		payment.Status,
		payment.TransactionID,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	return err
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(id int) (*models.Payment, error) {
	query := `
		SELECT id, order_id, payment_method, amount, status, transaction_id, 
		       paid_at, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	payment := &models.Payment{}
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentMethod,
		&payment.Amount,
		&payment.Status,
		&payment.TransactionID,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payments not found")
	}

	return payment, err
}

// GetByOrderID retrieves a payment by order ID
func (r *PaymentRepository) GetByOrderID(orderID int) (*models.Payment, error) {
	query := `
		SELECT id, order_id, payment_method, amount, status, transaction_id, 
		       paid_at, created_at, updated_at
		FROM payments
		WHERE order_id = $1
	`

	payment := &models.Payment{}
	err := r.db.QueryRow(query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentMethod,
		&payment.Amount,
		&payment.Status,
		&payment.TransactionID,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}

	return payment, err
}

// GetByTransactionID retrieves a payment by transaction ID
func (r *PaymentRepository) GetByTransactionID(transactionID string) (*models.Payment, error) {
	query := `
		SELECT id, order_id, payment_method, amount, status, transaction_id, 
		       paid_at, created_at, updated_at
		FROM payments
		WHERE transaction_id = $1
	`

	payment := &models.Payment{}
	err := r.db.QueryRow(query, transactionID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentMethod,
		&payment.Amount,
		&payment.Status,
		&payment.TransactionID,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}

	return payment, err
}

// GetAll retrieves all payment with pagination
func (r *PaymentRepository) GetAll(limit, offset int) ([]models.Payment, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM payments`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get payment
	query := `
		SELECT id, order_id, payment_method, amount, status, transaction_id, 
		       paid_at, created_at, updated_at
		FROM payments
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

	payment := []models.Payment{}
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.PaymentMethod,
			&p.Amount,
			&p.Status,
			&p.TransactionID,
			&p.PaidAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		payment = append(payment, p)
	}

	return payment, total, nil
}

// GetByStatus retrieves payment by status
func (r *PaymentRepository) GetByStatus(status string, limit, offset int) ([]models.Payment, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM payments WHERE status = $1`
	if err := r.db.QueryRow(countQuery, status).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get payment
	query := `
		SELECT id, order_id, payment_method, amount, status, transaction_id, 
		       paid_at, created_at, updated_at
		FROM payments
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

	payment := []models.Payment{}
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.PaymentMethod,
			&p.Amount,
			&p.Status,
			&p.TransactionID,
			&p.PaidAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		payment = append(payment, p)
	}

	return payment, total, nil
}

// GetByPaymentMethod retrieves payments by payment method
func (r *PaymentRepository) GetByPaymentMethod(paymentMethod string, limit, offset int) ([]models.Payment, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM payments WHERE payment_method = $1`
	if err := r.db.QueryRow(countQuery, paymentMethod).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get payments
	query := `
		SELECT id, order_id, payment_method, amount, status, transaction_id, 
		       paid_at, created_at, updated_at
		FROM payments
		WHERE payment_method = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, paymentMethod, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
    if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
    }
}()

	payment := []models.Payment{}
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.PaymentMethod,
			&p.Amount,
			&p.Status,
			&p.TransactionID,
			&p.PaidAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		payment = append(payment, p)
	}

	return payment, total, nil
}

// GetByDateRange retrieves payment within a date range
func (r *PaymentRepository) GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]models.Payment, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM payments WHERE created_at BETWEEN $1 AND $2`
	if err := r.db.QueryRow(countQuery, startDate, endDate).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get payment
	query := `
		SELECT id, order_id, payment_method, amount, status, transaction_id, 
		       paid_at, created_at, updated_at
		FROM payments
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

	payment := []models.Payment{}
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.PaymentMethod,
			&p.Amount,
			&p.Status,
			&p.TransactionID,
			&p.PaidAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		payment = append(payment, p)
	}

	return payment, total, nil
}

// Update updates a payment
func (r *PaymentRepository) Update(payment *models.Payment) error {
	query := `
		UPDATE payments
		SET payment_method = $1, amount = $2, status = $3, 
		    transaction_id = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`

	return r.db.QueryRow(
		query,
		payment.PaymentMethod,
		payment.Amount,
		payment.Status,
		payment.TransactionID,
		payment.ID,
	).Scan(&payment.UpdatedAt)
}

// UpdateStatus updates payment status
func (r *PaymentRepository) UpdateStatus(id int, status string) error {
	query := `
		UPDATE payments
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
		return fmt.Errorf("payment not found")
	}

	return nil
}

// MarkAsPaid marks a payment as paid
func (r *PaymentRepository) MarkAsPaid(id int, transactionID string) error {
	query := `
		UPDATE payments
		SET status = $1, paid_at = $2, transaction_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	result, err := r.db.Exec(query, models.PaymentStatusSuccess, time.Now(), transactionID, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// MarkAsFailed marks a payment as failed
func (r *PaymentRepository) MarkAsFailed(id int) error {
	query := `
		UPDATE payments
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, models.PaymentStatusFailed, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// MarkAsRefunded marks a payment as refunded
func (r *PaymentRepository) MarkAsRefunded(id int) error {
	query := `
		UPDATE payments
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, models.PaymentStatusRefunded, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// Delete deletes a payment (hard delete)
func (r *PaymentRepository) Delete(id int) error {
	query := `DELETE FROM payments WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// TransactionIDExists checks if a transaction ID already exists
func (r *PaymentRepository) TransactionIDExists(transactionID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM payments WHERE transaction_id = $1)`
	var exists bool
	err := r.db.QueryRow(query, transactionID).Scan(&exists)
	return exists, err
}

// CountByStatus counts payment by status
func (r *PaymentRepository) CountByStatus(status string) (int64, error) {
	query := `SELECT COUNT(*) FROM payments WHERE status = $1`
	var count int64
	err := r.db.QueryRow(query, status).Scan(&count)
	return count, err
}

// CountByPaymentMethod counts payment by payment method
func (r *PaymentRepository) CountByPaymentMethod(paymentMethod string) (int64, error) {
	query := `SELECT COUNT(*) FROM payments WHERE payment_method = $1`
	var count int64
	err := r.db.QueryRow(query, paymentMethod).Scan(&count)
	return count, err
}

// GetTotalRevenue calculates total revenue from successful payment
func (r *PaymentRepository) GetTotalRevenue() (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE status = $1`
	var total float64
	err := r.db.QueryRow(query, models.PaymentStatusSuccess).Scan(&total)
	return total, err
}

// GetRevenueByDateRange calculates revenue within a date range
func (r *PaymentRepository) GetRevenueByDateRange(startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0) 
		FROM payments
		WHERE status = $1 AND created_at BETWEEN $2 AND $3
	`
	var total float64
	err := r.db.QueryRow(query, models.PaymentStatusSuccess, startDate, endDate).Scan(&total)
	return total, err
}

// GetRevenueByPaymentMethod calculates revenue by payment method
func (r *PaymentRepository) GetRevenueByPaymentMethod(paymentMethod string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0) 
		FROM payments
		WHERE status = $1 AND payment_method = $2
	`
	var total float64
	err := r.db.QueryRow(query, models.PaymentStatusSuccess, paymentMethod).Scan(&total)
	return total, err
}
