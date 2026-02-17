package models

import "time"

// ========================================
// REQUEST DTOs (Data Transfer Objects)
// ========================================

// RegisterRequest for user registration
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

// LoginRequest for user authentication
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// CreateRoleRequest for creating a new role
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// UpdateRoleRequest for updating a role
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateProductRequest for creating a new product
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Slug        string  `json:"slug" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	SKU         string  `json:"sku"`
	ImageURL    string  `json:"image_url"`
	CategoryIDs []int   `json:"category_ids"`
}

// UpdateProductRequest for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Stock       int     `json:"stock" validate:"omitempty,gte=0"`
	SKU         string  `json:"sku"`
	ImageURL    string  `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
	CategoryIDs []int   `json:"category_ids"`
}

// CreateCategoryRequest for creating a new category
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
	ParentID    *int   `json:"parent_id"`
}

// UpdateCategoryRequest for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    *int   `json:"parent_id"`
	IsActive    *bool  `json:"is_active"`
}

// AddToCartRequest for adding an item to cart
type AddToCartRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
}

// UpdateCartItemRequest for updating cart item quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

// CreateOrderRequest for creating a new order
type CreateOrderRequest struct {
	ShippingAddress string `json:"shipping_address" validate:"required"`
	ShippingPhone   string `json:"shipping_phone" validate:"required"`
	Notes           string `json:"notes"`
}

// UpdateOrderStatusRequest for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

// CreatePaymentRequest for processing payment
type CreatePaymentRequest struct {
	OrderID       int    `json:"order_id" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
}

// UpdateProfileRequest for updating user profile
type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// ChangePasswordRequest for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// ConfirmPaymentRequest for confirming payment
type ConfirmPaymentRequest struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	ProofURL      string `json:"proof_url"`
}

// PaymentCallbackRequest for payment gateway callback
type PaymentCallbackRequest struct {
	TransactionID string  `json:"transaction_id" validate:"required"`
	OrderID       int     `json:"order_id" validate:"required"`
	Status        string  `json:"status" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
}

// ========================================
// FILTER DTOs
// ========================================

// OrderFilter for filtering orders
type OrderFilter struct {
	UserID    int        `json:"user_id,omitempty"`
	Status    string     `json:"status,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Page      int        `json:"page,omitempty"`
	Limit     int        `json:"limit,omitempty"`
}

// PaymentFilter for filtering payments
type PaymentFilter struct {
	OrderID   int        `json:"order_id,omitempty"`
	Status    string     `json:"status,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Page      int        `json:"page,omitempty"`
	Limit     int        `json:"limit,omitempty"`
}

// ProductFilter for filtering products
type ProductFilter struct {
	CategoryID int     `json:"category_id,omitempty"`
	Keyword    string  `json:"keyword,omitempty"`
	MinPrice   float64 `json:"min_price,omitempty"`
	MaxPrice   float64 `json:"max_price,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	Page       int     `json:"page,omitempty"`
	Limit      int     `json:"limit,omitempty"`
}
