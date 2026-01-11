package service

import (
	"fmt"

	"e-commerce-api/server/internal/models"
	"e-commerce-api/server/internal/repository"
	"e-commerce-api/server/internal/utils"
)

// CartService handles cart business logic
type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

// NewCartService creates a new cart service
func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// GetOrCreateCart gets or creates a cart for user
func (s *CartService) GetOrCreateCart(userID int) (*models.CartResponse, error) {
	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		utils.Error("CartService.GetOrCreateCart: Failed - UserID=%d, Error=%v", userID, err)
		return nil, fmt.Errorf("failed to get cart")
	}

	// Get cart items with product details
	items, err := s.cartRepo.GetCartItems(cart.ID)
	if err != nil {
		utils.Error("CartService.GetOrCreateCart: Failed to get items - CartID=%d, Error=%v", cart.ID, err)
		return nil, fmt.Errorf("failed to get cart items")
	}

	// Calculate totals
	totalPrice, totalQuantity := s.calculateTotals(items)

	return &models.CartResponse{
		ID:            cart.ID,
		UserID:        cart.UserID,
		Items:         items,
		TotalPrice:    totalPrice,
		TotalQuantity: totalQuantity,
	}, nil
}

// AddItem adds item to cart
func (s *CartService) AddItem(userID, productID, quantity int) (*models.CartResponse, error) {
	// Validate quantity
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	// Get product
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	// Check if product is active
	if !product.IsActive {
		return nil, fmt.Errorf("product is not available")
	}

	// Check stock availability
	if product.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock. Available: %d", product.Stock)
	}

	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		utils.Error("CartService.AddItem: Failed to get cart - UserID=%d, Error=%v", userID, err)
		return nil, fmt.Errorf("failed to add item to cart")
	}

	// Check if item already exists in cart
	existingItem, err := s.cartRepo.GetCartItemByCartAndProduct(cart.ID, productID)
	if err == nil && existingItem != nil {
		// Item exists, update quantity
		newQuantity := existingItem.Quantity + quantity
		
		// Check stock for new quantity
		if product.Stock < newQuantity {
			return nil, fmt.Errorf("insufficient stock. Available: %d, Already in cart: %d", product.Stock, existingItem.Quantity)
		}

		if err := s.cartRepo.UpdateItemQuantity(existingItem.ID, newQuantity); err != nil {
			utils.Error("CartService.AddItem: Failed to update quantity - Error=%v", err)
			return nil, fmt.Errorf("failed to add item to cart")
		}
	} else {
		// Item doesn't exist, add new
		cart := &models.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			Quantity:  quantity,
			Price:     product.Price,
		}

		if err := s.cartRepo.AddItem(cart.ID, productID, quantity, product.Price); err != nil {
			utils.Error("CartService.AddItem: Failed to add item - Error=%v", err)
			return nil, fmt.Errorf("failed to add item to cart")
		}
	}

	utils.Info("Item added to cart: UserID=%d, ProductID=%d, Quantity=%d", userID, productID, quantity)

	// Return updated cart
	return s.GetOrCreateCart(userID)
}

// UpdateItemQuantity updates cart item quantity
func (s *CartService) UpdateItemQuantity(userID, itemID, quantity int) (*models.CartResponse, error) {
	// Validate quantity
	if quantity < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}

	// If quantity is 0, remove item
	if quantity == 0 {
		return s.RemoveItem(userID, itemID)
	}

	// Get cart item
	item, err := s.cartRepo.GetCartItemByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("cart item not found")
	}

	// Verify ownership (get cart to check userID)
	cart, err := s.cartRepo.GetByID(item.CartID)
	if err != nil || cart.UserID != userID {
		return nil, fmt.Errorf("cart item not found")
	}

	// Get product to check stock
	product, err := s.productRepo.GetByID(item.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	// Check stock availability
	if product.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock. Available: %d", product.Stock)
	}

	// Update quantity
	if err := s.cartRepo.UpdateItemQuantity(itemID, quantity); err != nil {
		utils.Error("CartService.UpdateItemQuantity: Failed - ItemID=%d, Error=%v", itemID, err)
		return nil, fmt.Errorf("failed to update cart item")
	}

	utils.Info("Cart item updated: UserID=%d, ItemID=%d, NewQuantity=%d", userID, itemID, quantity)

	// Return updated cart
	return s.GetOrCreateCart(userID)
}

// RemoveItem removes item from cart
func (s *CartService) RemoveItem(userID, itemID int) (*models.CartResponse, error) {
	// Get cart item
	item, err := s.cartRepo.GetCartItemByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("cart item not found")
	}

	// Verify ownership
	cart, err := s.cartRepo.GetByID(item.CartID)
	if err != nil || cart.UserID != userID {
		return nil, fmt.Errorf("cart item not found")
	}

	// Remove item
	if err := s.cartRepo.RemoveItem(itemID); err != nil {
		utils.Error("CartService.RemoveItem: Failed - ItemID=%d, Error=%v", itemID, err)
		return nil, fmt.Errorf("failed to remove cart item")
	}

	utils.Info("Cart item removed: UserID=%d, ItemID=%d", userID, itemID)

	// Return updated cart
	return s.GetOrCreateCart(userID)
}

// ClearCart clears all items from cart
func (s *CartService) ClearCart(userID int) error {
	// Get cart
	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart")
	}

	// Clear cart
	if err := s.cartRepo.ClearCart(cart.ID); err != nil {
		utils.Error("CartService.ClearCart: Failed - UserID=%d, Error=%v", userID, err)
		return fmt.Errorf("failed to clear cart")
	}

	utils.Info("Cart cleared: UserID=%d", userID)
	return nil
}

// SyncCartPrices syncs cart item prices with current product prices
func (s *CartService) SyncCartPrices(userID int) error {
	// Get cart
	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart")
	}

	// Sync prices
	if err := s.cartRepo.SyncCartItemPrices(cart.ID); err != nil {
		utils.Error("CartService.SyncCartPrices: Failed - UserID=%d, Error=%v", userID, err)
		return fmt.Errorf("failed to sync cart prices")
	}

	utils.Info("Cart prices synced: UserID=%d", userID)
	return nil
}

// ValidateCartForCheckout validates cart before checkout
func (s *CartService) ValidateCartForCheckout(userID int) error {
	// Get cart
	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart")
	}

	// Get cart items
	items, err := s.cartRepo.GetCartItems(cart.ID)
	if err != nil {
		return fmt.Errorf("failed to get cart items")
	}

	// Check if cart is empty
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}

	// Validate each item
	for _, item := range items {
		// Get current product info
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return fmt.Errorf("product '%s' is no longer available", item.Product.Name)
		}

		// Check if product is active
		if !product.IsActive {
			return fmt.Errorf("product '%s' is no longer available", product.Name)
		}

		// Check stock
		if product.Stock < item.Quantity {
			return fmt.Errorf("insufficient stock for '%s'. Available: %d, Requested: %d", product.Name, product.Stock, item.Quantity)
		}

		// Check if price changed significantly (more than 10%)
		priceChange := (item.Price - product.Price) / product.Price * 100
		if priceChange > 10 || priceChange < -10 {
			return fmt.Errorf("price for '%s' has changed significantly. Please review your cart", product.Name)
		}
	}

	return nil
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// calculateTotals calculates total price and quantity
func (s *CartService) calculateTotals(items []models.CartItemWithProduct) (totalPrice float64, totalQuantity int) {
	for _, item := range items {
		totalPrice += item.Price * float64(item.Quantity)
		totalQuantity += item.Quantity
	}
	return totalPrice, totalQuantity
}