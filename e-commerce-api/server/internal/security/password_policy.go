package security

import (
	"crypto/subtle"
	"regexp"
	"unicode"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

// PasswordPolicy defines password security requirements
type PasswordPolicy struct {
	MinLength      int
	MaxLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireDigit   bool
	RequireSpecial bool
}

// DefaultPasswordPolicy returns standard password policy
func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:      8,
		MaxLength:      128,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: true,
	}
}

// ValidatePassword checks if password meets policy requirements
func (p *PasswordPolicy) ValidatePassword(password string) []string {
	var violations []string

	// Length check
	if len(password) < p.MinLength {
		violations = append(violations, "Password must be at least 8 characters")
	}
	if len(password) > p.MaxLength {
		violations = append(violations, "Password must not exceed 128 characters")
	}

	// Character type checks
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case !unicode.IsLetter(char) && !unicode.IsDigit(char):
			hasSpecial = true
		}
	}

	if p.RequireUpper && !hasUpper {
		violations = append(violations, "Password must contain at least one uppercase letter")
	}
	if p.RequireLower && !hasLower {
		violations = append(violations, "Password must contain at least one lowercase letter")
	}
	if p.RequireDigit && !hasDigit {
		violations = append(violations, "Password must contain at least one digit")
	}
	if p.RequireSpecial && !hasSpecial {
		violations = append(violations, "Password must contain at least one special character")
	}

	return violations
}

// IsCommonPassword checks against common password patterns
func IsCommonPassword(password string) bool {
	commonPatterns := []string{
		"^[0-9]{6,}$", // Only digits
		"^[a-z]{6,}$", // Only lowercase
		"^[A-Z]{6,}$", // Only uppercase
		"^qwerty",     // Keyboard patterns
		"^password",   // Common words
		"^admin",
		"^user",
		"^123456",
		"^password123",
	}

	for _, pattern := range commonPatterns {
		if regexp.MustCompile(pattern).MatchString(password) {
			return true
		}
	}

	return false
}

// ConstantTimeCompare safely compares two passwords
func ConstantTimeCompare(hashedPassword, plainPassword string) bool {
	return subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(plainPassword)) == 1
}

// GetPasswordPolicyFromRequest converts request to password policy
func GetPasswordPolicyFromRequest(req *models.CreateRoleRequest) *PasswordPolicy {
	return DefaultPasswordPolicy()
}
