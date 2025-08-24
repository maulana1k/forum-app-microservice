package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/database"
	"github.com/maulana1k/forum-app/internal/models"
	"github.com/maulana1k/forum-app/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// SignUp godoc
// @Summary Create a new user
// @Description Signup with username, email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.SignupRequest true "User info"
// @Success 200 {object} map[string]string
// @Router /signup [post]
func SignUp(c *fiber.Ctx) error {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

	user := models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, _ := utils.GenerateJWT(user.ID)
	return c.JSON(fiber.Map{"token": token})
}

// SignIn godoc
// @Summary Login user
// @Description Signin with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body dto.SigninRequest true "User credentials"
// @Success 200 {object} map[string]string
// @Router /signin [post]
func SignIn(c *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, _ := utils.GenerateJWT(user.ID)
	return c.JSON(fiber.Map{"token": token})
}
