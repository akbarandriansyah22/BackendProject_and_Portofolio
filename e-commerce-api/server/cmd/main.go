package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/config"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/handler"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/middleware"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/repository"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/service"
)

func main() {
	cfg := config.MustLoad()

	zlog, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	loggerObs := observability.NewZapLoggerFrom(zlog)
	defer func() {
		if syncErr := loggerObs.Sync(); syncErr != nil {
			log.Printf("failed to sync logger: %v", syncErr)
		}
	}()
	loggerObs.Info("starting e-commerce API")

	observability.InitMetrics(cfg.App.Version)

	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		loggerObs.Fatal("failed to connect database: %v", err)
	}
	if err := db.Ping(); err != nil {
		loggerObs.Fatal("database unreachable: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			loggerObs.Error("failed to close db: %v", closeErr)
		}
	}()

	userRepo := repository.NewUserRepositoryPort(repository.NewUserRepository(db, zlog))
	roleRepo := repository.NewRoleRepositoryPort(repository.NewRoleRepository(db, zlog))
	productRepo := repository.NewProductRepository(db, zlog)
	categoryRepo := repository.NewCategoryRepositoryPort(repository.NewCategoryRepository(db, zlog))
	cartRepo := repository.NewCartRepositoryPort(repository.NewCartRepository(db, zlog))
	orderRepo := repository.NewOrderRepositoryPort(repository.NewOrderRepository(db, zlog))
	paymentRepo := repository.NewPaymentRepositoryPort(repository.NewPaymentRepository(db, zlog))

	authService := service.NewAuthService(userRepo, roleRepo, cfg.JWT.Secret, loggerObs)
	productService := service.NewProductService(productRepo, categoryRepo, loggerObs)
	categoryService := service.NewCategoryService(categoryRepo, productRepo, loggerObs)
	cartService := service.NewCartService(cartRepo, productRepo, loggerObs)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, paymentRepo, loggerObs)

	authHandler := handler.NewAuthHandler(authService, loggerObs)
	productHandler := handler.NewProductHandler(productService, loggerObs)
	categoryHandler := handler.NewCategoryHandler(categoryService, loggerObs)
	cartHandler := handler.NewCartHandler(cartService, loggerObs)
	orderHandler := handler.NewOrderHandler(orderService, loggerObs)

	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  60 * time.Second,
		BodyLimit:    2 * 1024 * 1024,
	})

	app.Use(
		recover.New(),
		logger.New(),
		middleware.CORS(),
	)

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		nextErr := c.Next()
		duration := time.Since(start).Seconds()

		route := c.Route().Path
		if route == "" {
			route = "unknown"
		}
		status := strconv.Itoa(c.Response().StatusCode())

		observability.HttpRequestsTotal.WithLabelValues(c.Method(), route, status).Inc()
		observability.HttpRequestDuration.WithLabelValues(c.Method(), route).Observe(duration)
		if c.Response().StatusCode() >= 400 {
			observability.HttpErrorsTotal.WithLabelValues(c.Method(), route, status).Inc()
		}
		return nextErr
	})

	metricsToken := os.Getenv("METRICS_TOKEN")
	observability.RegisterMetricsEndpoint(app, metricsToken)

	app.Get("/live", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/ready", func(c *fiber.Ctx) error {
		if pingErr := db.Ping(); pingErr != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unavailable",
				"db":     "unreachable",
			})
		}
		return c.JSON(fiber.Map{"status": "ready", "db": "ok"})
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "ok"
		httpStatus := fiber.StatusOK
		appStatus := "ok"
		if pingErr := db.Ping(); pingErr != nil {
			dbStatus = "unreachable"
			httpStatus = fiber.StatusServiceUnavailable
			appStatus = "degraded"
		}
		return c.Status(httpStatus).JSON(fiber.Map{
			"status":  appStatus,
			"db":      dbStatus,
			"version": cfg.App.Version,
		})
	})

	authRateLimiter := middleware.NewRateLimitMiddleware(
		middleware.DefaultAuthRateLimitConfig(loggerObs),
	)
	apiRateLimiter := middleware.NewRateLimitMiddleware(
		middleware.DefaultAPIRateLimitConfig(loggerObs),
	)

	registerAuth := func(g fiber.Router) {
		g.Post("/register", authHandler.Register)
		g.Post("/login", authHandler.Login)
	}
	registerAuth(app.Group("/auth", authRateLimiter))
	registerAuth(app.Group("/api/auth", authRateLimiter))

	publicAPI := app.Group("/api", apiRateLimiter)
	publicAPI.Get("/products", productHandler.ListProducts)
	publicAPI.Get("/products/search", productHandler.SearchProducts)
	publicAPI.Get("/products/slug/:slug", productHandler.GetProductBySlug)
	publicAPI.Get("/products/category/:id", productHandler.GetProductsByCategory)
	publicAPI.Get("/products/:id", productHandler.GetProductByID)
	publicAPI.Get("/categories", categoryHandler.ListCategories)
	publicAPI.Get("/categories/:id", categoryHandler.GetCategoryByID)
	publicAPI.Get("/categories/:id/products", categoryHandler.GetProductsByCategory)
	publicAPI.Get("/categories/:id/subcategories", categoryHandler.GetSubCategories)

	protected := app.Group("/api",
		middleware.Auth(cfg.JWT.Secret, loggerObs),
		apiRateLimiter,
	)
	protected.Get("/auth/profile", authHandler.GetProfile)
	protected.Put("/auth/profile", authHandler.UpdateProfile)
	protected.Put("/auth/change-password", authHandler.ChangePassword)

	protected.Get("/cart", cartHandler.GetCart)
	protected.Post("/cart/items", cartHandler.AddItem)
	protected.Delete("/cart/items/:id", cartHandler.RemoveItem)
	protected.Delete("/cart", cartHandler.ClearCart)

	protected.Get("/orders", orderHandler.ListOrders)
	protected.Post("/orders", orderHandler.CreateOrder)
	protected.Get("/orders/number/:orderNumber", orderHandler.GetByOrderNumber)
	protected.Get("/orders/:id", orderHandler.GetByID)
	protected.Post("/orders/:id/cancel", orderHandler.CancelOrder)

	admin := app.Group("/api/admin",
		middleware.Auth(cfg.JWT.Secret, loggerObs),
		middleware.RequireRole(loggerObs, middleware.RoleAdmin),
		apiRateLimiter,
	)
	admin.Post("/products", productHandler.CreateProduct)
	admin.Put("/products/:id", productHandler.UpdateProduct)
	admin.Delete("/products/:id", productHandler.DeleteProduct)
	admin.Put("/products/:id/activate", productHandler.ActivateProduct)
	admin.Put("/products/:id/deactivate", productHandler.DeactivateProduct)
	admin.Patch("/products/:id/stock", productHandler.UpdateStock)

	admin.Post("/categories", categoryHandler.CreateCategory)
	admin.Put("/categories/:id", categoryHandler.UpdateCategory)
	admin.Delete("/categories/:id", categoryHandler.DeleteCategory)
	admin.Patch("/categories/:id/status", categoryHandler.ToggleCategoryStatus)
	admin.Get("/categories/stats", categoryHandler.GetCategoryStats)

	admin.Get("/orders", orderHandler.GetAllOrders)
	admin.Put("/orders/:id/status", orderHandler.UpdateStatus)
	admin.Get("/orders/stats", orderHandler.GetOrderStats)

	if err := app.Listen(cfg.GetServerAddress()); err != nil {
		loggerObs.Fatal("failed to start server: %v", err)
	}
}
