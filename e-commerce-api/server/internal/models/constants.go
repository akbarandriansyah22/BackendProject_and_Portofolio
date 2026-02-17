package models

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
