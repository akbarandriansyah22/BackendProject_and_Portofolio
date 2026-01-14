package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"e-commerce-api/server/internal/config"
	"e-commerce-api/server/internal/database"
	"e-commerce-api/server/internal/handler"
	"e-commerce-api/server/internal/middleware"
)

func main() {
	// Load .env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found")
	}

	// Load config
	cfg, err := config.LoadWithValidation()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	cfg.PrintConfig()

	// Initialize database
	initDatabase(cfg)
	defer database.Close()

	// Setup Fiber app
	app := setupApp(cfg)

	// Register routes
	registerRoutes(app, cfg)

	// Start server
	log.Printf("Server starting on %s", cfg.GetServerAddress())
	if err := app.Listen(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func initDatabase(cfg *config.Config) {
	log.Println("Initializing database...")

	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	_, err := database.InitDB(dbConfig, 5)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	database.LogConnectionStats()
	log.Println("Database initialized successfully!")
}

func setupApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		if err := database.HealthCheck(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "unhealthy",
				"error":  err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
		})
	})

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to E-Commerce API",
			"version": cfg.App.Version,
		})
	})

	return app
}

func registerRoutes(app *fiber.App, cfg *config.Config) {
	// API group
	api := app.Group("/api")

	// ============ AUTH ROUTES ============
	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// Protected auth routes
	authProtected := auth.Group("/")
	authProtected.Use(middleware.Auth(cfg.JWT.Secret))
	authProtected.Get("/profile", handler.GetProfile)
	authProtected.Put("/profile", handler.UpdateProfile)
	authProtected.Put("/change-password", handler.ChangePassword)

	// ============ PRODUCT ROUTES ============
	products := api.Group("/products")
	products.Get("/", handler.GetAllProducts)
	products.Get("/:id", handler.GetProductByID)
	products.Get("/search", handler.SearchProducts)
	products.Get("/featured", handler.GetFeaturedProducts)

	// Admin product routes
	adminProducts := api.Group("/admin/products")
	adminProducts.Use(middleware.Auth(cfg.JWT.Secret))
	adminProducts.Use(middleware.RequireRole(1))
	adminProducts.Post("/", handler.CreateProduct)
	adminProducts.Put("/:id", handler.UpdateProduct)
	adminProducts.Delete("/:id", handler.DeleteProduct)
	adminProducts.Patch("/:id/stock", handler.UpdateProductStock)
	adminProducts.Patch("/:id/status", handler.ToggleProductStatus)

	// ============ CATEGORY ROUTES ============
	categories := api.Group("/categories")
	categories.Get("/", handler.GetAllCategories)
	categories.Get("/:id", handler.GetCategoryByID)

	// Admin category routes
	adminCategories := api.Group("/admin/categories")
	adminCategories.Use(middleware.Auth(cfg.JWT.Secret))
	adminCategories.Use(middleware.RequireRole(1))
	adminCategories.Post("/", handler.CreateCategory)
	adminCategories.Put("/:id", handler.UpdateCategory)
	adminCategories.Delete("/:id", handler.DeleteCategory)
	adminCategories.Patch("/:id/status", handler.ToggleCategoryStatus)

	// ============ CART ROUTES ============
	cart := api.Group("/cart")
	cart.Use(middleware.Auth(cfg.JWT.Secret))
	cart.Get("/", handler.GetCart)
	cart.Post("/items", handler.AddItem)
	cart.Put("/items/:id", handler.UpdateItem)
	cart.Delete("/items/:id", handler.RemoveItem)
	cart.Delete("/", handler.ClearCart)

	// ============ ORDER ROUTES ============
	orders := api.Group("/orders")
	orders.Use(middleware.Auth(cfg.JWT.Secret))
	orders.Get("/", handler.GetAllOrders)
	orders.Get("/:id", handler.GetOrderByID)
	orders.Post("/checkout", handler.CreateFromCart)
	orders.Post("/:id/cancel", handler.CancelOrder)

	// Admin order routes
	adminOrders := api.Group("/admin/orders")
	adminOrders.Use(middleware.Auth(cfg.JWT.Secret))
	adminOrders.Use(middleware.RequireRole(1))
	adminOrders.Get("/", handler.GetAllAdminOrders)
	adminOrders.Put("/:id/status", handler.UpdateOrderStatus)

	log.Println("Routes registered successfully!")
}

