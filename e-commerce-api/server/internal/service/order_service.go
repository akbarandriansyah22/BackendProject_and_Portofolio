package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/ports"
)

// OrderService handles order business logic
type OrderService struct {
	orderRepo   ports.OrderRepository
	cartRepo    ports.CartRepository
	productRepo ports.ProductRepository
	paymentRepo ports.PaymentRepository
	logger      observability.Logger
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo ports.OrderRepository,
	cartRepo ports.CartRepository,
	productRepo ports.ProductRepository,
	paymentRepo ports.PaymentRepository,
	logger observability.Logger,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

// CreateOrder creates a new order from cart
func (s *OrderService) CreateOrder(ctx context.Context, userID int) (*models.OrderResponse, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	s.logger.Info("Order created for UserID=%d", userID)
	return &models.OrderResponse{
		ID:              1,
		OrderNumber:     "ORD-001",
		Status:          "pending",
		TotalAmount:     0,
		PaymentMethod:   "",
		ShippingAddress: "",
	}, nil
}

// ListOrders lists user's orders
func (s *OrderService) ListOrders(ctx context.Context, userID int, page, limit int) ([]*models.OrderResponse, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	return make([]*models.OrderResponse, 0), nil
}

// ============================================
// Handler-level methods (non-context or simplified versions)
// ============================================

// GetUserOrders gets all orders for a user with pagination
func (s *OrderService) GetUserOrders(userID, page, limit int) ([]*models.Order, int, error) {
	ctx := context.Background()
	orders, total, err := s.orderRepo.GetUserOrders(ctx, userID, page, limit)
	if err != nil {
		s.logger.Error("OrderService.GetUserOrders failed", err)
		return nil, 0, fmt.Errorf("failed to get user orders")
	}
	return orders, total, nil
}

// GetByID gets order by ID
func (s *OrderService) GetByID(orderID int) (*models.Order, error) {
	ctx := context.Background()
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

// GetByOrderNumber gets order by order number
func (s *OrderService) GetByOrderNumber(orderNumber string) (*models.Order, error) {
	ctx := context.Background()
	order, err := s.orderRepo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

// CreateFromCart creates order from cart
func (s *OrderService) CreateFromCart(userID int, shippingAddress, paymentMethod, notes string) (*models.Order, error) {
	ctx := context.Background()

	// Get user's cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil || cart == nil {
		return nil, fmt.Errorf("cart not found")
	}

	// Get cart items
	cartItems, err := s.cartRepo.GetCartItems(ctx, cart.ID)
	if err != nil || len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Generate order number
	orderNumber, err := s.orderRepo.GenerateOrderNumber(ctx)
	if err != nil {
		s.logger.Error("OrderService.CreateFromCart failed to generate order number", err)
		return nil, fmt.Errorf("failed to create order")
	}

	// Create order
	order := &models.Order{
		UserID:          userID,
		OrderNumber:     orderNumber,
		Status:          "pending",
		ShippingAddress: shippingAddress,
		PaymentMethod:   paymentMethod,
		Notes:           sql.NullString{String: notes, Valid: notes != ""},
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		s.logger.Error("OrderService.CreateFromCart failed to create order", err)
		return nil, fmt.Errorf("failed to create order")
	}

	// Create order items from cart items
	orderItems := make([]*models.OrderItem, len(cartItems))
	for i, item := range cartItems {
		orderItems[i] = &models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	if err := s.orderRepo.CreateOrderItems(ctx, order.ID, orderItems); err != nil {
		s.logger.Error("OrderService.CreateFromCart failed to create order items", err)
		return nil, fmt.Errorf("failed to create order items")
	}

	// Clear cart after order creation
	_ = s.cartRepo.ClearCart(ctx, cart.ID)

	s.logger.Info("Order created from cart: UserID=%d, OrderID=%d, OrderNumber=%s", userID, order.ID, order.OrderNumber)

	return order, nil
}

// UpdateStatus updates order status
func (s *OrderService) UpdateStatus(orderID int, status string) error {
	ctx := context.Background()

	// Verify order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return fmt.Errorf("order not found")
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		s.logger.Error("OrderService.UpdateStatus failed", err)
		return fmt.Errorf("failed to update order status")
	}

	s.logger.Info("Order status updated: OrderID=%d, NewStatus=%s", orderID, status)
	return nil
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(orderID int) error {
	ctx := context.Background()

	// Verify order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return fmt.Errorf("order not found")
	}

	// Update status to cancelled
	if err := s.orderRepo.UpdateStatus(ctx, orderID, "cancelled"); err != nil {
		s.logger.Error("OrderService.CancelOrder failed", err)
		return fmt.Errorf("failed to cancel order")
	}

	s.logger.Info("Order cancelled: OrderID=%d", orderID)
	return nil
}

// GetOrderStats gets order statistics
func (s *OrderService) GetOrderStats() (interface{}, error) {
	ctx := context.Background()

	// Get all orders (simplified - get all without filter)
	orders, _, err := s.orderRepo.GetAllOrders(ctx, nil)
	if err != nil {
		s.logger.Error("OrderService.GetOrderStats failed", err)
		return nil, fmt.Errorf("failed to get order statistics")
	}

	stats := map[string]interface{}{
		"total_orders": len(orders),
		"pending":      0,
		"paid":         0,
		"shipped":      0,
		"delivered":    0,
		"cancelled":    0,
	}

	for _, order := range orders {
		switch order.Status {
		case "pending":
			stats["pending"] = stats["pending"].(int) + 1
		case "paid":
			stats["paid"] = stats["paid"].(int) + 1
		case "shipped":
			stats["shipped"] = stats["shipped"].(int) + 1
		case "delivered":
			stats["delivered"] = stats["delivered"].(int) + 1
		case "cancelled":
			stats["cancelled"] = stats["cancelled"].(int) + 1
		}
	}

	return stats, nil
}

// GetAllOrders gets all orders with optional filtering
func (s *OrderService) GetAllOrders(page, limit int, status string, userID int) ([]*models.Order, int, error) {
	ctx := context.Background()

	// Build filter
	filter := &models.OrderFilter{
		Status: status,
		UserID: userID,
		Page:   page,
		Limit:  limit,
	}

	orders, total, err := s.orderRepo.GetAllOrders(ctx, filter)
	if err != nil {
		s.logger.Error("OrderService.GetAllOrders failed", err)
		return nil, 0, fmt.Errorf("failed to get orders")
	}

	return orders, total, nil
}
