package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DB holds the database connection pool
var DB *sql.DB

// Connect establishes a connection to the database
func Connect(config Config) (*sql.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)                  // Maximum number of open connections
	db.SetMaxIdleConns(5)                   // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute)  // Maximum lifetime of a connection
	db.SetConnMaxIdleTime(10 * time.Minute) // Maximum idle time of a connection

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully!")
	log.Printf("📊 Database: %s@%s:%s/%s", config.User, config.Host, config.Port, config.DBName)

	// Store in global variable
	DB = db

	return db, nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		log.Println("🔌 Closing database connection...")
		return DB.Close()
	}
	return nil
}

// Ping checks if the database connection is alive
func Ping() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	return DB.Ping()
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}

// Stats returns database connection pool statistics
func Stats() sql.DBStats {
	if DB == nil {
		return sql.DBStats{}
	}
	return DB.Stats()
}

// HealthCheck performs a health check on the database
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	// Ping the database
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Check if we can execute a simple query
	var result int
	err := DB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected query result: %d", result)
	}

	return nil
}

// InitDB initializes the database with retry logic
func InitDB(config Config, maxRetries int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = Connect(config)
		if err == nil {
			return db, nil
		}

		log.Printf("⚠️  Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * 2 * time.Second
			log.Printf("⏳ Retrying in %v...", waitTime)
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// Transaction helper function
type TxFunc func(*sql.Tx) error

// WithTransaction executes a function within a database transaction
func WithTransaction(db *sql.DB, fn TxFunc) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %w", err, rbErr)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// LogConnectionStats logs the current connection pool statistics
func LogConnectionStats() {
	if DB == nil {
		log.Println("⚠️  Database connection is nil")
		return
	}

	stats := DB.Stats()
	log.Printf("📊 DB Stats - MaxOpen: %d, Open: %d, InUse: %d, Idle: %d",
		stats.MaxOpenConnections,
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
	)
}

// WaitForDB waits until the database is ready
func WaitForDB(db *sql.DB, timeout time.Duration) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-timeoutTimer.C:
			return fmt.Errorf("database did not become ready within %v", timeout)
		case <-ticker.C:
			if err := db.Ping(); err == nil {
				log.Println("✅ Database is ready!")
				return nil
			}
			log.Println("⏳ Waiting for database to be ready...")
		}
	}
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}
