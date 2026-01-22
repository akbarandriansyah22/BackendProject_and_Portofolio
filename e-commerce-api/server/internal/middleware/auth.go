package middleware

import (
	"strings"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims stored in JWT token
// Claims = data yang disimpan di dalam token
type JWTClaims struct {
	UserID int    `json:"user_id"` // ID user
	Email  string `json:"email"`   // Email user
	RoleID int    `json:"role_id"` // Role ID (1=admin, 2=customer, dst)
	jwt.RegisteredClaims
}

// Auth middleware protects routes that require authentication
// Middleware ini digunakan untuk protect route yang butuh login
func Auth(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// STEP 1: Ambil Authorization header dari request
		// Header format: "Authorization: Bearer <token>"
		authHeader := c.Get("Authorization")
		
		// STEP 2: Cek apakah header ada
		if authHeader == "" {
			utils.Warn("Auth failed: No authorization header from IP %s", c.IP())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Authorization header required",
			})
		}

		// STEP 3: Parse header untuk extract token
		// Format header: "Bearer <token>"
		// Split untuk ambil token-nya saja
		parts := strings.Split(authHeader, " ")
		
		// Validasi format header
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Warn("Auth failed: Invalid authorization header format from IP %s", c.IP())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid authorization header format. Expected: Bearer <token>",
			})
		}

		// STEP 4: Ambil token string
		tokenString := parts[1]

		// STEP 5: Parse dan validasi JWT token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Cek signing method (harus HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				utils.Warn("Auth failed: Invalid signing method")
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			// Return secret key untuk validasi signature
			return []byte(jwtSecret), nil
		})

		// STEP 6: Handle error saat parsing token
		if err != nil {
			utils.Warn("Auth failed: Token parsing error - %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid or expired token",
			})
		}

		// STEP 7: Extract claims dari token
		claims, ok := token.Claims.(*JWTClaims)
		
		// STEP 8: Validasi claims dan token validity
		if !ok || !token.Valid {
			utils.Warn("Auth failed: Invalid token claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid token",
			})
		}

		// STEP 9: SUKSES! Simpan user info di context
		// Context = tempat simpan data yang bisa diakses di handler
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("roleID", claims.RoleID)
		c.Locals("user", claims) // Simpan semua claims sekaligus

		// Log successful authentication
		utils.Debug("Auth success: UserID=%d, Email=%s from IP %s", 
			claims.UserID, claims.Email, c.IP())

		// STEP 10: Lanjutkan ke handler berikutnya
		return c.Next()
	}
}

// RequireRole middleware checks if user has specific role
// Middleware ini cek apakah user punya role tertentu (misal: admin only)
func RequireRole(roleID int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// STEP 1: Ambil roleID dari context (sudah di-set oleh Auth middleware)
		userRoleID := c.Locals("roleID")
		
		// STEP 2: Cek apakah roleID ada (berarti user sudah login)
		if userRoleID == nil {
			utils.Warn("Access denied: No role ID in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Unauthorized",
			})
		}

		// STEP 3: Cek apakah role sesuai
		if userRoleID.(int) != roleID {
			userID := c.Locals("userID")
			utils.Warn("Access denied: UserID=%v attempted to access role %d endpoint", userID, roleID)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Forbidden: Insufficient permissions",
			})
		}

		// STEP 4: Role sesuai, lanjutkan
		return c.Next()
	}
}

// RequireAdmin middleware ensures only admin can access
// Shortcut untuk require role admin (roleID = 1)
func RequireAdmin() fiber.Handler {
	return RequireRole(1) // Assume roleID 1 = admin
}

// Optional middleware allows both authenticated and guest users
// Middleware ini tidak block jika tidak ada token (optional auth)
func Optional(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil header
		authHeader := c.Get("Authorization")
		
		// Jika tidak ada header, skip (allow guest)
		if authHeader == "" {
			return c.Next()
		}

		// Jika ada header, validate seperti biasa
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next() // Invalid format, tapi tetap allow (as guest)
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(jwtSecret), nil
		})

		// Jika parsing berhasil, simpan user info
		if err == nil {
			if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
				c.Locals("userID", claims.UserID)
				c.Locals("email", claims.Email)
				c.Locals("roleID", claims.RoleID)
				c.Locals("user", claims)
			}
		}

		// Lanjutkan ke handler (baik ada token atau tidak)
		return c.Next()
	}
}

// GetUserID retrieves user ID from context
// Helper function untuk ambil userID dari context di handler
func GetUserID(c *fiber.Ctx) (int, bool) {
	userID := c.Locals("userID")
	if userID == nil {
		return 0, false
	}
	id, ok := userID.(int)
	return id, ok
}

// GetUserEmail retrieves user email from context
// Helper function untuk ambil email dari context
func GetUserEmail(c *fiber.Ctx) (string, bool) {
	email := c.Locals("email")
	if email == nil {
		return "", false
	}
	emailStr, ok := email.(string)
	return emailStr, ok
}

// GetUserRoleID retrieves user role ID from context
// Helper function untuk ambil roleID dari context
func GetUserRoleID(c *fiber.Ctx) (int, bool) {
	roleID := c.Locals("roleID")
	if roleID == nil {
		return 0, false
	}
	id, ok := roleID.(int)
	return id, ok
}

// GetUserClaims retrieves full user claims from context
// Helper function untuk ambil semua claims sekaligus
func GetUserClaims(c *fiber.Ctx) (*JWTClaims, bool) {
	user := c.Locals("user")
	if user == nil {
		return nil, false
	}
	claims, ok := user.(*JWTClaims)
	return claims, ok
}

// IsAdmin checks if current user is admin
// Helper function untuk cek apakah user adalah admin
func IsAdmin(c *fiber.Ctx) bool {
	roleID, ok := GetUserRoleID(c)
	if !ok {
		return false
	}
	return roleID == 1 // Assume roleID 1 = admin
}