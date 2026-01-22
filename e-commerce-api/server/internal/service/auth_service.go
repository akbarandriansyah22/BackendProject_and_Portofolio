package service

import (
	"fmt"

	"github.com/akbarandriansyah22/Devops_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/Devops_Portofolio/e-commerce-api/server/internal/repository"
	"github.com/akbarandriansyah22/Devops_Portofolio/e-commerce-api/server/internal/utils"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  *repository.UserRepository
	roleRepo  *repository.RoleRepository
	jwtSecret string
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration business logic
func (s *AuthService) Register(req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Validate email
	if !isValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Validate password strength
	if err := utils.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Validate full name
	if req.FullName == "" {
		return nil, fmt.Errorf("full name is required")
	}

	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		utils.Error("AuthService.Register: Failed to check email existence - %v", err)
		return nil, fmt.Errorf("failed to register user")
	}
	if exists {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error("AuthService.Register: Failed to hash password - %v", err)
		return nil, fmt.Errorf("failed to register user")
	}

	// Get customer role (default role)
	// Assume roleID 2 = customer
	customerRoleID := 2

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:     req.FullName,
		RoleID:       customerRoleID,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		utils.Error("AuthService.Register: Failed to create user - %v", err)
		return nil, fmt.Errorf("failed to register user")
	}

	utils.Info("User registered successfully: UserID=%d, Email=%s", user.ID, user.Email)

	// Generate JWT token
	token, err := utils.GenerateAccessToken(user.ID, user.Email, user.RoleID, user.Name, s.jwtSecret)
	if err != nil {
		utils.Error("AuthService.Register: Failed to generate token - %v", err)
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
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		utils.Warn("AuthService.Login: Failed login attempt for email=%s", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		utils.Warn("AuthService.Login: Inactive user attempted login - UserID=%d", user.ID)
		return nil, fmt.Errorf("account is inactive")
	}

	// Verify password
	if !utils.VerifyPassword(req.Password, user.PasswordHash) {
		utils.Warn("AuthService.Login: Invalid password for email=%s", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateAccessToken(user.ID, user.Email, user.RoleID, user.Name, s.jwtSecret)
	if err != nil {
		utils.Error("AuthService.Login: Failed to generate token - %v", err)
		return nil, fmt.Errorf("failed to generate token")
	}

	utils.Info("User logged in successfully: UserID=%d, Email=%s", user.ID, user.Email)

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
func (s *AuthService) GetProfile(userID int) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		utils.Error("AuthService.GetProfile: User not found - UserID=%d", userID)
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
func (s *AuthService) UpdateProfile(userID int, fullName, email string) (*models.UserResponse, error) {
	// Get current user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Validate and update email if changed
	if email != "" && email != user.Email {
		if !isValidEmail(email) {
			return nil, fmt.Errorf("invalid email format")
		}

		// Check if new email already exists
		exists, err := s.userRepo.EmailExistsExcludingUser(email, userID)
		if err != nil {
			utils.Error("AuthService.UpdateProfile: Failed to check email - %v", err)
			return nil, fmt.Errorf("failed to update profile")
		}
		if exists {
			return nil, fmt.Errorf("email already in use")
		}

		user.Email = email
	}

	// Update full name if provided
	if fullName != "" {
		user.Name = fullName
	}

	// Update user in database
	if err := s.userRepo.Update(user); err != nil {
		utils.Error("AuthService.UpdateProfile: Failed to update user - %v", err)
		return nil, fmt.Errorf("failed to update profile")
	}

	utils.Info("Profile updated: UserID=%d", userID)

	return &models.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.Name,
		RoleID:   user.RoleID,
		IsActive: user.IsActive,
	}, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID int, oldPassword, newPassword string) error {
	// Validate input
	if oldPassword == "" || newPassword == "" {
		return fmt.Errorf("old password and new password are required")
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if !utils.VerifyPassword(oldPassword, user.PasswordHash) {
		return fmt.Errorf("old password is incorrect")
	}

	// Validate new password strength
	if err := utils.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		utils.Error("AuthService.ChangePassword: Failed to hash password - %v", err)
		return fmt.Errorf("failed to change password")
	}

	// Update password
	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		utils.Error("AuthService.ChangePassword: Failed to update password - %v", err)
		return fmt.Errorf("failed to change password")
	}

	utils.Info("Password changed: UserID=%d", userID)
	return nil
}

// ValidateToken validates JWT token and returns user ID
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	claims, err := utils.ParseToken(tokenString, s.jwtSecret)
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
	// Simple email validation
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