package service

import (
	"context"
	"fmt"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/ports"
)

// CartService handles cart business logic
type CartService struct {
	cartRepo    ports.CartRepository
	productRepo ports.ProductRepository
	logger      observability.Logger
}

// NewCartService creates a new cart service
func NewCartService(
	cartRepo ports.CartRepository,
	productRepo ports.ProductRepository,
	logger observability.Logger,
) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		logger:      logger,
	}
}

// GetCart gets user's cart
func (s *CartService) GetCart(ctx context.Context, userID int) (*models.CartResponse, error) {
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("CartService.GetCart failed", err)
		return nil, fmt.Errorf("failed to get cart")
	}

	if cart == nil {
		return nil, fmt.Errorf("cart not found")
	}

	_, err = s.cartRepo.GetCartItems(ctx, cart.ID)
	if err != nil {
		s.logger.Error("CartService.GetCart failed to get items", err)
		return nil, fmt.Errorf("failed to get cart items")
	}

	return &models.CartResponse{
		ID:            cart.ID,
		UserID:        cart.UserID,
		Items:         make([]models.CartItemWithProduct, 0),
		TotalPrice:    0,
		TotalQuantity: 0,
	}, nil
}

// AddItem adds item to cart
func (s *CartService) AddItem(ctx context.Context, userID, productID, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	// Get product
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil || product == nil {
		return fmt.Errorf("product not found")
	}

	// Get or create cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		cart, err = s.cartRepo.Create(ctx, userID)
		if err != nil {
			s.logger.Error("CartService.AddItem failed to create cart", err)
			return fmt.Errorf("failed to create cart")
		}
	}

	return s.cartRepo.AddItem(ctx, cart.ID, productID, quantity)
}

// RemoveItem removes item from cart
func (s *CartService) RemoveItem(ctx context.Context, userID, cartItemID int) error {
	return s.cartRepo.RemoveItem(ctx, 0, cartItemID)
}

// ClearCart clears all items from cart
func (s *CartService) ClearCart(ctx context.Context, userID int) error {
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil || cart == nil {
		return fmt.Errorf("cart not found")
	}

	return s.cartRepo.ClearCart(ctx, cart.ID)
}
