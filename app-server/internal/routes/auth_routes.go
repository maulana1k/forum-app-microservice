package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/database"
	"github.com/maulana1k/forum-app/internal/handler"
	"github.com/maulana1k/forum-app/internal/repository"
	"github.com/maulana1k/forum-app/internal/service"
)

func SetupAuthRoutes(app *fiber.App) {
	authRepo := repository.NewAuthRepository(database.DB)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	v1 := app.Group("/api/v1/auth")

	v1.Post("/signup", authHandler.SignUp)

	v1.Post("/signin", authHandler.SignIn)
}
