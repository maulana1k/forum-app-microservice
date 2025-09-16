package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/app/container"
	"github.com/maulana1k/forum-app/internal/app/handler"
)

func RegisterAuthRoutes(api fiber.Router, c *container.Container) {
	authHandler := handler.NewAuthHandler(c.AuthService, c.UserService)

	authV1 := api.Group("/v1/auth")

	authV1.Post("/signup", authHandler.SignUp)

	authV1.Post("/signin", authHandler.SignIn)
}
