package service

import (
	"context"
	"fmt"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/ports"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/security"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  ports.UserRepository
	roleRepo  ports.RoleRepository
	jwtSecret string
	logger    observability.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo ports.UserRepository,
	roleRepo ports.RoleRepository,
	jwtSecret string,
	logger observability.Logger,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// Register handles user registration business logic
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Validate email
	if !isValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Validate password strength
	if err := security.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Validate full name
	if req.FullName == "" {
		return nil, fmt.Errorf("full name is required")
	}

	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("AuthService.Register: Failed to hash password", err)
		return nil, fmt.Errorf("failed to register user")
	}

	// Get customer role (default role ID = 2)
	customerRoleID := 2

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.FullName,
		RoleID:       customerRoleID,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("AuthService.Register: Failed to create user", err)
		return nil, fmt.Errorf("failed to register user")
	}

	s.logger.Info("User registered successfully: UserID=%d, Email=%s", user.ID, user.Email)

	// Generate JWT token
	token, err := security.GenerateToken(user.ID, user.Email, user.RoleID, user.Name, s.jwtSecret, 24)
	if err != nil {
		s.logger.Error("AuthService.Register: Failed to generate token", err)
		return nil, fmt.Errorf("failed to generate token")
	}

	// Return login response
	return &models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.Name,
			RoleID:   user.RoleID,
			IsActive: user.IsActive,
		},
	}, nil
}

// Login handles user login business logic
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		s.logger.Warn("AuthService.Login: Failed login attempt for email=%s", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		s.logger.Warn("AuthService.Login: Inactive user attempted login - UserID=%d", user.ID)
		return nil, fmt.Errorf("account is inactive")
	}

	// Verify password
	if !security.VerifyPassword(req.Password, user.PasswordHash) {
		s.logger.Warn("AuthService.Login: Invalid password for email=%s", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := security.GenerateToken(user.ID, user.Email, user.RoleID, user.Name, s.jwtSecret, 24)
	if err != nil {
		s.logger.Error("AuthService.Login: Failed to generate token", err)
		return nil, fmt.Errorf("failed to generate token")
	}

	s.logger.Info("User logged in successfully: UserID=%d, Email=%s", user.ID, user.Email)

	// Return login response
	return &models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.Name,
			RoleID:   user.RoleID,
			IsActive: user.IsActive,
		},
	}, nil
}

// GetProfile gets user profile by ID
func (s *AuthService) GetProfile(ctx context.Context, userID int) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		s.logger.Error("AuthService.GetProfile: User not found - UserID=%d", err)
		return nil, fmt.Errorf("user not found")
	}

	return &models.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.Name,
		RoleID:   user.RoleID,
		IsActive: user.IsActive,
	}, nil
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(ctx context.Context, userID int, fullName, email string) (*models.UserResponse, error) {
	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Validate and update email if changed
	if email != "" && email != user.Email {
		if !isValidEmail(email) {
			return nil, fmt.Errorf("invalid email format")
		}

		// Check if new email already exists
		existingUser, _ := s.userRepo.GetByEmail(ctx, email)
		if existingUser != nil {
			return nil, fmt.Errorf("email already in use")
		}

		user.Email = email
	}

	// Update full name if provided
	if fullName != "" {
		user.Name = fullName
	}

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("AuthService.UpdateProfile: Failed to update user", err)
		return nil, fmt.Errorf("failed to update profile")
	}

	s.logger.Info("Profile updated: UserID=%d", userID)

	return &models.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.Name,
		RoleID:   user.RoleID,
		IsActive: user.IsActive,
	}, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	// Validate input
	if oldPassword == "" || newPassword == "" {
		return fmt.Errorf("old password and new password are required")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if !security.VerifyPassword(oldPassword, user.PasswordHash) {
		return fmt.Errorf("old password is incorrect")
	}

	// Validate new password strength
	if err := security.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := security.HashPassword(newPassword)
	if err != nil {
		s.logger.Error("AuthService.ChangePassword: Failed to hash password", err)
		return fmt.Errorf("failed to change password")
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		s.logger.Error("AuthService.ChangePassword: Failed to update password", err)
		return fmt.Errorf("failed to change password")
	}

	s.logger.Info("Password changed: UserID=%d", userID)
	return nil
}

// ValidateToken validates JWT token and returns user ID
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	claims, err := security.ParseToken(tokenString, s.jwtSecret)
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// isValidEmail validates email format
func isValidEmail(email string) bool {
	if len(email) < 3 {
		return false
	}

	atIndex := -1
	dotIndex := -1

	for i, char := range email {
		if char == '@' {
			atIndex = i
		}
		if char == '.' && i > atIndex {
			dotIndex = i
		}
	}

	return atIndex > 0 && dotIndex > atIndex+1 && dotIndex < len(email)-1
}
