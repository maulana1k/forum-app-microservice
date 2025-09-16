package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/app/container"
	"github.com/maulana1k/forum-app/internal/app/handler"
)

// RegisterPublicPostRoutes registers only public post routes (no auth required)
func RegisterPublicPostRoutes(api fiber.Router, c *container.Container) {
	postHandler := handler.NewPostHandler(c.PostService)

	v1 := api.Group("/v1/posts")
	v1.Get("/", postHandler.GetAllPosts)
	v1.Get("/:id", postHandler.GetPostByID)
	v1.Get("/user/:id", postHandler.GetUserPosts)
}

// RegisterProtectedPostRoutes registers routes that require authentication
func RegisterProtectedPostRoutes(api fiber.Router, c *container.Container, middleware fiber.Handler) {
	postHandler := handler.NewPostHandler(c.PostService)

	v1 := api.Group("/v1/posts")
	v1.Use(middleware) // apply middleware to all protected endpoints

	v1.Post("/", postHandler.CreatePost)
	v1.Post("/:id/like", postHandler.LikePost)
	// v1.Post("/:id/bookmark", postHandler.BookmarkPost)
	v1.Put("/:id", postHandler.UpdatePost)
	v1.Delete("/:id", postHandler.DeletePost)
	v1.Delete("/:id/unlike", postHandler.UnlikePost)
	// v1.Delete("/:id/unbookmark", postHandler.UnbookmarkPost)
}
