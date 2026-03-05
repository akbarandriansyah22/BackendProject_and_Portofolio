package config

import "os"

// SecurityConfig holds all security-related configuration
type SecurityConfig struct {
	JWTSecret         string
	JWTExpiration     int
	PasswordMinLength int
	PasswordMaxLength int
	SessionTimeout    int
	MaxLoginAttempts  int
	LockoutDuration   int
	RateLimitPerMin   int
}

// NewSecurityConfig returns default security configuration
func NewSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTExpiration:     86400, // 24 hours
		PasswordMinLength: 8,
		PasswordMaxLength: 128,
		SessionTimeout:    900, // 15 minutes
		MaxLoginAttempts:  5,
		LockoutDuration:   900, // 15 minutes
		RateLimitPerMin:   60,
	}
}
