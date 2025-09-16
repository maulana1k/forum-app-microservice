package database

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/maulana1k/forum-app/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	StartTime time.Time
	DSN       string
}

var (
	instance *DB
	once     sync.Once
)

func GetDB() (*gorm.DB, error) {
	if instance == nil {
		return nil, errors.New("database not initialized. call database first")
	}
	return instance.DB, nil
}

func NewDBInstance(DSN string) *DB {
	once.Do(func() {
		instance = &DB{DSN: DSN}
		instance.Connect()
	})
	return instance
}

func Ping() (time.Duration, error) {
	if instance == nil || instance.DB == nil {
		return 0, errors.New("database not initialized")
	}

	sqlDB, err := instance.DB.DB()
	if err != nil {
		return 0, err
	}

	start := time.Now()
	if err := sqlDB.Ping(); err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

func Uptime() string {
	if instance == nil {
		return "DB not initialized"
	}
	return time.Since(instance.StartTime).String()
}

func NewDB(DSN string) *DB {
	return &DB{DSN: DSN}
}

func (db *DB) Connect() *gorm.DB {
	db.StartTime = time.Now()

	var err error
	db.DB, err = gorm.Open(postgres.Open(db.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully with GORM")

	// Auto migrate all models
	tableMigration := []any{
		&models.User{},
		&models.Post{},
		&models.Replies{},
		&models.PostInteractions{},
	}

	if err := db.DB.AutoMigrate(tableMigration...); err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}

	return db.DB
}

func (db *DB) Close() {
	if db.DB != nil {
		sqlDB, err := db.DB.DB()
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
