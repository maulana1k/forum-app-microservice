package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/app/container"
	"github.com/maulana1k/forum-app/internal/app/handler"
)

func RegisterUserRoutes(app fiber.Router, c *container.Container) {
	userHandler := handler.NewUserHandler(c.UserService)

	v1 := app.Group("/v1/users")
	v1.Get("/", userHandler.GetAllUsers)
	v1.Get("/:id", userHandler.GetUserByID)
}
