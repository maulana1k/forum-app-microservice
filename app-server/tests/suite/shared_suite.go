package shared_suite

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maulana1k/forum-app/internal/app/container"
	"github.com/maulana1k/forum-app/internal/app/routes"
	"github.com/maulana1k/forum-app/internal/domain/models"
	"github.com/maulana1k/forum-app/internal/pkg/utils"
	"github.com/maulana1k/forum-app/tests/helper"
	"gorm.io/gorm"
)

type SharedSuite struct {
	App   *fiber.App
	Tx    *gorm.DB
	once  sync.Once
	Token string
}

var shared *SharedSuite

func GetSharedSuite() *SharedSuite {
	if shared == nil {
		shared = &SharedSuite{}
		shared.once.Do(func() {
			helper.Setup() // ensure test DB exists
			db := helper.GetTestDB()
			if db == nil {
				panic("test DB not initialized")
			}

			// Begin transaction for suite
			tx := db.Begin()

			// Setup Fiber app
			app := fiber.New()
			c := container.NewContainer(tx, nil)
			routes.Register(app, c)

			shared.App = app
			shared.Tx = tx
			shared.Token = CreateTestUser(shared.Tx)
		})
	}
	return shared
}

// Cleanup shared resources
func TeardownSharedSuite() {
	if shared != nil && shared.Tx != nil {
		shared.Tx.Rollback() // rollback all changes
	}
	helper.Teardown() // drop test DB
}

func CreateTestUser(tx *gorm.DB) string {
	user := models.User{
		ID:       uuid.New(),
		Email:    "testuser@example.com",
		Password: "test", // store hashed password if needed
		Username: "Test User",
	}

	if err := tx.Create(&user).Error; err != nil {
		panic("Failed to create test user: " + err.Error())
	}

	// generate JWT for this user
	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		panic("Failed to generate JWT: " + err.Error())
	}

	return token
}
