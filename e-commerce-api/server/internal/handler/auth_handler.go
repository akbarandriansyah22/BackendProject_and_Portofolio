package handler

import (
	"strings"

	"e-commerce-api/server/internal/middleware"
	"e-commerce-api/server/internal/models"
	"e-commerce-api/server/internal/repository"
	"e-commerce-api/server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userRepo  *repository.UserRepository
	roleRepo  *repository.RoleRepository
	jwtSecret string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration
// POST /api/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Warn("Register: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate input
	var validationErrors []utils.ValidationError

	// Validate email
	if req.Email == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	} else if !isValidEmail(req.Email) {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}

	// Validate password
	if req.Password == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	} else if err := utils.ValidatePasswordStrength(req.Password); err != nil {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "password",
			Message: err.Error(),
		})
	}

	// Validate full name
	if req.FullName == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field:   "full_name",
			Message: "Full name is required",
		})
	}

	// Return validation errors if any
	if len(validationErrors) > 0 {
		return utils.ValidationErrorsResponse(c, "Validation failed", validationErrors)
	}

	// Check if email already exists
	exists, err := h.userRepo.EmailExists(req.Email)
	if err != nil {
		utils.Error("Register: Failed to check email existence - %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to register user")
	}
	if exists {
		return utils.ConflictResponse(c, "Email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error("Register: Failed to hash password - %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to register user")
	}

	// Get customer role (default role)
	// Assume roleID 2 = customer (you can adjust this)
	customerRoleID := 2

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:     req.FullName,
		RoleID:       customerRoleID,
		IsActive:     true,
	}

	if err := h.userRepo.Create(user); err != nil {
		utils.Error("Register: Failed to create user - %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to register user")
	}

	utils.Info("User registered successfully: ID=%d, Email=%s", user.ID, user.Email)

	// Generate JWT token
	token, err := utils.GenerateAccessToken(user.ID, user.Email, user.RoleID, user.Name, h.jwtSecret)
	if err != nil {
		utils.Error("Register: Failed to generate token - %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to generate token")
	}

	// Return success response
	return utils.CreatedResponse(c, "User registered successfully", fiber.Map{
		"token": token,
		"user": models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.Name,
			RoleID:   user.RoleID,
			IsActive: user.IsActive,
		},
	})
}

// Login handles user login
// POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Warn("Login: Invalid request body - %v", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return utils.BadRequestResponse(c, "Email and password are required")
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		utils.Warn("Login: Failed login attempt for email=%s", req.Email)
		return utils.UnauthorizedResponse(c, "Invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		utils.Warn("Login: Inactive user attempted login - UserID=%d", user.ID)
		return utils.ForbiddenResponse(c, "Account is inactive")
	}

	// Verify password
	if !utils.VerifyPassword(req.Password, user.PasswordHash) {
		utils.Warn("Login: Invalid password for email=%s", req.Email)
		return utils.UnauthorizedResponse(c, "Invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateAccessToken(user.ID, user.Email, user.RoleID, user.Name, h.jwtSecret)
	if err != nil {
		utils.Error("Login: Failed to generate token - %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to generate token")
	}

	utils.Info("User logged in successfully: UserID=%d, Email=%s", user.ID, user.Email)

	// Return success response
	return utils.SuccessResponse(c, "Login successful", models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.Name,
			RoleID:   user.RoleID,
			IsActive: user.IsActive,
		},
	})
}

// GetProfile gets current user profile
// GET /api/auth/profile
// Protected: Requires authentication
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Get user from database
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		utils.Error("GetProfile: User not found - UserID=%d", userID)
		return utils.NotFoundResponse(c, "User not found")
	}

	// Return user profile
	return utils.SuccessResponse(c, "Profile retrieved successfully", models.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.Name,
		RoleID:   user.RoleID,
		IsActive: user.IsActive,
	})
}

// UpdateProfile updates current user profile
// PUT /api/auth/profile
// Protected: Requires authentication
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Parse request body
	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Get current user
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return utils.NotFoundResponse(c, "User not found")
	}

	// Validate email if changed
	if req.Email != "" && req.Email != user.Email {
		if !isValidEmail(req.Email) {
			return utils.BadRequestResponse(c, "Invalid email format")
		}

		// Check if new email already exists
		exists, err := h.userRepo.EmailExistsExcludingUser(req.Email, userID)
		if err != nil {
			utils.Error("UpdateProfile: Failed to check email - %v", err)
			return utils.InternalServerErrorResponse(c, "Failed to update profile")
		}
		if exists {
			return utils.ConflictResponse(c, "Email already in use")
		}

		user.Email = req.Email
	}

	// Update full name if provided
	if req.FullName != "" {
		user.Name = req.FullName
	}

	// Update user in database
	if err := h.userRepo.Update(user); err != nil {
		utils.Error("UpdateProfile: Failed to update user - %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to update profile")
	}

	utils.Info("Profile updated: UserID=%d", userID)

	return utils.SuccessResponse(c, "Profile updated successfully", models.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.Name,
		RoleID:   user.RoleID,
		IsActive: user.IsActive,
	})
}

// ChangePassword changes user password
// PUT /api/auth/change-password
// Protected: Requires authentication
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.UnauthorizedResponse(c, "Unauthorized")
	}

	// Parse request body
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate input
	if req.OldPassword == "" || req.NewPassword == "" {
		return utils.BadRequestResponse(c, "Old password and new password are required")
	}

	// Get user
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return utils.NotFoundResponse(c, "User not found")
	}

	// Verify old password
	if !utils.VerifyPassword(req.OldPassword, user.PasswordHash) {
		return utils.BadRequestResponse(c, "Old password is incorrect")
	}

	// Validate new password strength
	if err := utils.ValidatePasswordStrength(req.NewPassword); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.Error("ChangePassword: Failed to hash password - %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to change password")
	}

	// Update password
	if err := h.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		utils.Error("ChangePassword: Failed to update password - %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to change password")
	}

	utils.Info("Password changed: UserID=%d", userID)

	return utils.SuccessMessage(c, "Password changed successfully")
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// isValidEmail validates email format
func isValidEmail(email string) bool {
	// Simple email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}