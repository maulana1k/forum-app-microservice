package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/handler"
)

func SetupAuthRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Post("/signup", handler.SignUp)

	v1.Post("/signin", handler.SignIn)
}
