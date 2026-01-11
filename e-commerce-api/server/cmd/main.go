package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"e-commerce-api/server/internal/config"
	"e-commerce-api/server/internal/database"
)

func main() {
	// Load .env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("⚠️  No .env file found")
	}

	// Load config
	cfg, err := config.LoadWithValidation()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}
	cfg.PrintConfig()

	// Initialize database
	db := initDatabase(cfg)
	defer database.Close()
	
	// Untuk sementara, gunakan _ untuk avoid unused variable
	_ = db
	
	log.Println("✅ Application initialized successfully!")
	
	// Setup Gin router
	r := setupRouter(cfg)

	// Start server
	log.Printf("🚀 Server starting on %s", cfg.GetServerAddress())
	if err := r.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}

func initDatabase(cfg *config.Config) *sql.DB {
	log.Println("\n🔍 Initializing database...")

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
		log.Fatalf("❌ Database initialization failed: %v", err)
	}

	if err := database.WaitForDB(db, 30*time.Second); err != nil {
		log.Fatalf("❌ Database not ready: %v", err)
	}

	database.LogConnectionStats()
	log.Println("✅ Database initialized successfully!")

	return db
}

func setupRouter(cfg *config.Config) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		if err := database.HealthCheck(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"status":  "healthy",
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
		})
	})

	// Simple test endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to E-Commerce API",
			"version": cfg.App.Version,
		})
	})
	
	return r
}