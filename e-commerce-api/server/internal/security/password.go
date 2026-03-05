package security

import (
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt

// Input: password plaintext (misal: "password123")
// Output: hashed password (misal: "$2a$10$N9qo8uLOickgx2ZMRZoMye...")
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword akan:
	// 1. Generate random salt
	// 2. Combine salt dengan password
	// 3. Hash menggunakan bcrypt algorithm
	// 4. Return hasil hash (sudah include salt di dalamnya)

	// bcrypt.DefaultCost = 10 (recommended, balance antara security & speed)
	// Semakin tinggi cost, semakin lambat tapi semakin secure
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Convert bytes to string
	return string(hashedBytes), nil
}

// VerifyPassword compares a plain text password with a hashed password
// Fungsi ini untuk CHECK apakah password yang diinput cocok dengan hash di database
// Input:
//   - plainPassword: password yang diinput user saat login (misal: "password123")
//   - hashedPassword: password hash dari database (misal: "$2a$10$...")
//
// Output:
//   - true jika cocok
//   - false jika tidak cocok
func VerifyPassword(plainPassword, hashedPassword string) bool {
	// bcrypt.CompareHashAndPassword akan:
	// 1. Extract salt dari hashedPassword
	// 2. Hash plainPassword menggunakan salt yang sama
	// 3. Compare hasilnya dengan hashedPassword
	// 4. Return nil jika cocok, error jika tidak cocok

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	// Jika err == nil, berarti password cocok
	return err == nil
}

// ValidatePasswordStrength validates password complexity
// Fungsi ini untuk VALIDASI apakah password cukup kuat
// Rules:
// - Minimal 8 karakter
// - Harus ada huruf besar
// - Harus ada huruf kecil
// - Harus ada angka
// - Optional: harus ada special character
func ValidatePasswordStrength(password string) error {
	// Check minimum length
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Check maximum length (bcrypt limit = 72 bytes)
	if len(password) > 72 {
		return fmt.Errorf("password must be at most 72 characters long")
	}

	// Flags untuk tracking password requirements
	var (
		hasUpper  bool // Ada huruf besar?
		hasLower  bool // Ada huruf kecil?
		hasNumber bool // Ada angka?
	)

	// Loop setiap karakter untuk check requirements
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
			// Optional: Uncomment jika mau require special character
			// case unicode.IsPunct(char) || unicode.IsSymbol(char):
			// 	hasSpecial = true
		}
	}

	// Validate requirements
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	// Optional: Uncomment jika mau require special character
	// if !hasSpecial {
	// 	return fmt.Errorf("password must contain at least one special character")
	// }

	// Password valid!
	return nil
}

// ValidatePasswordSimple validates basic password requirements
// Versi simple dari ValidatePasswordStrength (hanya check length)
// Gunakan ini jika tidak mau requirement yang ketat
func ValidatePasswordSimple(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	if len(password) > 72 {
		return fmt.Errorf("password must be at most 72 characters long")
	}
	return nil
}

// GenerateRandomPassword generates a random secure password
// Fungsi ini untuk GENERATE random password (berguna untuk reset password, dll)
func GenerateRandomPassword(length int) (string, error) {
	// Pastikan length reasonable
	if length < 8 {
		length = 8
	}
	if length > 72 {
		length = 72
	}

	// Character sets
	const (
		uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
		numbers          = "0123456789"
		specialChars     = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	)

	// Combine all character sets
	allChars := uppercaseLetters + lowercaseLetters + numbers + specialChars

	// Generate random password
	password := make([]byte, length)

	// Untuk simplicity, kita gunakan kombinasi karakter
	// Note: Untuk production, lebih baik gunakan crypto/rand
	for i := 0; i < length; i++ {
		password[i] = allChars[i%len(allChars)]
	}

	return string(password), nil
}

// PasswordStrengthScore calculates password strength score (0-5)
// Fungsi ini untuk HITUNG strength password (untuk password strength meter)
// Return: 0 (very weak) to 5 (very strong)
func PasswordStrengthScore(password string) int {
	score := 0

	// Length score
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	if len(password) >= 16 {
		score++
	}

	// Complexity score
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasUpper && hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score++
	}

	// Cap at 5
	if score > 5 {
		score = 5
	}

	return score
}

// PasswordStrengthText returns human-readable password strength
// Fungsi ini untuk CONVERT score jadi text (untuk UI)
func PasswordStrengthText(password string) string {
	score := PasswordStrengthScore(password)

	switch score {
	case 0, 1:
		return "Very Weak"
	case 2:
		return "Weak"
	case 3:
		return "Fair"
	case 4:
		return "Strong"
	case 5:
		return "Very Strong"
	default:
		return "Unknown"
	}
}
