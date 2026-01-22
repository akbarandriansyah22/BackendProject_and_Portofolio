package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/repository"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/utils"
)

// OrderService handles order business logic
type OrderService struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
	paymentRepo *repository.PaymentRepository
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo *repository.OrderRepository,
	cartRepo *repository.CartRepository,
	productRepo *repository.ProductRepository,
	paymentRepo *repository.PaymentRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		paymentRepo: paymentRepo,
	}
}

// GetUserOrders gets all orders for a user with pagination
func (s *OrderService) GetUserOrders(userID, page, limit int) ([]models.Order, int64, error) {
	offset := (page - 1) * limit

	orders, total, err := s.orderRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		utils.Error("OrderService.GetUserOrders: Failed - UserID=%d, Error=%v", userID, err)
		return nil, 0, fmt.Errorf("failed to get orders")
	}

	return orders, total, nil
}

// GetByID gets order by ID
func (s *OrderService) GetByID(orderID int) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

// GetByOrderNumber gets order by order number
func (s *OrderService) GetByOrderNumber(orderNumber string) (*models.Order, error) {
	order, err := s.orderRepo.GetByOrderNumber(orderNumber)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

// CreateFromCart creates order from user's cart
func (s *OrderService) CreateFromCart(
    userID int,
    shippingAddress, paymentMethod, notes string,
) (*models.Order, error) {

    // Get cart
    cart, err := s.cartRepo.GetOrCreateCart(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get cart")
    }

    // Get cart items
    cartItems, err := s.cartRepo.GetCartItems(cart.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to get cart items")
    }

    if len(cartItems) == 0 {
        return nil, fmt.Errorf("cart is empty")
    }

    // Validate & calculate total
    var totalAmount float64
    for _, item := range cartItems {
        product, err := s.productRepo.GetByID(item.ProductID)
        if err != nil {
            return nil, fmt.Errorf("product is no longer available")
        }

        if !product.IsActive {
            return nil, fmt.Errorf("product '%s' is not available", product.Name)
        }

        if product.Stock < item.Quantity {
            return nil, fmt.Errorf(
                "insufficient stock for '%s'. Available: %d, Requested: %d",
                product.Name, product.Stock, item.Quantity,
            )
        }

        totalAmount += item.Price * float64(item.Quantity)
    }

    // Create order
    order := &models.Order{
        UserID:          userID,
        OrderNumber:     s.generateOrderNumber(),
        TotalAmount:     totalAmount,
        Status:          "pending",
        ShippingAddress: shippingAddress,
        PaymentMethod:   paymentMethod,
        Notes:           sql.NullString{String: notes, Valid: notes != ""},
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    if err := s.orderRepo.Create(order); err != nil {
        utils.Error("OrderService.CreateFromCart: create order failed: %v", err)
        return nil, fmt.Errorf("failed to create order")
    }

    // Create order items + decrement stock
    for _, cartItem := range cartItems {

        orderItem := &models.OrderItem{
            OrderID:   order.ID,
            ProductID: cartItem.ProductID,
            Quantity:  cartItem.Quantity,
            Price:     cartItem.Price,
            Subtotal:  cartItem.Price * float64(cartItem.Quantity),
        }

        if err := s.orderRepo.AddOrderItem(orderItem); err != nil {
            utils.Error("add order item failed (order_id=%d): %v", order.ID, err)

            if rbErr := s.orderRepo.Delete(order.ID); rbErr != nil {
                utils.Error("rollback failed (order_id=%d): %v", order.ID, rbErr)
            }
            return nil, fmt.Errorf("failed to create order items")
        }

        if err := s.productRepo.DecrementStock(cartItem.ProductID, cartItem.Quantity); err != nil {
            utils.Error("decrement stock failed (product_id=%d): %v", cartItem.ProductID, err)

            if rbErr := s.orderRepo.Delete(order.ID); rbErr != nil {
                utils.Error("rollback failed (order_id=%d): %v", order.ID, rbErr)
            }
            return nil, fmt.Errorf("failed to update stock")
        }
    }

    // Clear cart (non-critical)
    if err := s.cartRepo.ClearCart(cart.ID); err != nil {
        utils.Error("clear cart failed (cart_id=%d): %v", cart.ID, err)
    }

    fullOrder, err := s.orderRepo.GetByID(order.ID)
    if err != nil {
        return nil, fmt.Errorf("order created but failed to fetch details")
    }

    return fullOrder, nil
}


// UpdateStatus updates order status
func (s *OrderService) UpdateStatus(orderID int, status string) error {
	// Check if order exists
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	// Update status
	if err := s.orderRepo.UpdateStatus(orderID, status); err != nil {
		utils.Error("OrderService.UpdateStatus: Failed - OrderID=%d, Status=%s, Error=%v", orderID, status, err)
		return fmt.Errorf("failed to update order status")
	}

	utils.Info("Order status updated: OrderID=%d, OldStatus=%s, NewStatus=%s", orderID, order.Status, status)
	return nil
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(orderID int) error {
	// Get order
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	// Check if order can be cancelled
	if order.Status != "pending" && order.Status != "paid" {
		return fmt.Errorf("order cannot be cancelled. Current status: %s", order.Status)
	}

	// Get order items to restore stock
	orderItems, err := s.orderRepo.GetOrderItems(orderID)
	if err != nil {
		utils.Warn("OrderService.CancelOrder: Failed to get order items - OrderID=%d", orderID)
	} else {
		// Restore stock for each item
		for _, item := range orderItems {
			if err := s.productRepo.IncrementStock(item.ProductID, item.Quantity); err != nil {
				utils.Error("OrderService.CancelOrder: Failed to restore stock - ProductID=%d, Error=%v", item.ProductID, err)
				// Continue anyway
			}
		}
	}

	// Update order status to cancelled
	if err := s.orderRepo.UpdateStatus(orderID, "cancelled"); err != nil {
		utils.Error("OrderService.CancelOrder: Failed to update status - OrderID=%d, Error=%v", orderID, err)
		return fmt.Errorf("failed to cancel order")
	}

	utils.Info("Order cancelled: OrderID=%d", orderID)
	return nil
}

// GetOrderStats gets order statistics
func (s *OrderService) GetOrderStats() (map[string]interface{}, error) {
	// Get total orders
	totalOrders, err := s.orderRepo.CountTotal()
	if err != nil {
		utils.Error("OrderService.GetOrderStats: Failed to count total - Error=%v", err)
		return nil, fmt.Errorf("failed to get order statistics")
	}

	// Get count by status
	pendingCount, _ := s.orderRepo.CountByStatus("pending")
	paidCount, _ := s.orderRepo.CountByStatus("paid")
	shippedCount, _ := s.orderRepo.CountByStatus("shipped")
	deliveredCount, _ := s.orderRepo.CountByStatus("delivered")
	cancelledCount, _ := s.orderRepo.CountByStatus("cancelled")

	// Get total revenue
	totalRevenue, _ := s.orderRepo.GetTotalRevenue()

	// Get today's orders
	todayOrders, _ := s.orderRepo.CountTodayOrders()
	todayRevenue, _ := s.orderRepo.GetTodayRevenue()

	stats := map[string]interface{}{
		"total_orders":      totalOrders,
		"pending_orders":    pendingCount,
		"paid_orders":       paidCount,
		"shipped_orders":    shippedCount,
		"delivered_orders":  deliveredCount,
		"cancelled_orders":  cancelledCount,
		"total_revenue":     totalRevenue,
		"today_orders":      todayOrders,
		"today_revenue":     todayRevenue,
	}

	return stats, nil
}

// GetAllOrders gets all orders with pagination and filters (for admin)
// ✅ FIX: Rename parameter 'something' to 'userID'
func (s *OrderService) GetAllOrders(page, limit int, status string, userID int) ([]models.Order, int64, error) {
	offset := (page - 1) * limit

	// Get orders with filters
	filters := make(map[string]interface{})
	if status != "" {
		filters["status"] = status
	}
	// ✅ FIX: Now 'userID' is properly defined and used
	if userID > 0 {
		filters["user_id"] = userID
	}

	orders, total, err := s.orderRepo.GetAllWithFilters(filters, limit, offset)
	if err != nil {
		utils.Error("OrderService.GetAllOrders: Failed - Error=%v", err)
		return nil, 0, fmt.Errorf("failed to get orders")
	}

	return orders, total, nil
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// generateOrderNumber generates a unique order number
func (s *OrderService) generateOrderNumber() string {
	// Format: ORD-YYYYMMDD-XXX
	now := time.Now()
	date := now.Format("20060102")
	
	// Get count of orders today
	todayCount, _ := s.orderRepo.CountTodayOrders()
	
	// Generate number with padding
	number := todayCount + 1
	orderNumber := fmt.Sprintf("ORD-%s-%03d", date, number)
	
	return orderNumber
}