package helper

import (
	"github.com/google/uuid"
	"github.com/maulana1k/forum-app/internal/domain/models"
	"github.com/maulana1k/forum-app/internal/pkg/utils"
	"gorm.io/gorm"
)

// CreateTestUser inserts a test user and returns a JWT token for it
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
