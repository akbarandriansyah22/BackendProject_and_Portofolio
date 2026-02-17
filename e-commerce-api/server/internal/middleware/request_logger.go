package middleware

import (
	"time"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/gofiber/fiber/v2"
)

func RequestLogger(logger observability.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		logger.Info(
			"http_request",
			"method=%s path=%s status=%d latency=%s ip=%s",
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			latency,
			c.IP(),
		)

		return err
	}
}
