package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/app/container"
	"github.com/maulana1k/forum-app/internal/app/handler"
)

func RegisterRecommendationRoutes(app fiber.Router, c *container.Container, middleware fiber.Handler) {
	recHandler := handler.NewRecommendationHandler(c.RecommendationService)

	v1 := app.Group("/v1/recommendation")

	v1.Use(middleware)

	v1.Get("/posts", recHandler.GetRecommendedPosts)
}
