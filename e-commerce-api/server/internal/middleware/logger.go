package middleware

import (
	"time"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Logger middleware for Fiber
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		statusCode := c.Response().StatusCode()

		// Get client IP
		clientIP := c.IP()

		// Log request
		utils.LogRequest(c.Method(), c.Path(), statusCode, duration, clientIP)

		return err
	}
}
