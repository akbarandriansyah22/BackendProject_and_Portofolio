package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"e-commerce-api/server/internal/config"
	"e-commerce-api/server/internal/database"
	"e-commerce-api/server/internal/handler"
	"e-commerce-api/server/internal/middleware"
	"e-commerce-api/server/internal/repository"
	"e-commerce-api/server/internal/service"
	"e-commerce-api/server/internal/utils"
)

func main() {
	// ============================================
	// 1. INITIALIZE LOGGER
	// ============================================
	if err := utils.InitLogger("logs", true); err != nil {
		log.Fatalf(" Failed to initialize logger: %v", err)
	}
	defer utils.CloseLogger()
	
	utils.Info("Starting E-Commerce API Server...")

	// ============================================
	// 2. LOAD ENVIRONMENT VARIABLES
	// ============================================
	envPaths := []string{".env", "../.env", "../../.env"}
	envLoaded := false
	
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			utils.Info(" Loaded .env from: %s", path)
			envLoaded = true
			break
		}
	}
	
	if !envLoaded {
		utils.Warn("⚠️  No .env file found, using environment variables")
	}

	// ============================================
	// 3. LOAD CONFIGURATION
	// ============================================
	cfg, err := config.LoadWithValidation()
	if err != nil {
		utils.Fatal(" Config error: %v", err)
	}
	cfg.PrintConfig()

	// ============================================
	// 4. INITIALIZE DATABASE
	// ============================================
	db := initDatabase(cfg)
	defer func() {
		utils.Info(" Closing database connection...")
		if err := database.Close(); err != nil {
			utils.Error(" Error closing database: %v", err)
		}
	}()

	// ============================================
	// 5. INITIALIZE ALL LAYERS
	// ============================================
	utils.Info(" Initializing repositories...")
	
	// REPOSITORIES
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	utils.Info(" Initializing services...")
	
	// SERVICES
	authService := service.NewAuthService(userRepo, roleRepo, cfg.JWT.Secret)
	productService := service.NewProductService(productRepo, categoryRepo)
    categoryService := service.NewCategoryService(categoryRepo,productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, paymentRepo)

	utils.Info(" Initializing handlers...")
	
	// HANDLERS
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)

	// ============================================
	// 6. SETUP FIBER APP
	// ============================================
	utils.Info("Setting up Fiber app...")
	app := setupApp(cfg)

	// ============================================
	// 7. REGISTER ROUTES
	// ============================================
	utils.Info("Registering routes...")
	registerRoutes(app, authHandler, productHandler, categoryHandler, cartHandler, orderHandler, cfg)

	// ============================================
	// 8. START SERVER WITH GRACEFUL SHUTDOWN
	// ============================================
	startServerWithGracefulShutdown(app, cfg)
}

// ============================================
// INITIALIZATION FUNCTIONS
// ============================================

func initDatabase(cfg *config.Config) *sql.DB {
	utils.Info(" Initializing database...")

	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := database.InitDB(dbConfig, 5)
	if err != nil {
		utils.Fatal("Database initialization failed: %v", err)
	}

	database.LogConnectionStats()
	utils.Info("Database initialized successfully")

	return db
}

func setupApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// CORS middleware
	if cfg.IsDevelopment() {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
			AllowCredentials: false,
		}))
	} else {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.CORS.AllowedOrigins[0],
			AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
			AllowCredentials: true,
		}))
	}

	// Health check
	app.Get("/health", healthCheckHandler)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to E-Commerce API",
			"version": cfg.App.Version,
			"status":  "running",
		})
	})

	return app
}

func registerRoutes(
	app *fiber.App, 
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	categoryHandler *handler.CategoryHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	cfg *config.Config,
) {
	// API group
	api := app.Group("/api")

	// ============================================
	// AUTH ROUTES
	// ============================================
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected auth routes
	authProtected := auth.Group("/")
	authProtected.Use(middleware.Auth(cfg.JWT.Secret))
	authProtected.Get("/profile", authHandler.GetProfile)
	authProtected.Put("/profile", authHandler.UpdateProfile)
	authProtected.Put("/change-password", authHandler.ChangePassword)

	// ============================================
	// PRODUCT ROUTES (Public)
	// ============================================
	products := api.Group("/products")
	products.Get("/", productHandler.GetAllProducts)
	products.Get("/:id", productHandler.GetProductByID)
	products.Get("/slug/:slug", productHandler.GetProductBySlug)
	products.Get("/search", productHandler.SearchProducts)
	products.Get("/featured", productHandler.GetFeaturedProducts)
	products.Get("/category/:category_id", productHandler.GetProductsByCategory)

	// ============================================
	// PRODUCT ROUTES (Admin Only)
	// ============================================
	adminProducts := api.Group("/admin/products")
	adminProducts.Use(middleware.Auth(cfg.JWT.Secret))
	adminProducts.Use(middleware.RequireRole(1))
	adminProducts.Post("/", productHandler.CreateProduct)
	adminProducts.Put("/:id", productHandler.UpdateProduct)
	adminProducts.Delete("/:id", productHandler.DeleteProduct)
	adminProducts.Patch("/:id/stock", productHandler.UpdateProductStock)
	adminProducts.Patch("/:id/status", productHandler.ToggleProductStatus)
	adminProducts.Get("/low-stock", productHandler.GetLowStockProducts)
	adminProducts.Get("/stats", productHandler.GetProductStats)

	// ============================================
	// CATEGORY ROUTES (Public)
	// ============================================
	categories := api.Group("/categories")
	categories.Get("/", categoryHandler.GetAllCategories)
	categories.Get("/:id", categoryHandler.GetCategoryByID)
	categories.Get("/:id/products", categoryHandler.GetProductsByCategory)

	// ============================================
	// CATEGORY ROUTES (Admin Only)
	// ============================================
	adminCategories := api.Group("/admin/categories")
	adminCategories.Use(middleware.Auth(cfg.JWT.Secret))
	adminCategories.Use(middleware.RequireRole(1))
	adminCategories.Post("/", categoryHandler.CreateCategory)
	adminCategories.Put("/:id", categoryHandler.UpdateCategory)
	adminCategories.Delete("/:id", categoryHandler.DeleteCategory)
	adminCategories.Patch("/:id/status", categoryHandler.ToggleCategoryStatus)

	// ============================================
	// CART ROUTES (Protected)
	// ============================================
	cart := api.Group("/cart")
	cart.Use(middleware.Auth(cfg.JWT.Secret))
	cart.Get("/", cartHandler.GetCart)
	cart.Post("/items", cartHandler.AddItem)
	cart.Put("/items/:id", cartHandler.UpdateItem)
	cart.Delete("/items/:id", cartHandler.RemoveItem)
	cart.Delete("/", cartHandler.ClearCart)
	cart.Post("/sync-prices", cartHandler.SyncPrices)
	cart.Get("/validate", cartHandler.ValidateForCheckout)

	// ============================================
	// ORDER ROUTES (Protected)
	// ============================================
	orders := api.Group("/orders")
	orders.Use(middleware.Auth(cfg.JWT.Secret))
	orders.Get("/", orderHandler.GetAll)
	orders.Get("/:id", orderHandler.GetByID)
	orders.Get("/number/:orderNumber", orderHandler.GetByOrderNumber)
	orders.Post("/checkout", orderHandler.CreateFromCart)
	orders.Post("/:id/cancel", orderHandler.CancelOrder)

	// ============================================
	// ORDER ROUTES (Admin Only)
	// ============================================
	adminOrders := api.Group("/admin/orders")
	adminOrders.Use(middleware.Auth(cfg.JWT.Secret))
	adminOrders.Use(middleware.RequireRole(1))
	adminOrders.Get("/", orderHandler.GetAllOrders)
	adminOrders.Put("/:id/status", orderHandler.UpdateStatus)
	adminOrders.Get("/stats", orderHandler.GetOrderStats)

	utils.Info("Routes registered successfully")
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func healthCheckHandler(c *fiber.Ctx) error {
	if err := database.HealthCheck(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "healthy",
		"app":     "E-Commerce API",
		"version": "2.0.0",
		"db":      "connected",
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	utils.Error("Error occurred: %v", err)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}

func startServerWithGracefulShutdown(app *fiber.App, cfg *config.Config) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		utils.Info(" Server starting on %s", cfg.GetServerAddress())
		if err := app.Listen(cfg.GetServerAddress()); err != nil {
			utils.Fatal(" Server error: %v", err)
		}
	}()

	<-quit
	utils.Info(" Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		utils.Error("Server shutdown error: %v", err)
	}

	if err := database.Close(); err != nil {
		utils.Error("Error closing database: %v", err)
	}

	utils.CloseLogger()
	utils.Info(" Server stopped gracefully")
}