package helper

import (
	"fmt"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/maulana1k/forum-app/internal/provider/database"
	"gorm.io/gorm"
)

const (
	testDBName = "testdb"
	baseDBName = "postgres"
)

var TestDB database.DB

func GetTestDB() *gorm.DB {
	if TestDB.DB == nil {
		TestDB = ConnectToDB(testDBName)
	}
	return TestDB.DB
}

func Setup() {

	// Connect to base postgres db to create test database
	baseDB := ConnectToDB(baseDBName)
	// defer baseDB.Close()

	// Drop test database if exists and create fresh one
	result := baseDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if result.Error != nil {
		log.Fatalf("Failed to drop test database: %v", result.Error)
	}

	result = baseDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if result.Error != nil {
		log.Fatalf("Failed to drop test database: %v", result.Error)
	}

	// Connect to test database
	TestDB = ConnectToDB(testDBName)

	if TestDB.DB == nil {
		log.Fatalf("testDB.DB is nil after Setup")
	}
}

func Teardown() {

	baseDB := ConnectToDB(baseDBName)
	defer baseDB.Close()

	if err := baseDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName)).Error; err != nil {
		log.Printf("Failed to drop test database during teardown: %v", err)
	}

	if TestDB.DB != nil {
		TestDB.Close()
	}
}

func ConnectToDB(dbName string) database.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "dev")
	password := getEnv("DB_PASSWORD", "dev")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslmode)

	db := database.NewDBInstance(dsn)

	return *db
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
