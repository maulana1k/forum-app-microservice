package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/database"
	"github.com/maulana1k/forum-app/internal/handler"
	"github.com/maulana1k/forum-app/internal/repository"
	"github.com/maulana1k/forum-app/internal/service"
	"github.com/maulana1k/forum-app/internal/utils"
)

func SetupPostRoutes(app *fiber.App) {
	postRepo := repository.NewPostRepository(database.DB)
	postService := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postService)

	v1 := app.Group("/api/v1/posts")

	v1.Get("/", postHandler.GetAllPosts)
	v1.Get("/:id", postHandler.GetPostByID)
	v1.Get("/user/:id", postHandler.GetUserPosts)

	v1.Use(utils.Protected())

	v1.Post("/", postHandler.CreatePost)
	v1.Post("/:id/like", postHandler.LikePost)
	v1.Post("/:id/repost", postHandler.RepostPost)
	v1.Post("/:id/bookmark", postHandler.BookmarkPost)

	v1.Put("/:id", postHandler.UpdatePost)

	v1.Delete("/:id", postHandler.DeletePost)
	v1.Delete("/:id/unlike", postHandler.UnlikePost)
	v1.Delete("/:id/unreport", postHandler.UnrepostPost)
	v1.Delete("/:id/unbookmark", postHandler.UnbookmarkPost)
}
