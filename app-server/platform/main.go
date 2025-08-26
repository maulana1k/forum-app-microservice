package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

var DBHost string
var DSN string

func main() {
	// Accept migration action: up, down, status, reset
	action := flag.String("action", "up", "Migration action: up, down, status, reset")
	flag.Parse()

	InitConfiguration()

	// Open database
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://platform/migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Execute migration based on action
	switch *action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("Migration up completed successfully!")

	case "down":
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("Migration down completed successfully!")

	case "reset":
		if err := m.Drop(); err != nil {
			log.Fatalf("Migration reset failed: %v", err)
		}
		fmt.Println("Migration reset completed successfully!")

	case "status":
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Migration version: %v, dirty: %v\n", version, dirty)

	default:
		fmt.Println("Invalid action. Use 'up', 'down', 'status', or 'reset'.")
		os.Exit(1)
	}
}

func InitConfiguration() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	EnvMode := os.Getenv("DOCKER_ENV")
	if EnvMode == "true" {
		DBHost = "postgres"
	} else {
		DBHost = getEnv("POSTGRES_HOST", "localhost")
	}

	DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		DBHost,
		getEnv("POSTGRES_USER", "test"),
		getEnv("POSTGRES_PASSWORD", "test"),
		getEnv("POSTGRES_DB", "forumdb"),
		getEnv("POSTGRES_PORT", "5432"),
	)
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
