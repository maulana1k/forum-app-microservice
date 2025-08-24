package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/maulana1k/forum-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var StartTime time.Time

func Uptime() string {
	return time.Since(StartTime).String()
}

func Connect() {
	StartTime = time.Now()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_USER", "forum"),
		getEnv("POSTGRES_PASSWORD", "password"),
		getEnv("POSTGRES_DB", "forumdb"),
		getEnv("POSTGRES_PORT", "5432"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully with GORM")

	// Auto migrate User model
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}
}

func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Error getting underlying sql.DB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
