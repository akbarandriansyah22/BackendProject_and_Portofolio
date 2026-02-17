package handler

import (
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/middleware"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/service"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *service.AuthService
	logger      observability.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	authService *service.AuthService,
	logger observability.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register handles user registration
// POST /api/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// 1. Parse request body
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("Register: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// 2. Call service (all business logic in service)
	response, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	// 3. Return success response
	return utils.CreatedResponse(c, "User registered successfully", response)
}

// Login handles user login
// POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// 1. Parse request body
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("Login: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// 2. Call service (all business logic in service)
	response, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	// 3. Return success response
	return utils.SuccessResponse(c, "Login successful", response)
}

// GetProfile gets current user profile
// GET /api/auth/profile
// Protected: Requires authentication
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// 1. Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// 2. Call service
	profile, err := h.authService.GetProfile(c.Context(), userID)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	// 3. Return success response
	return utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

// UpdateProfile updates current user profile
// PUT /api/auth/profile
// Protected: Requires authentication
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	// 1. Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// 2. Parse request body
	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("UpdateProfile: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// 3. Call service (all business logic in service)
	profile, err := h.authService.UpdateProfile(c.Context(), userID, req.FullName, req.Email)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	// 4. Return success response
	return utils.SuccessResponse(c, "Profile updated successfully", profile)
}

// ChangePassword changes user password
// PUT /api/auth/change-password
// Protected: Requires authentication
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	// 1. Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// 2. Parse request body
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("ChangePassword: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// 3. Call service (all business logic in service)
	if err := h.authService.ChangePassword(c.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		return h.handleAuthError(c, err)
	}

	// 4. Return success response
	return utils.SuccessMessage(c, "Password changed successfully")
}

// ============================================
// ERROR HANDLING HELPER
// ============================================

// handleAuthError maps service errors to HTTP responses
func (h *AuthHandler) handleAuthError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// Map specific errors to appropriate HTTP status codes
	switch errMsg {
	// Bad Request (400)
	case "invalid email format",
		"password must be at least 8 characters long",
		"password must contain at least one uppercase letter",
		"password must contain at least one lowercase letter",
		"password must contain at least one number",
		"email and password are required",
		"old password and new password are required",
		"old password is incorrect",
		"full name is required":
		return utils.BadRequestResponse(c, errMsg)

	// Unauthorized (401)
	case "invalid email or password":
		return utils.UnauthorizedResponse(c, errMsg)

	// Forbidden (403)
	case "account is inactive":
		return utils.ForbiddenResponse(c, errMsg)

	// Not Found (404)
	case "user not found":
		return utils.NotFoundResponse(c, errMsg)

	// Conflict (409)
	case "email already registered",
		"email already in use":
		return utils.ConflictResponse(c, errMsg)

	// Internal Server Error (500) - default
	default:
		h.logger.Error("AuthHandler: Unhandled error - %v", err)
		return utils.InternalServerErrorResponse(c, "An error occurred")
	}
}
