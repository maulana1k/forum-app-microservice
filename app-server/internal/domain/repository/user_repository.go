package repository

import (
	"github.com/google/uuid"
	"github.com/maulana1k/forum-app/internal/domain/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll() ([]models.User, error)
	CreateUserProfile(profile *models.User) error
	GetUserProfileByUserID(userID uuid.UUID) (*models.User, error)
	UpdateUserProfile(profile *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) CreateUserProfile(profile *models.User) error {
	return r.db.Create(profile).Error
}

func (r *userRepository) GetUserProfileByUserID(userID uuid.UUID) (*models.User, error) {
	var profile models.User
	err := r.db.Where("id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userRepository) UpdateUserProfile(profile *models.User) error {
	return r.db.Save(profile).Error
}
