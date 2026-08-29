package utils

import (
	"github.com/gofiber/fiber/v2"
)

// ============================================
// STANDARD RESPONSE STRUCTURES
// ============================================

// Response represents standard API response
// Struktur standar untuk semua response API
type Response struct {
	Success bool        `json:"success"`           // Status: true/false
	Message string      `json:"message,omitempty"` // Pesan untuk user
	Data    interface{} `json:"data,omitempty"`    // Data payload (optional)
	Error   string      `json:"error,omitempty"`   // Error message (optional)
}

// PaginatedResponse represents paginated API response
// Struktur untuk response dengan pagination
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"` // Array of items
	Meta    Pagination  `json:"meta"` // Pagination info
}

// Pagination represents pagination metadata
// Metadata untuk pagination (page, limit, total, dll)
type Pagination struct {
	Page       int   `json:"page"`        // Current page (1, 2, 3, ...)
	Limit      int   `json:"limit"`       // Items per page
	TotalItems int64 `json:"total_items"` // Total items in database
	TotalPages int   `json:"total_pages"` // Total pages
}

// ValidationError represents field validation error
// Struktur untuk validation error (field-level)
type ValidationError struct {
	Field   string `json:"field"`   // Field name yang error
	Message string `json:"message"` // Error message
}

// ValidationErrorResponse represents validation error response
// Response untuk validation errors (multiple fields)
type ValidationErrorResponse struct {
	Success bool              `json:"success"` // Always false
	Message string            `json:"message"` // General error message
	Errors  []ValidationError `json:"errors"`  // Array of field errors
}

// ============================================
// SUCCESS RESPONSES
// ============================================

// Success sends a success response with data
// Fungsi untuk kirim response sukses dengan data
// Contoh: SuccessResponse(c, "Product created", product)
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessMessage sends a success response without data
// Fungsi untuk kirim response sukses tanpa data (hanya pesan)
// Contoh: SuccessMessage(c, "Product deleted successfully")
func SuccessMessage(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
	})
}

// CreatedResponse sends a 201 Created response
// Fungsi untuk kirim response 201 (created)
// Contoh: CreatedResponse(c, "User registered", user)
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ============================================
// ERROR RESPONSES
// ============================================

// ErrorResponse sends an error response
// Fungsi untuk kirim error response dengan status code custom
// Contoh: ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// BadRequestResponse sends a 400 Bad Request response
// Fungsi untuk kirim error 400 (bad request)
// Contoh: BadRequestResponse(c, "Invalid email format")
func BadRequestResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// UnauthorizedResponse sends a 401 Unauthorized response
// Fungsi untuk kirim error 401 (unauthorized - belum login)
// Contoh: UnauthorizedResponse(c, "Please login first")
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// ForbiddenResponse sends a 403 Forbidden response
// Fungsi untuk kirim error 403 (forbidden - tidak ada akses)
// Contoh: ForbiddenResponse(c, "Admin only")
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// NotFoundResponse sends a 404 Not Found response
// Fungsi untuk kirim error 404 (not found)
// Contoh: NotFoundResponse(c, "Product not found")
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// ConflictResponse sends a 409 Conflict response
// Fungsi untuk kirim error 409 (conflict - duplicate data)
// Contoh: ConflictResponse(c, "Email already exists")
func ConflictResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// InternalServerErrorResponse sends a 500 Internal Server Error response
// Fungsi untuk kirim error 500 (server error)
// Contoh: InternalServerErrorResponse(c, "Database connection failed")
func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// ============================================
// VALIDATION ERROR RESPONSE
// ============================================

// ValidationErrorsResponse sends a validation error response
// Fungsi untuk kirim validation errors (multiple fields)
// Contoh:
//
//	errors := []ValidationError{
//	  {Field: "email", Message: "Invalid email format"},
//	  {Field: "password", Message: "Password too short"},
//	}
//	ValidationErrorsResponse(c, "Validation failed", errors)
func ValidationErrorsResponse(c *fiber.Ctx, message string, errors []ValidationError) error {
	return c.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// SingleValidationError creates a single validation error response
// Fungsi shortcut untuk kirim 1 validation error
// Contoh: SingleValidationError(c, "email", "Invalid email format")
func SingleValidationError(c *fiber.Ctx, field, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
		Success: false,
		Message: "Validation failed",
		Errors: []ValidationError{
			{Field: field, Message: message},
		},
	})
}

// ============================================
// PAGINATED RESPONSE
// ============================================

// PaginatedSuccessResponse sends a paginated response
// Fungsi untuk kirim response dengan pagination
// Contoh: PaginatedSuccessResponse(c, "Products retrieved", products, page, limit, total)
func PaginatedSuccessResponse(c *fiber.Ctx, message string, data interface{}, page, limit int, total int64) error {
	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	// Handle edge case: no data
	if total == 0 {
		totalPages = 0
	}

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: totalPages,
		},
	})
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// GetPaginationParams extracts pagination parameters from query
// Fungsi helper untuk extract page & limit dari query string
// Default: page=1, limit=10
// Contoh: page, limit := GetPaginationParams(c)
func GetPaginationParams(c *fiber.Ctx) (page, limit int) {
	// Get page (default = 1)
	page = c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	// Get limit (default = 10, max = 100)
	limit = c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Prevent abuse
	}

	return page, limit
}

// CalculateOffset calculates database offset from page and limit
// Fungsi helper untuk hitung offset untuk SQL query
// Contoh: offset := CalculateOffset(2, 10) // Result: 10
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// ============================================
// CUSTOM RESPONSE BUILDERS
// ============================================

// NewResponse creates a new Response
// Builder untuk custom response
func NewResponse(success bool, message string, data interface{}, errorMsg string) Response {
	return Response{
		Success: success,
		Message: message,
		Data:    data,
		Error:   errorMsg,
	}
}

// NewPaginatedResponse creates a new PaginatedResponse
// Builder untuk custom paginated response
func NewPaginatedResponse(success bool, message string, data interface{}, page, limit int, total int64) PaginatedResponse {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	if total == 0 {
		totalPages = 0
	}

	return PaginatedResponse{
		Success: success,
		Message: message,
		Data:    data,
		Meta: Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}
}

// ============================================
// COMMON ERROR MESSAGES (CONSTANTS)
// ============================================

const (
	// General errors
	ErrInternalServer = "Internal server error"
	ErrInvalidInput   = "Invalid input"
	ErrUnauthorized   = "Unauthorized"
	ErrForbidden      = "Forbidden"
	ErrNotFound       = "Resource not found"

	// Auth errors
	ErrInvalidCredentials = "Invalid email or password"
	ErrEmailExists        = "Email already exists"
	ErrInvalidToken       = "Invalid or expired access credential"

	// Validation errors
	ErrValidationFailed = "Validation failed"
	ErrInvalidEmail     = "Invalid email format"
	ErrPasswordTooShort = "Password must be at least 8 characters"

	// Database errors
	ErrDatabaseConnection = "Database connection failed"
	ErrRecordNotFound     = "Record not found"
	ErrDuplicateEntry     = "Duplicate entry"
)
