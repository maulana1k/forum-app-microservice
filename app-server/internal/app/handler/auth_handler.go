package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/app/dto"
	"github.com/maulana1k/forum-app/internal/domain/service"
)

type AuthHandler struct {
	service.AuthService
	service.UserService
}

func NewAuthHandler(authService service.AuthService, userService service.UserService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		UserService: userService,
	}
}

// SignUp godoc
//
//	@Summary		Create a new user
//	@Description	Signup with username, email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.SignupRequest	true	"User info"
//	@Success		200		{object}	map[string]string
//	@Router			/v1/auth/signup [post]
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var body dto.SignupRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse body",
		})
	}

	user, token, err := h.AuthService.SignUp(body.Username, body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

// SignIn godoc
//
//	@Summary		Login user
//	@Description	Signin with email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		dto.SigninRequest	true	"User credentials"
//	@Success		200			{object}	map[string]string
//	@Router			/v1/auth/signin [post]
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {

	var body dto.SigninRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Cannot parse body",
		})
	}

	token, err := h.AuthService.SignIn(body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
