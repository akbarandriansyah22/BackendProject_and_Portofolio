package handler

import (
	"strconv"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/middleware"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/service"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CartHandler handles cart-related requests
type CartHandler struct {
	cartService *service.CartService
	logger      observability.Logger
}

// NewCartHandler creates a new cart handler
func NewCartHandler(
	cartService *service.CartService,
	logger observability.Logger,
) *CartHandler {
	return &CartHandler{
		cartService: cartService,
		logger:      logger,
	}
}

// GetCart gets current user's cart
// GET /api/cart
// Protected: Requires authentication
func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Get cart
	cart, err := h.cartService.GetCart(c.Context(), userID)
	if err != nil {
		h.logger.Error("CartHandler.GetCart: Failed - UserID=%d, Error=%v", userID, err)
		return utils.InternalServerErrorResponse(c, "Failed to get cart")
	}

	return utils.SuccessResponse(c, "Cart retrieved successfully", cart)
}

// AddItem adds item to cart
// POST /api/cart/items
// Protected: Requires authentication
// Body: {"product_id": 1, "quantity": 2}
func (h *CartHandler) AddItem(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Parse request body
	var req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("CartHandler.AddItem: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate input
	var validationErrors []utils.ValidationError

	if req.ProductID <= 0 {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "product_id",
			Message: "Product ID is required and must be greater than 0",
		})
	}

	if req.Quantity <= 0 {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "quantity",
			Message: "Quantity is required and must be greater than 0",
		})
	}

	if len(validationErrors) > 0 {
		return utils.ValidationErrorsResponse(c, "Validation failed", validationErrors)
	}

	// Add item to cart

	if err := h.cartService.AddItem(c.Context(), userID, req.ProductID, req.Quantity); err != nil {
		// Check error type
		errMsg := err.Error()
		if errMsg == "product not found" {
			return utils.NotFoundResponse(c, "Product not found")
		}
		if errMsg == "product is not available" {
			return utils.BadRequestResponse(c, "Product is not available")
		}
		// Stock related errors
		return utils.BadRequestResponse(c, err.Error())
	}

	// Get updated cart
	cart, err := h.cartService.GetCart(c.Context(), userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get updated cart")
	}

	h.logger.Info("Item added to cart: UserID=%d, ProductID=%d, Quantity=%d", userID, req.ProductID, req.Quantity)

	return utils.SuccessResponse(c, "Item added to cart successfully", cart)
}

// RemoveItem removes item from cart
// DELETE /api/cart/items/:id
// Protected: Requires authentication
func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Get item ID from URL parameter
	itemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid item ID")
	}

	// Remove item
	if err := h.cartService.RemoveItem(c.Context(), userID, itemID); err != nil {
		if err.Error() == "cart item not found" {
			return utils.NotFoundResponse(c, "Cart item not found")
		}
		h.logger.Error("CartHandler.RemoveItem: Failed - UserID=%d, ItemID=%d, Error=%v", userID, itemID, err)
		return utils.InternalServerErrorResponse(c, "Failed to remove item")
	}

	// Get updated cart
	cart, err := h.cartService.GetCart(c.Context(), userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get updated cart")
	}

	h.logger.Info("Cart item removed: UserID=%d, ItemID=%d", userID, itemID)

	return utils.SuccessResponse(c, "Item removed successfully", cart)
}

// ClearCart clears all items from cart
// DELETE /api/cart
// Protected: Requires authentication
func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Clear cart
	if err := h.cartService.ClearCart(c.Context(), userID); err != nil {
		h.logger.Error("CartHandler.ClearCart: Failed - UserID=%d, Error=%v", userID, err)
		return utils.InternalServerErrorResponse(c, "Failed to clear cart")
	}

	h.logger.Info("Cart cleared: UserID=%d", userID)

	return utils.SuccessMessage(c, "Cart cleared successfully")
}
