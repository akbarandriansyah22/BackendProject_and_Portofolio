package handler

import (
	"strconv"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/middleware"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/service"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// OrderHandler handles order-related requests
type OrderHandler struct {
	orderService *service.OrderService
	logger       observability.Logger
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(
	orderService *service.OrderService,
	logger observability.Logger,
) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
	}
}

// GetAll gets all orders for current user
// GET /api/orders
// Protected: Requires authentication
func (h *OrderHandler) GetAll(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Get pagination parameters
	page, limit := utils.GetPaginationParams(c)

	// Get orders
	orders, total, err := h.orderService.GetUserOrders(userID, page, limit)
	if err != nil {
		h.logger.Error("OrderHandler.GetAll: Failed - UserID=%d, Error=%v", userID, err)
		return utils.InternalServerErrorResponse(c, "Failed to get orders")
	}

	return utils.PaginatedSuccessResponse(c, "Orders retrieved successfully", orders, page, limit, int64(total))
}

// ListOrders is an alias for GetAll for route handler
// GET /api/orders
// Protected: Requires authentication
func (h *OrderHandler) ListOrders(c *fiber.Ctx) error {
	return h.GetAll(c)
}

// CreateOrder creates a new order from cart
// POST /api/orders
// Protected: Requires authentication
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	return h.CreateFromCart(c)
}

// GetByID gets order by ID
// GET /api/orders/:id
// Protected: Requires authentication
func (h *OrderHandler) GetByID(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Get order ID from URL parameter
	orderID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid order ID")
	}

	// Get order
	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return utils.NotFoundResponse(c, "Order not found")
		}
		h.logger.Error("OrderHandler.GetByID: Failed - OrderID=%d, Error=%v", orderID, err)
		return utils.InternalServerErrorResponse(c, "Failed to get order")
	}

	// Verify ownership (user can only view their own orders)
	if order.UserID != userID {
		h.logger.Warn("OrderHandler.GetByID: Access denied - UserID=%d attempted to access OrderID=%d", userID, orderID)
		return utils.ForbiddenResponse(c, "Access denied")
	}

	return utils.SuccessResponse(c, "Order retrieved successfully", order)
}

// GetByOrderNumber gets order by order number
// GET /api/orders/number/:orderNumber
// Protected: Requires authentication
func (h *OrderHandler) GetByOrderNumber(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Get order number from URL parameter
	orderNumber := c.Params("orderNumber")
	if orderNumber == "" {
		return utils.BadRequestResponse(c, "Order number is required")
	}

	// Get order
	order, err := h.orderService.GetByOrderNumber(orderNumber)
	if err != nil {
		if err.Error() == "order not found" {
			return utils.NotFoundResponse(c, "Order not found")
		}
		h.logger.Error("OrderHandler.GetByOrderNumber: Failed - OrderNumber=%s, Error=%v", orderNumber, err)
		return utils.InternalServerErrorResponse(c, "Failed to get order")
	}

	// Verify ownership
	if order.UserID != userID {
		h.logger.Warn("OrderHandler.GetByOrderNumber: Access denied - UserID=%d attempted to access order %s", userID, orderNumber)
		return utils.ForbiddenResponse(c, "Access denied")
	}

	return utils.SuccessResponse(c, "Order retrieved successfully", order)
}

// CreateFromCart creates order from cart (checkout)
// POST /api/orders/checkout
// Protected: Requires authentication
//
//	Body: {
//	  "shipping_address": "Jl. Example No. 123",
//	  "payment_method": "bank_transfer"
//	}
func (h *OrderHandler) CreateFromCart(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Parse request body
	var req struct {
		ShippingAddress string `json:"shipping_address"`
		PaymentMethod   string `json:"payment_method"`
		Notes           string `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("OrderHandler.CreateFromCart: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate input
	var validationErrors []utils.ValidationError

	if req.ShippingAddress == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "shipping_address",
			Message: "Shipping address is required",
		})
	}

	if req.PaymentMethod == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "payment_method",
			Message: "Payment method is required",
		})
	}

	// Validate payment method
	validPaymentMethods := map[string]bool{
		"bank_transfer": true,
		"credit_card":   true,
		"e_wallet":      true,
		"cod":           true, // Cash on Delivery
	}

	if !validPaymentMethods[req.PaymentMethod] {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "payment_method",
			Message: "Invalid payment method. Valid options: bank_transfer, credit_card, e_wallet, cod",
		})
	}

	if len(validationErrors) > 0 {
		return utils.ValidationErrorsResponse(c, "Validation failed", validationErrors)
	}

	// Create order from cart
	order, err := h.orderService.CreateFromCart(userID, req.ShippingAddress, req.PaymentMethod, req.Notes)
	if err != nil {
		errMsg := err.Error()
		// Handle specific errors
		if errMsg == "cart is empty" || errMsg == "failed to get cart" {
			return utils.BadRequestResponse(c, errMsg)
		}
		// Stock/validation errors
		return utils.BadRequestResponse(c, err.Error())
	}

	h.logger.Info("Order created from cart: UserID=%d, OrderID=%d, OrderNumber=%s", userID, order.ID, order.OrderNumber)

	return utils.CreatedResponse(c, "Order created successfully", order)
}

// UpdateStatus updates order status (Admin only)
// PUT /api/orders/:id/status
// Protected: Requires admin role
// Body: {"status": "shipped"}
func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {
	// Get order ID from URL parameter
	orderID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid order ID")
	}

	// Parse request body
	var req struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("OrderHandler.UpdateStatus: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"delivered": true,
		"cancelled": true,
	}

	if !validStatuses[req.Status] {
		return utils.BadRequestResponse(c, "Invalid status. Valid options: pending, paid, shipped, delivered, cancelled")
	}

	// Update status
	if err := h.orderService.UpdateStatus(orderID, req.Status); err != nil {
		if err.Error() == "order not found" {
			return utils.NotFoundResponse(c, "Order not found")
		}
		h.logger.Error("OrderHandler.UpdateStatus: Failed - OrderID=%d, Status=%s, Error=%v", orderID, req.Status, err)
		return utils.InternalServerErrorResponse(c, "Failed to update order status")
	}

	// Get updated order
	order, _ := h.orderService.GetByID(orderID)

	h.logger.Info("Order status updated: OrderID=%d, NewStatus=%s", orderID, req.Status)

	return utils.SuccessResponse(c, "Order status updated successfully", order)
}

// CancelOrder cancels an order
// POST /api/orders/:id/cancel
// Protected: Requires authentication
func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Get order ID from URL parameter
	orderID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid order ID")
	}

	// Get order to verify ownership
	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return utils.NotFoundResponse(c, "Order not found")
		}
		return utils.InternalServerErrorResponse(c, "Failed to get order")
	}

	// Verify ownership
	if order.UserID != userID {
		h.logger.Warn("OrderHandler.CancelOrder: Access denied - UserID=%d attempted to cancel OrderID=%d", userID, orderID)
		return utils.ForbiddenResponse(c, "Access denied")
	}

	// Check if order can be cancelled
	if order.Status != "pending" && order.Status != "paid" {
		return utils.BadRequestResponse(c, "Order cannot be cancelled. Current status: "+order.Status)
	}

	// Cancel order
	if err := h.orderService.CancelOrder(orderID); err != nil {
		h.logger.Error("OrderHandler.CancelOrder: Failed - OrderID=%d, Error=%v", orderID, err)
		return utils.InternalServerErrorResponse(c, "Failed to cancel order")
	}

	// Get updated order
	updatedOrder, _ := h.orderService.GetByID(orderID)

	h.logger.Info("Order cancelled: OrderID=%d, UserID=%d", orderID, userID)

	return utils.SuccessResponse(c, "Order cancelled successfully", updatedOrder)
}

// GetOrderStats gets order statistics (Admin only)
// GET /api/orders/stats
// Protected: Requires admin role
func (h *OrderHandler) GetOrderStats(c *fiber.Ctx) error {
	// Get statistics
	stats, err := h.orderService.GetOrderStats()
	if err != nil {
		h.logger.Error("OrderHandler.GetOrderStats: Failed - Error=%v", err)
		return utils.InternalServerErrorResponse(c, "Failed to get order statistics")
	}

	return utils.SuccessResponse(c, "Order statistics retrieved successfully", stats)
}

// GetAllOrders gets all orders (Admin only)
// GET /api/admin/orders
// Protected: Requires admin role
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	// Get pagination parameters
	page, limit := utils.GetPaginationParams(c)

	// Get filter parameters
	status := c.Query("status")     // Filter by status
	userID := c.QueryInt("user_id") // Filter by user

	// Get all orders
	orders, total, err := h.orderService.GetAllOrders(page, limit, status, userID)
	if err != nil {
		h.logger.Error("OrderHandler.GetAllOrders: Failed - Error=%v", err)
		return utils.InternalServerErrorResponse(c, "Failed to get orders")
	}

	return utils.PaginatedSuccessResponse(c, "Orders retrieved successfully", orders, page, limit, int64(total))
}
