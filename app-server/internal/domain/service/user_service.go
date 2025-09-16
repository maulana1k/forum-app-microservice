package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/maulana1k/forum-app/internal/domain/models"
	"github.com/maulana1k/forum-app/internal/domain/repository"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	CreateUserProfile(userID uuid.UUID, username string) error
	GetUserProfile(userID uuid.UUID) (*models.User, error)
	UpdateUserProfile(userID uuid.UUID, displayName, bio, location, avatarURL string) error
	// FollowUser(userID uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) CreateUserProfile(userID uuid.UUID, username string) error {
	// Check if profile already exists
	existingProfile, err := s.userRepo.GetUserProfileByUserID(userID)
	if err == nil && existingProfile != nil {
		return errors.New("user profile already exists")
	}

	profile := &models.User{
		ID:          userID,
		DisplayName: username, // Default displayname is same as username
		AvatarURL:   "",       // Empty by default
		Bio:         "",       // Empty by default
		Location:    "",       // Empty by default
	}

	return s.userRepo.CreateUserProfile(profile)
}

func (s *userService) GetUserProfile(userID uuid.UUID) (*models.User, error) {
	profile, err := s.userRepo.GetUserProfileByUserID(userID)
	if err != nil {
		return nil, errors.New("user profile not found")
	}
	return profile, nil
}

func (s *userService) UpdateUserProfile(userID uuid.UUID, displayName, bio, location, avatarURL string) error {
	profile, err := s.userRepo.GetUserProfileByUserID(userID)
	if err != nil {
		return errors.New("user profile not found")
	}

	// Update fields only if they are provided
	if displayName != "" {
		profile.DisplayName = displayName
	}
	if bio != "" {
		profile.Bio = bio
	}
	if location != "" {
		profile.Location = location
	}
	if avatarURL != "" {
		profile.AvatarURL = avatarURL
	}

	return s.userRepo.UpdateUserProfile(profile)
}
