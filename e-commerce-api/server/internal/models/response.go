package models

import (
	"time"
)

// ========================================
// RESPONSE DTOs
// ========================================

// UserResponse is the user response without sensitive data
type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone,omitempty"`
	Address   string    `json:"address,omitempty"`
	RoleID    int       `json:"role_id"`
	RoleName  string    `json:"role_name,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginResponse after successful login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// CartItemWithProduct includes product details
type CartItemWithProduct struct {
	ID        int       `json:"id"`
	CartID    int       `json:"cart_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Product   Product   `json:"product"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CartResponse includes cart items and total
type CartResponse struct {
	ID            int                   `json:"id"`
	UserID        int                   `json:"user_id"`
	Items         []CartItemWithProduct `json:"items"`
	TotalPrice    float64               `json:"total_price"`
	TotalQuantity int                   `json:"total_quantity"`
}

// OrderItemWithProduct includes product details
type OrderItemWithProduct struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Subtotal  float64   `json:"subtotal"`
	Product   Product   `json:"product"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderWithItems includes order items
type OrderWithItems struct {
	Order Order                  `json:"order"`
	Items []OrderItemWithProduct `json:"items"`
}

// ProductWithCategories includes category details
type ProductWithCategories struct {
	Product    Product    `json:"product"`
	Categories []Category `json:"categories"`
}

// ========================================
// API RESPONSE WRAPPERS
// ========================================

// information Error
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIResponse is the standard response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse for paginated data
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int64       `json:"total"`
}

// ========================================
// RESPONSE DTOs - DOMAIN SPECIFIC
// ========================================

// RoleResponse for role data
type RoleResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductResponse for product listing
type ProductResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	ImageURL  string    `json:"image_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductDetailResponse for detailed product info
type ProductDetailResponse struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Description string             `json:"description,omitempty"`
	Price       float64            `json:"price"`
	Stock       int                `json:"stock"`
	SKU         string             `json:"sku,omitempty"`
	ImageURL    string             `json:"image_url,omitempty"`
	IsActive    bool               `json:"is_active"`
	Categories  []CategoryResponse `json:"categories,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// CategoryResponse for category data
type CategoryResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	ParentID    *int      `json:"parent_id,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryDetailResponse for detailed category info
type CategoryDetailResponse struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Description string             `json:"description,omitempty"`
	ParentID    *int               `json:"parent_id,omitempty"`
	IsActive    bool               `json:"is_active"`
	Children    []CategoryResponse `json:"children,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// CategoryTreeResponse for hierarchical category data
type CategoryTreeResponse struct {
	ID       int                     `json:"id"`
	Name     string                  `json:"name"`
	Slug     string                  `json:"slug"`
	IsActive bool                    `json:"is_active"`
	Children []*CategoryTreeResponse `json:"children,omitempty"`
}

// OrderResponse for order listing
type OrderResponse struct {
	ID              int       `json:"id"`
	OrderNumber     string    `json:"order_number"`
	Status          string    `json:"status"`
	TotalAmount     float64   `json:"total_amount"`
	PaymentMethod   string    `json:"payment_method"`
	ShippingAddress string    `json:"shipping_address"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// OrderDetailResponse for detailed order info
type OrderDetailResponse struct {
	ID              int                    `json:"id"`
	OrderNumber     string                 `json:"order_number"`
	UserID          int                    `json:"user_id"`
	Status          string                 `json:"status"`
	TotalAmount     float64                `json:"total_amount"`
	PaymentMethod   string                 `json:"payment_method"`
	ShippingAddress string                 `json:"shipping_address"`
	ShippingPhone   string                 `json:"shipping_phone"`
	Notes           string                 `json:"notes,omitempty"`
	Items           []OrderItemWithProduct `json:"items,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// PaymentResponse for payment data
type PaymentResponse struct {
	ID            int       `json:"id"`
	OrderID       int       `json:"order_id"`
	PaymentMethod string    `json:"payment_method"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PaymentDetailResponse for detailed payment info
type PaymentDetailResponse struct {
	ID            int        `json:"id"`
	OrderID       int        `json:"order_id"`
	PaymentMethod string     `json:"payment_method"`
	Amount        float64    `json:"amount"`
	Status        string     `json:"status"`
	TransactionID string     `json:"transaction_id,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// PaymentStatusResponse for payment status
type PaymentStatusResponse struct {
	PaymentID     int        `json:"payment_id"`
	OrderID       int        `json:"order_id"`
	Status        string     `json:"status"`
	Amount        float64    `json:"amount"`
	TransactionID string     `json:"transaction_id,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// UserProfileResponse for user profile
type UserProfileResponse struct {
	ID              int        `json:"id"`
	Email           string     `json:"email"`
	Name            string     `json:"name"`
	Phone           string     `json:"phone,omitempty"`
	Address         string     `json:"address,omitempty"`
	RoleID          int        `json:"role_id"`
	RoleName        string     `json:"role_name,omitempty"`
	IsActive        bool       `json:"is_active"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
