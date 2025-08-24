package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/maulana1k/forum-app/internal/configs"
	"github.com/maulana1k/forum-app/internal/database"
	"github.com/maulana1k/forum-app/internal/routes"

	_ "github.com/maulana1k/forum-app/docs"
)

// @title Forum App API
// @version 1.0
// @description This is a simple forum app server API.
// @host localhost:8080
// @BasePath /api/v1
func Run() {

	configs.LoadConfig()

	database.Connect()
	defer database.Close()

	app := fiber.New(configs.AppConfig)

	app.Use(logger.New(configs.LoggerConfig))
	app.Use(cors.New(configs.CorsConfig))

	routes.SetupPublicRoutes(app)
	routes.SetupAuthRoutes(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Printf("Server stopped: %v\n", err)
		}
	}()

	gracefulShutdown(app, 10*time.Second)
}

func gracefulShutdown(app *fiber.App, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Error during shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}
