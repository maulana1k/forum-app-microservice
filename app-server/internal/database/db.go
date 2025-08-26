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
var DBHost string
var DSN string

func Uptime() string {
	return time.Since(StartTime).String()
}

func Connect() {
	StartTime = time.Now()
	EnvMode := os.Getenv("DOCKER_ENV")
	if EnvMode == "true" {
		DBHost = "postgres"
	} else {
		DBHost = getEnv("POSTGRES_HOST", "localhost")
	}

	DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		DBHost,
		getEnv("POSTGRES_USER", "forum"),
		getEnv("POSTGRES_PASSWORD", "password"),
		getEnv("POSTGRES_DB", "forumdb"),
		getEnv("POSTGRES_PORT", "5432"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully with GORM")

	// Auto migrate all models
	tableMigration := []any{
		&models.User{},
		&models.UserProfile{},
		&models.Post{},
		&models.PostLikes{},
		&models.Replies{},
		&models.RepostByUser{},
		&models.BookmarkByUser{},
	}

	if err := DB.AutoMigrate(tableMigration...); err != nil {
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
