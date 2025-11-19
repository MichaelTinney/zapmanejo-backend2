package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pool to prevent "too many connections" errors
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// DigitalOcean basic PostgreSQL typically allows 22 connections
	// Reserve 2 for admin/monitoring, use up to 20 for the app
	sqlDB.SetMaxIdleConns(10)                  // Keep 10 connections ready
	sqlDB.SetMaxOpenConns(20)                  // Max 20 concurrent connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Recycle connections after 1 hour
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Close idle connections after 10 min

	// Verify connection is working
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("✓ Connected to DigitalOcean PostgreSQL with connection pooling")
	AutoMigrate()
}
