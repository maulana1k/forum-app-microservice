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
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/swagger"

	"github.com/maulana1k/forum-app/internal/app/container"
	"github.com/maulana1k/forum-app/internal/app/middleware"
	"github.com/maulana1k/forum-app/internal/app/routes"
	"github.com/maulana1k/forum-app/internal/config"
	"github.com/maulana1k/forum-app/internal/pkg/utils"
	"github.com/maulana1k/forum-app/internal/provider/broker"
	"github.com/maulana1k/forum-app/internal/provider/database"
	"github.com/maulana1k/forum-app/internal/provider/grpc"

	_ "github.com/maulana1k/forum-app/docs"
)

func Run() {

	cfg := config.LoadConfig()
	/*

	 */
	grpc := grpc.NewGRPCClient(cfg.GRPCAddress)
	db := database.NewDBInstance(cfg.DBUri)
	broker := broker.NewRabbitMQ(cfg.BrokerAddress)
	c := container.NewContainer(db.DB, grpc, broker)
	// prometheus := fiberprometheus.New("forum_app")
	defer db.Close()

	utils.InitLogger()

	app := fiber.New(cfg.AppConfig)
	// prometheus.RegisterAt(app, "/metrics")

	// app.Use(prometheus.Middleware)
	app.Use(middleware.LogrusMiddleware)
	app.Use(cors.New(cfg.CorsConfig))
	app.Use(pprof.New())

	routes.Register(app, c)

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/metrics", monitor.New(monitor.Config{Title: "Forum App Metrics", Refresh: time.Second}))
	// app.Get("/metrics", adaptor.HTTPHandlerFunc(promhttp.Handler().ServeHTTP))

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
