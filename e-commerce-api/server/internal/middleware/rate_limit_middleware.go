package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/security"
)

// RateLimitConfig menyimpan konfigurasi rate limiter middleware.
type RateLimitConfig struct {
	Capacity        int64
	Window          time.Duration
	CleanupInterval time.Duration
	KeyFunc         func(c *fiber.Ctx) string
	Logger          observability.Logger
}

// DefaultAuthRateLimitConfig — ketat untuk /auth/*
// 10 request per menit per IP, mencegah brute force login
func DefaultAuthRateLimitConfig(logger observability.Logger) RateLimitConfig {
	return RateLimitConfig{
		Capacity:        10,
		Window:          time.Minute,
		CleanupInterval: 5 * time.Minute,
		KeyFunc:         IPKeyFunc,
		Logger:          logger,
	}
}

// DefaultAPIRateLimitConfig — standar untuk /api/*
// 60 request per menit per userID
func DefaultAPIRateLimitConfig(logger observability.Logger) RateLimitConfig {
	return RateLimitConfig{
		Capacity:        60,
		Window:          time.Minute,
		CleanupInterval: 5 * time.Minute,
		KeyFunc:         UserIDKeyFunc,
		Logger:          logger,
	}
}

// IPKeyFunc — gunakan IP sebagai key.
// Cocok untuk endpoint publik sebelum user login.
func IPKeyFunc(c *fiber.Ctx) string {
	return "ip:" + c.IP()
}

// UserIDKeyFunc — gunakan userID dari JWT sebagai key.
// Lebih adil karena tidak menghukum semua user di balik NAT yang sama.
func UserIDKeyFunc(c *fiber.Ctx) string {
	if userID, ok := c.Locals("userID").(int); ok {
		return fmt.Sprintf("uid:%d", userID)
	}
	// Fallback ke IP jika userID tidak tersedia
	return "ip:" + c.IP()
}

// NewRateLimitMiddleware membuat Fiber middleware dari RateLimitConfig.
//
// Cara pakai di main.go:
//
//	authRL := middleware.NewRateLimitMiddleware(middleware.DefaultAuthRateLimitConfig(logger))
//	app.Group("/auth").Use(authRL)
//
//	apiRL := middleware.NewRateLimitMiddleware(middleware.DefaultAPIRateLimitConfig(logger))
//	app.Group("/api").Use(middleware.Auth(cfg.JWT.Secret, logger), apiRL)
func NewRateLimitMiddleware(cfg RateLimitConfig) fiber.Handler {
	if cfg.CleanupInterval == 0 {
		cfg.CleanupInterval = 5 * time.Minute
	}
	if cfg.KeyFunc == nil {
		cfg.KeyFunc = IPKeyFunc
	}

	// Buat rate limiter dengan cleanup goroutine aktif
	limiter := security.NewRateLimiterWithCleanup(
		cfg.Capacity,
		cfg.Window,
		cfg.CleanupInterval,
	)

	return func(c *fiber.Ctx) error {
		key := cfg.KeyFunc(c)

		if !limiter.IsAllowed(key) {
			remaining := limiter.GetRemaining(key)

			// Log violation
			if cfg.Logger != nil {
				cfg.Logger.Warn(
					"rate_limit_exceeded",
					"key=%s ip=%s path=%s method=%s remaining=%d",
					key,
					c.IP(),
					c.Path(),
					c.Method(),
					remaining,
				)
			}

			// Set standard rate limit headers (RFC 6585)
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Capacity))
			c.Set("X-RateLimit-Remaining", "0")

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many requests. Please slow down and try again later.",
			})
		}

		// Set informational headers untuk request yang diizinkan
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Capacity))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.GetRemaining(key)))

		return c.Next()
	}
}