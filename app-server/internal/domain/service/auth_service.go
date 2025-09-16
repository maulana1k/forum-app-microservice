package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/maulana1k/forum-app/internal/domain/models"
	"github.com/maulana1k/forum-app/internal/domain/repository"
	"github.com/maulana1k/forum-app/internal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(username, email, password string) (models.User, string, error)
	SignIn(email, password string) (string, error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (s *authService) SignUp(username, email, password string) (models.User, string, error) {
	// Trim inputs
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	// Validate required fields
	if username == "" || email == "" || password == "" {
		return models.User{}, "", errors.New("username, email, and password cannot be empty")
	}

	// Check if email or username already exists
	if exists, _ := s.authRepo.IsEmailExists(email); exists {
		return models.User{}, "", errors.New("email already in use")
	}
	if exists, _ := s.authRepo.IsUsernameExists(username); exists {
		return models.User{}, "", errors.New("username already in use")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return models.User{}, "", err
	}

	// Create user model with safe defaults
	user := models.User{
		ID:          uuid.New(),
		Username:    username,
		DisplayName: username, // default display name
		Email:       email,
		Password:    string(hashedPassword),
		Role:        "user", // default role
		AvatarURL:   "",
		Bio:         "",
		Location:    "",
	}

	if err := s.authRepo.CreateUser(&user); err != nil {
		return models.User{}, "", err
	}

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return models.User{}, "", err
	}

	return user, token, nil
}

func (s *authService) SignIn(email, password string) (string, error) {
	// Get user by email
	user, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}
