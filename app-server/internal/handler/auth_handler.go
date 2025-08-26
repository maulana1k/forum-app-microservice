package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/dto"
	"github.com/maulana1k/forum-app/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
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
//	@Router			/auth/signup [post]
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var body dto.SignupRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse body",
		})
	}

	token, err := h.authService.SignUp(body.Username, body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
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
//	@Router			/auth/signin [post]
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {

	var body dto.SigninRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Cannot parse body",
		})
	}

	token, err := h.authService.SignIn(body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
