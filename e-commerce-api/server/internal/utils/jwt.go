package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ============================================
// JWT CLAIMS STRUCTURE
// ============================================

// JWTClaims represents the claims stored in JWT token
// Claims = data yang disimpan di dalam token
type JWTClaims struct {
	UserID   int    `json:"user_id"`   // ID user dari database
	Email    string `json:"email"`     // Email user
	RoleID   int    `json:"role_id"`   // Role ID (1=admin, 2=customer, dll)
	FullName string `json:"full_name"` // Nama lengkap user
	jwt.RegisteredClaims
}

// ============================================
// TOKEN GENERATION
// ============================================

// GenerateToken generates a new JWT token
// Fungsi untuk GENERATE JWT token (dipanggil saat login sukses)
// Input: userID, email, roleID, fullName, secret, expiresIn (dalam jam)
// Output: token string
func GenerateToken(userID int, email string, roleID int, fullName string, secret string, expiresIn int) (string, error) {
	// Set expiration time
	// expiresIn dalam jam, convert ke time.Duration
	expirationTime := time.Now().Add(time.Duration(expiresIn) * time.Hour)

	// Create claims
	claims := &JWTClaims{
		UserID:   userID,
		Email:    email,
		RoleID:   roleID,
		FullName: fullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Token expired time
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // Token issued time
			NotBefore: jwt.NewNumericDate(time.Now()),     // Token valid after
			Issuer:    "e-commerce-api",                   // Issuer name
			Subject:   email,                               // Subject (email)
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateAccessToken generates an access token (short-lived)
// Fungsi shortcut untuk generate access token (expired dalam 24 jam)
// Contoh: token, _ := GenerateAccessToken(user, secret)
func GenerateAccessToken(userID int, email string, roleID int, fullName string, secret string) (string, error) {
	return GenerateToken(userID, email, roleID, fullName, secret, 24) // 24 hours
}

// GenerateRefreshToken generates a refresh token (long-lived)
// Fungsi shortcut untuk generate refresh token (expired dalam 7 hari)
// Contoh: refreshToken, _ := GenerateRefreshToken(user, secret)
func GenerateRefreshToken(userID int, email string, roleID int, fullName string, secret string) (string, error) {
	return GenerateToken(userID, email, roleID, fullName, secret, 24*7) // 7 days
}

// ============================================
// TOKEN VERIFICATION
// ============================================

// ParseToken parses and validates a JWT token
// Fungsi untuk PARSE dan VALIDATE JWT token
// Input: token string, secret
// Output: claims (jika valid), error (jika invalid)
func ParseToken(tokenString string, secret string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
	// Validate signing method 
	if token.Method != jwt.SigningMethodHS256 {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(secret), nil
})


	// Check for parsing errors
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// VerifyToken verifies if a token is valid
// Fungsi untuk VERIFY apakah token valid (tanpa extract claims)
// Return: true jika valid, false jika invalid
func VerifyToken(tokenString string, secret string) bool {
	_, err := ParseToken(tokenString, secret)
	return err == nil
}

// ============================================
// TOKEN EXTRACTION
// ============================================

// ExtractUserID extracts user ID from token
// Fungsi helper untuk extract userID dari token
func ExtractUserID(tokenString string, secret string) (int, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ExtractEmail extracts email from token
// Fungsi helper untuk extract email dari token
func ExtractEmail(tokenString string, secret string) (string, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}

// ExtractRoleID extracts role ID from token
// Fungsi helper untuk extract roleID dari token
func ExtractRoleID(tokenString string, secret string) (int, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return 0, err
	}
	return claims.RoleID, nil
}

// ============================================
// TOKEN VALIDATION CHECKS
// ============================================

// IsTokenExpired checks if token is expired
// Fungsi untuk cek apakah token sudah expired
func IsTokenExpired(tokenString string, secret string) bool {
	claims, err := ParseToken(tokenString, secret)
	if err != nil || claims.ExpiresAt == nil { 
		return true // If can't parse, consider expired
	}

	// Check if token is expired
	return claims.ExpiresAt.Time.Before(time.Now())
}

// GetTokenExpiration gets token expiration time
// Fungsi untuk get waktu expired token
func GetTokenExpiration(tokenString string, secret string) (time.Time, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return time.Time{}, err
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, fmt.Errorf("token has no expiration")
	}

	return claims.ExpiresAt.Time, nil
}

// GetTokenRemainingTime gets remaining valid time of token
// Fungsi untuk get sisa waktu valid token
func GetTokenRemainingTime(tokenString string, secret string) (time.Duration, error) {
	expiration, err := GetTokenExpiration(tokenString, secret)
	if err != nil {
		return 0, err
	}

	remaining := time.Until(expiration)
	if remaining < 0 {
		return 0, fmt.Errorf("token expired")
	}

	return remaining, nil
}

// ============================================
// TOKEN REFRESH
// ============================================

// RefreshToken refreshes an existing token
// Fungsi untuk REFRESH token (generate token baru dari token lama)
// Use case: Token hampir expired, generate token baru
func RefreshToken(oldTokenString string, secret string, expiresIn int) (string, error) {
	// Parse old token (even if expired)
	claims, err := ParseToken(oldTokenString, secret)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Generate new token with same claims but new expiration
	return GenerateToken(claims.UserID, claims.Email, claims.RoleID, claims.FullName, secret, expiresIn)
}

// ============================================
// TOKEN BLACKLIST HELPERS
// ============================================

// GetTokenIdentifier gets a unique identifier for the token
// Fungsi untuk get unique ID token (untuk blacklist)
// Return: hash dari token untuk disimpan di blacklist
func GetTokenIdentifier(tokenString string) string {
	// Simple implementation: return first 32 chars
	// In production: use hash (SHA256)
	if len(tokenString) > 32 {
		return tokenString[:32]
	}
	return tokenString
}

// ============================================
// VALIDATION HELPERS
// ============================================

// ValidateTokenFormat validates basic token format
// Fungsi untuk validate format token (tanpa verify signature)
func ValidateTokenFormat(tokenString string) bool {
	// Basic check: JWT should have 3 parts separated by dots
	// Format: header.payload.signature
	parts := 0
	for _, char := range tokenString {
		if char == '.' {
			parts++
		}
	}
	return parts == 2 && len(tokenString) > 20
}

// ============================================
// TOKEN INFO
// ============================================

// GetTokenInfo gets detailed information about token
// Fungsi untuk get informasi lengkap tentang token
func GetTokenInfo(tokenString string, secret string) (map[string]interface{}, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	var (
		issuedAt  time.Time
		expiresAt time.Time
		remaining time.Duration
		isExpired bool
	)

	if claims.IssuedAt != nil {
		issuedAt = claims.IssuedAt.Time
	}

	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
		isExpired = expiresAt.Before(time.Now())

		remaining = time.Until(expiresAt)
		if remaining < 0 {
			remaining = 0
		}
	} else {
		isExpired = true
	}

	info := map[string]interface{}{
		"user_id":    claims.UserID,
		"email":      claims.Email,
		"role_id":    claims.RoleID,
		"full_name":  claims.FullName,
		"issued_at":  issuedAt,
		"expires_at": expiresAt,
		"issuer":     claims.Issuer,
		"is_expired": isExpired,
		"remaining":  remaining.String(),
	}

	return info, nil
}


// ============================================
// CONSTANTS
// ============================================

const (
	// Token expiration times (in hours)
	AccessTokenExpiry  = 24      // 24 hours
	RefreshTokenExpiry = 24 * 7  // 7 days
	
	// Token types
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)
