package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/config"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
)

func main() {
	// =======================
	// CONFIG
	// =======================
	cfg := config.Load()

	// =======================
	// LOGGER
	// =======================
	loggerObs := observability.NewZapLogger()
	defer loggerObs.Sync()
	loggerObs.Info("starting e-commerce API")

	// =======================
	// PROMETHEUS METRICS
	// =======================
	observability.InitMetrics()

	// =======================
	// DATABASE
	// =======================
	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		loggerObs.Fatal("failed to connect database", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		loggerObs.Fatal("database unreachable", err)
	}

	// =======================
	// FIBER APP
	// =======================
	app := fiber.New(fiber.Config{
		AppName: "E-Commerce API",
	})

	// =======================
	// GLOBAL MIDDLEWARES
	// =======================
	app.Use(
		recover.New(),
		logger.New(),
		cors.New(),
	)

	// =======================
	// METRICS MIDDLEWARE (WAJIB SEBELUM ROUTES)
	// =======================
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()

		route := c.Route().Path
		if route == "" {
			route = "unknown"
		}

		status := strconv.Itoa(c.Response().StatusCode())

		observability.HttpRequestsTotal.WithLabelValues(
			c.Method(),
			route,
			status,
		).Inc()

		observability.HttpRequestDuration.WithLabelValues(
			c.Method(),
			route,
		).Observe(duration)

		if c.Response().StatusCode() >= 400 {
			observability.HttpErrorsTotal.WithLabelValues(
				c.Method(),
				route,
				status,
			).Inc()
		}

		return err
	})

	// =======================
	// METRICS ENDPOINT
	// =======================
	observability.RegisterMetricsEndpoint(app)

	// =======================
	// HEALTH ENDPOINT
	// =======================
	app.Get("/health", func(c *fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "error",
				"db":     "unreachable",
			})
		}

		return c.JSON(fiber.Map{
			"status":      "ok",
			"db":          "ok",
			"environment": cfg.Server.Environment,
			"timestamp":   time.Now().UTC(),
		})
	})

	// =======================
	// START SERVER
	// =======================
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}

