package common

import (
	"database/sql"
	"fmt"
	"os"
)

// ConnectDB menginisialisasi dan mengembalikan koneksi database
func ConnectDB() (*sql.DB, error) {
	// menggunakan os.Getenv untuk mendapatkan konfigurasi dari environment variables
	dbHost := getEnvOrDefault("DB_HOST", "127.0.0.1")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "akbar123")
	dbName := getEnvOrDefault("DB_NAME", "blogging_platform")
	
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Cek koneksi
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// getEnvOrDefault mendapatkan environment variable atau default value
func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}