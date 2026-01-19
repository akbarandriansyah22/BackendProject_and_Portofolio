package models

import (
	"database/sql"
	"time"
)

// struktur tabel pada Roles

type Role struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
// struktur tabel pada Users
type User struct {
	ID           int       `json:"id" db:"id"`
	RoleID       int       `json:"role_id" db:"role_id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name"`
	Phone        sql.NullString   `json:"phone" db:"phone"`
	Address      sql.NullString   `json:"address" db:"address"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	EmailVerifiedAt sql.NullTime `json:"email_verified_at" db:"email_verified_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
// struktur tabel pada Orders
type Order struct {
	ID             int       `json:"id" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	OrderNumber    string    `json:"order_number" db:"order_number"`
	Status         string    `json:"status" db:"status"`
	TotalAmount    float64   `json:"total_amount" db:"total_amount"`
	PaymentMethod  string    `json:"payment_method" db:"payment_method"`
	ShippingAddress string    `json:"shipping_address" db:"shipping_address"`
	ShippingPhone  string    `json:"shipping_phone" db:"shipping_phone"`
	Notes          sql.NullString    `json:"notes" db:"notes"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
// struktur tabel pada OrderItems
type OrderItem struct {
	ID        int       `json:"id" db:"id"`
	OrderID   int       `json:"order_id" db:"order_id"`
	ProductID int       `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	Subtotal  float64   `json:"subtotal" db:"subtotal"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
// struktur tabel pada Carts
type Cart struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
// struktur tabel pada CartItems
type CartItem struct {
	ID        int       `json:"id" db:"id"`
	CartID    int       `json:"cart_id" db:"cart_id"`
	ProductID int       `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
// struktur tabel pada Products
type Product struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description sql.NullString   `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	SKU         sql.NullString   `json:"sku,omitempty" db:"sku"`
	ImageURL    sql.NullString   `json:"image_url" db:"image_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
// struktur tabel pada Categories
type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description sql.NullString   `json:"description" db:"description"`
	ParentID    sql.NullInt32      `json:"parent_id" db:"parent_id"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
// struktur tabel pada ProductCategories
type ProductCategory struct {
	ID         int       `json:"id" db:"id"`
	ProductID  int       `json:"product_id" db:"product_id"`
	CategoryID int       `json:"category_id" db:"category_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
// struktur tabel pada Payments
type Payment struct {
	ID            int       `json:"id" db:"id"`
	OrderID       int       `json:"order_id" db:"order_id"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	Amount        float64     `json:"amount" db:"amount"`
	Status        string    `json:"status" db:"status"`
	TransactionID sql.NullString   `json:"transaction_id,omitempty" db:"transaction_id"`
	PaidAt        sql.NullTime `json:"paid_at" db:"paid_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}


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
	Name        string   `json:"name" validate:"required"`
	Slug        string   `json:"slug" validate:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	Stock       int      `json:"stock" validate:"required,gte=0"`
	SKU         string   `json:"sku"`
	ImageURL    string   `json:"image_url"`
	CategoryIDs []int    `json:"category_ids"`
}

// UpdateProductRequest for updating a product
type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price" validate:"omitempty,gt=0"`
	Stock       int      `json:"stock" validate:"omitempty,gte=0"`
	SKU         string   `json:"sku"`
	ImageURL    string   `json:"image_url"`
	IsActive    *bool    `json:"is_active"`
	CategoryIDs []int    `json:"category_ids"`
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
	Cart  Cart                 `json:"cart"`
	Items []CartItemWithProduct `json:"items"`
	Total float64               `json:"total"`
	ID 	int                   `json:"id"`
	UserID	int                   `json:"user_id"`
	TotalPrice    float64      `json:"total_price"`
	TotalQuantity int          `json:"total_quantity"`
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
	Order Order                 `json:"order"`
	Items []OrderItemWithProduct `json:"items"`
}

// ProductWithCategories includes category details
type ProductWithCategories struct {
	Product    Product   `json:"product"`
	Categories []Category `json:"categories"`
}

// ========================================
// API RESPONSE WRAPPERS
// ========================================

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
// CONSTANTS
// ========================================

// Order status constants
const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"
)

// Payment status constants
const (
	PaymentStatusPending  = "pending"
	PaymentStatusSuccess  = "success"
	PaymentStatusFailed   = "failed"
	PaymentStatusRefunded = "refunded"
)

// Payment method constants
const (
	PaymentMethodCreditCard   = "credit_card"
	PaymentMethodBankTransfer = "bank_transfer"
	PaymentMethodEWallet      = "e_wallet"
	PaymentMethodCOD          = "cod"
)

// Role constants
const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
	RoleSeller   = "seller"
)









