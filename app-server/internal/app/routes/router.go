package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/app/container"
	"github.com/maulana1k/forum-app/internal/pkg/utils"
)

func Register(app *fiber.App, c *container.Container) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", 301)
	})

	basePath := "/api"

	api := app.Group(basePath)

	api.Get("/health", HealthCheck)

	RegisterAuthRoutes(api, c)

	RegisterUserRoutes(api, c)

	RegisterRecommendationRoutes(api, c, utils.Protected())

	RegisterPublicPostRoutes(api, c)
	RegisterProtectedPostRoutes(api, c, utils.Protected())
}
