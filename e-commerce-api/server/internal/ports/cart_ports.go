package ports

import (
	"context"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

// CartRepository mendefinisikan kontrak untuk data access layer Cart
type CartRepository interface {
	// Create membuat cart baru
	Create(ctx context.Context, userID int) (*models.Cart, error)

	// GetByUserID mengambil cart user
	GetByUserID(ctx context.Context, userID int) (*models.Cart, error)

	// AddItem menambah item ke cart
	AddItem(ctx context.Context, cartID, productID, quantity int) error

	// UpdateItemQuantity memperbarui quantity item di cart
	UpdateItemQuantity(ctx context.Context, cartID, cartItemID, quantity int) error

	// RemoveItem menghapus item dari cart
	RemoveItem(ctx context.Context, cartID, cartItemID int) error

	// GetCartItems mengambil semua item di cart
	GetCartItems(ctx context.Context, cartID int) ([]*models.CartItem, error)

	// ClearCart menghapus semua item di cart
	ClearCart(ctx context.Context, cartID int) error

	// GetItemCount menghitung jumlah item di cart
	GetItemCount(ctx context.Context, cartID int) (int, error)

	// GetCartTotal menghitung total harga cart
	GetCartTotal(ctx context.Context, cartID int) (float64, error)

	// CheckItemExists memeriksa apakah product sudah di cart
	CheckItemExists(ctx context.Context, cartID, productID int) (bool, int, error)

	// Delete menghapus cart
	Delete(ctx context.Context, cartID int) error
}

// CartService mendefinisikan kontrak untuk business logic layer Cart
type CartService interface {
	// GetCart mengambil cart user
	GetCart(ctx context.Context, userID int) (*models.CartResponse, error)

	// AddToCart menambah item ke cart
	AddToCart(ctx context.Context, userID int, req *models.AddToCartRequest) (*models.CartResponse, error)

	// UpdateCartItem memperbarui item di cart
	UpdateCartItem(ctx context.Context, userID, cartItemID int, req *models.UpdateCartItemRequest) (*models.CartResponse, error)

	// RemoveFromCart menghapus item dari cart
	RemoveFromCart(ctx context.Context, userID, cartItemID int) (*models.CartResponse, error)

	// ClearCart menghapus semua item di cart
	ClearCart(ctx context.Context, userID int) error

	// GetCartSummary mengambil ringkasan cart
	GetCartSummary(ctx context.Context, userID int) (*models.CartResponse, error)
}
