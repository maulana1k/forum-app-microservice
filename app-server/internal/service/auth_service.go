package service

import (
	"errors"

	"github.com/maulana1k/forum-app/internal/models"
	"github.com/maulana1k/forum-app/internal/repository"
	"github.com/maulana1k/forum-app/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(username, email, password string) (string, error)
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

func (s *authService) SignUp(username, email, password string) (string, error) {
	// Check if email and username already exists
	if exists, _ := s.authRepo.IsEmailExists(email); exists {
		return "", errors.New("email already in use")
	}
	if exists, _ := s.authRepo.IsUsernameExists(username); exists {
		return "", errors.New("username already in use")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	// Create user model
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Save user to database
	if err := s.authRepo.CreateUser(user); err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
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
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
