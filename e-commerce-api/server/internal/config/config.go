package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	CORS     CORSConfig
	App      AppConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Environment  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// AppConfig holds general application configuration
type AppConfig struct {
	Name        string
	Version     string
	Environment string
	LogLevel    string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "ecommerce"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10) * time.Second,
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10) * time.Second,
			Environment:  getEnv("ENVIRONMENT", "development"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-this"),
			Expiration: getDurationEnv("JWT_EXPIRATION_HOURS", 24) * time.Hour,
		},
		CORS: CORSConfig{
			AllowedOrigins: getSliceEnv("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getSliceEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
			AllowedHeaders: getSliceEnv("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Authorization"}),
		},
		App: AppConfig{
			Name:        getEnv("APP_NAME", "E-Commerce API"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
	}
}

// LoadWithValidation loads and validates configuration
func LoadWithValidation() (*Config, error) {
	cfg := Load()

	// Validate required fields
	if cfg.Database.Password == "" {
		log.Println("⚠️  Warning: DB_PASSWORD is empty")
	}

	if cfg.JWT.Secret == "your-secret-key-change-this" {
		log.Println("⚠️  Warning: Using default JWT_SECRET. Please change in production!")
	}

	if cfg.Server.Environment == "production" {
		if cfg.Database.SSLMode == "disable" {
			log.Println("⚠️  Warning: SSL is disabled in production!")
		}
	}

	return cfg, nil
}

// IsDevelopment checks if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction checks if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// IsStaging checks if running in staging mode
func (c *Config) IsStaging() bool {
	return c.Server.Environment == "staging"
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

// PrintConfig prints the current configuration (without sensitive data)
func (c *Config) PrintConfig() {
	log.Println("📋 Application Configuration:")
	log.Printf("  App Name: %s", c.App.Name)
	log.Printf("  Version: %s", c.App.Version)
	log.Printf("  Environment: %s", c.Server.Environment)
	log.Printf("  Server: %s", c.GetServerAddress())
	log.Printf("  Database: %s@%s:%s/%s", c.Database.User, c.Database.Host, c.Database.Port, c.Database.DBName)
	log.Printf("  SSL Mode: %s", c.Database.SSLMode)
	log.Printf("  Log Level: %s", c.App.LogLevel)
	log.Printf("  JWT Expiration: %v", c.JWT.Expiration)
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets integer environment variable with default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv gets duration environment variable with default value
func getDurationEnv(key string, defaultValue int) time.Duration {
	return time.Duration(getIntEnv(key, defaultValue))
}



// getSliceEnv gets slice environment variable with default value
func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma
		var result []string
		current := ""
		for _, char := range value {
			if char == ',' {
				if current != "" {
					result = append(result, trimSpace(current))
					current = ""
				}
			} else {
				current += string(char)
			}
		}
		if current != "" {
			result = append(result, trimSpace(current))
		}
		return result
	}
	return defaultValue
}

// trimSpace removes leading and trailing spaces
func trimSpace(s string) string {
	start := 0
	end := len(s)

	// Trim leading spaces
	for start < end && s[start] == ' ' {
		start++
	}

	// Trim trailing spaces
	for end > start && s[end-1] == ' ' {
		end--
	}

	return s[start:end]
}