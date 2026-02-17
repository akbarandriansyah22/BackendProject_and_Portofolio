package ports

import (
	"context"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

// PaymentRepository mendefinisikan kontrak untuk data access layer Payment
type PaymentRepository interface {
	// Create membuat payment baru
	Create(ctx context.Context, payment *models.Payment) error

	// GetByID mengambil payment berdasarkan ID
	GetByID(ctx context.Context, id int) (*models.Payment, error)

	// GetByOrderID mengambil payment berdasarkan order ID
	GetByOrderID(ctx context.Context, orderID int) (*models.Payment, error)

	// GetByTransactionID mengambil payment berdasarkan transaction ID
	GetByTransactionID(ctx context.Context, transactionID string) (*models.Payment, error)

	// Update memperbarui data payment
	Update(ctx context.Context, payment *models.Payment) error

	// UpdateStatus memperbarui status payment
	UpdateStatus(ctx context.Context, paymentID int, status string) error

	// ConfirmPayment mengkonfirmasi payment
	ConfirmPayment(ctx context.Context, paymentID int, transactionID string) error

	// List mengambil semua payment dengan filter
	List(ctx context.Context, filter *models.PaymentFilter) ([]*models.Payment, int, error)

	// Delete menghapus payment
	Delete(ctx context.Context, id int) error
}

// PaymentService mendefinisikan kontrak untuk business logic layer Payment
type PaymentService interface {
	// CreatePayment membuat payment untuk order
	CreatePayment(ctx context.Context, orderID int, req *models.CreatePaymentRequest) (*models.PaymentResponse, error)

	// GetPaymentByID mengambil detail payment
	GetPaymentByID(ctx context.Context, paymentID int) (*models.PaymentDetailResponse, error)

	// GetPaymentByOrderID mengambil payment berdasarkan order
	GetPaymentByOrderID(ctx context.Context, orderID int) (*models.PaymentDetailResponse, error)

	// ConfirmPayment mengkonfirmasi payment (admin only)
	ConfirmPayment(ctx context.Context, paymentID int, req *models.ConfirmPaymentRequest) error

	// RejectPayment menolak payment (admin only)
	RejectPayment(ctx context.Context, paymentID int, reason string) error

	// ProcessPaymentCallback memproses callback dari payment gateway
	ProcessPaymentCallback(ctx context.Context, req *models.PaymentCallbackRequest) error

	// GetPaymentStatus mengambil status payment
	GetPaymentStatus(ctx context.Context, paymentID int) (*models.PaymentStatusResponse, error)

	// ListPayments mengambil semua payment dengan filter (admin only)
	ListPayments(ctx context.Context, filter *models.PaymentFilter) (*models.PaginatedResponse, error)
}
