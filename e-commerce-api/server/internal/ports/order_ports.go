package ports

import (
	"context"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

// OrderRepository mendefinisikan kontrak untuk data access layer Order
type OrderRepository interface {
	// Create membuat order baru
	Create(ctx context.Context, order *models.Order) error

	// CreateOrderItems membuat order items
	CreateOrderItems(ctx context.Context, orderID int, items []*models.OrderItem) error

	// GetByID mengambil order berdasarkan ID
	GetByID(ctx context.Context, id int) (*models.Order, error)

	// GetByOrderNumber mengambil order berdasarkan order number
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error)

	// GetOrderItems mengambil semua item dalam order
	GetOrderItems(ctx context.Context, orderID int) ([]*models.OrderItem, error)

	// Update memperbarui data order
	Update(ctx context.Context, order *models.Order) error

	// UpdateStatus memperbarui status order
	UpdateStatus(ctx context.Context, orderID int, status string) error

	// GetUserOrders mengambil semua order dari user
	GetUserOrders(ctx context.Context, userID int, page, limit int) ([]*models.Order, int, error)

	// GetAllOrders mengambil semua order (admin only)
	GetAllOrders(ctx context.Context, filter *models.OrderFilter) ([]*models.Order, int, error)

	// GenerateOrderNumber menghasilkan unique order number
	GenerateOrderNumber(ctx context.Context) (string, error)

	// CalculateTotal menghitung total harga order
	CalculateTotal(ctx context.Context, items []*models.OrderItem) float64

	// Delete menghapus order (soft delete)
	Delete(ctx context.Context, id int) error
}

// OrderService mendefinisikan kontrak untuk business logic layer Order
type OrderService interface {
	// CreateOrder membuat order baru dari cart
	CreateOrder(ctx context.Context, userID int, req *models.CreateOrderRequest) (*models.OrderResponse, error)

	// GetOrderByID mengambil detail order
	GetOrderByID(ctx context.Context, userID, orderID int) (*models.OrderDetailResponse, error)

	// GetOrderByNumber mengambil order berdasarkan order number
	GetOrderByNumber(ctx context.Context, userID int, orderNumber string) (*models.OrderDetailResponse, error)

	// GetUserOrders mengambil semua order user
	GetUserOrders(ctx context.Context, userID int, page, limit int) (*models.PaginatedResponse, error)

	// CancelOrder membatalkan order
	CancelOrder(ctx context.Context, userID, orderID int) error

	// UpdateOrderStatus memperbarui status order (admin only)
	UpdateOrderStatus(ctx context.Context, orderID int, status string) error

	// GetAllOrders mengambil semua order dengan filter (admin only)
	GetAllOrders(ctx context.Context, filter *models.OrderFilter) (*models.PaginatedResponse, error)

	// ConfirmPayment mengkonfirmasi pembayaran order
	ConfirmPayment(ctx context.Context, orderID int, paymentProof string) error

	// ProcessOrder memproses order (admin only)
	ProcessOrder(ctx context.Context, orderID int) error

	// ShipOrder mengirim order (admin only)
	ShipOrder(ctx context.Context, orderID int, trackingNumber string) error

	// CompleteOrder menyelesaikan order
	CompleteOrder(ctx context.Context, orderID int) error
}
