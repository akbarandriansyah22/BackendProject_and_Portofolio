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
	
	AuditModel
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












