package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS  (Cross-Origin Resource Sharing)
// Middleware untuk handle CORS dengan konfigurasi default
// Allow all origins (untuk development)
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		// Allow all origins (wildcard)
		//  DEVELOPMENT ONLY! Ganti di production!
		AllowOrigins: "*",

		// Allow common HTTP methods
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",

		// Allow common headers
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Requested-With",

		// Expose headers to frontend
		ExposeHeaders: "Content-Length,Content-Type",

		// Allow credentials (cookies, authorization headers)
		AllowCredentials: false, // Set false jika AllowOrigins = "*"

		// Max age for preflight request cache (in seconds)
		MaxAge: 3600, // 1 hour
	})
}

// CORSWithOrigins returns a CORS middleware with specific allowed origins
// Middleware CORS dengan whitelist origins (untuk production)
func CORSWithOrigins(origins []string) fiber.Handler {
	// Convert slice to comma-separated string
	allowedOrigins := ""
	for i, origin := range origins {
		if i > 0 {
			allowedOrigins += ","
		}
		allowedOrigins += origin
	}

	return cors.New(cors.Config{
		// Allow specific origins only
		AllowOrigins: allowedOrigins,

		// Allow common HTTP methods
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",

		// Allow common headers + custom headers
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token",

		// Expose headers
		ExposeHeaders: "Content-Length,Content-Type,X-Total-Count",

		// Allow credentials (required for cookies & auth)
		AllowCredentials: true,

		// Max age for preflight cache
		MaxAge: 86400, // 24 hours
	})
}

// CORSProduction returns a CORS middleware for production
// Middleware CORS dengan konfigurasi production (strict)
func CORSProduction(allowedOrigins []string) fiber.Handler {
	// Convert slice to comma-separated string
	origins := ""
	for i, origin := range allowedOrigins {
		if i > 0 {
			origins += ","
		}
		origins += origin
	}

	return cors.New(cors.Config{
		// Strict: Only allow whitelisted origins
		AllowOrigins: origins,

		// Only allow necessary methods
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",

		// Only allow necessary headers
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",

		// Expose necessary headers
		ExposeHeaders: "Content-Length,Content-Type",

		// Allow credentials
		AllowCredentials: true,

		// Cache preflight for longer
		MaxAge: 86400, // 24 hours
	})
}

// CORSDevelopment returns a CORS middleware for development
// Middleware CORS untuk development (permissive)
func CORSDevelopment() fiber.Handler {
	return cors.New(cors.Config{
		// Allow all origins in development
		AllowOrigins: "*",

		// Allow all methods
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",

		// Allow all headers
		AllowHeaders: "*",

		// Expose all headers
		ExposeHeaders: "*",

		// No credentials with wildcard origin
		AllowCredentials: false,

		// Short cache for development (easier testing)
		MaxAge: 300, // 5 minutes
	})
}

// CORSCustom returns a CORS middleware with custom configuration
// Middleware CORS dengan konfigurasi custom (full control)
func CORSCustom(config cors.Config) fiber.Handler {
	return cors.New(config)
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// GetDefaultCORSConfig returns default CORS configuration
// Fungsi helper untuk get default CORS config (bisa di-customize)
func GetDefaultCORSConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		ExposeHeaders:    "Content-Length,Content-Type",
		AllowCredentials: false,
		MaxAge:           3600,
	}
}

// GetProductionCORSConfig returns production CORS configuration
// Fungsi helper untuk get production CORS config
func GetProductionCORSConfig(allowedOrigins []string) cors.Config {
	origins := ""
	for i, origin := range allowedOrigins {
		if i > 0 {
			origins += ","
		}
		origins += origin
	}

	return cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		ExposeHeaders:    "Content-Length,Content-Type",
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// ============================================
// ORIGIN VALIDATORS
// ============================================

// IsOriginAllowed checks if an origin is in the allowed list
// Fungsi helper untuk check apakah origin diizinkan
func IsOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

// GetAllowedOrigins returns list of allowed origins from environment
// Fungsi helper untuk get allowed origins dari config
// Contoh: origins := GetAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))
func GetAllowedOrigins(originsString string) []string {
	if originsString == "" {
		return []string{"*"}
	}

	// Parse comma-separated string
	origins := []string{}
	current := ""

	for _, char := range originsString {
		if char == ',' {
			if current != "" {
				origins = append(origins, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		origins = append(origins, current)
	}

	return origins
}
