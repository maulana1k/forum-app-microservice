package routes

import (
	"log"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/database"
)

func SetupPublicRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Hello": "Dunia!"})
	})

	app.Get("/health", HealthCheck)

}

// HealthCheck godoc
// @Summary Health check
// @Description Checks if the server and database are running
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *fiber.Ctx) error {
	// Default status
	status := "ok"
	dbStatus := "ok"
	var dbLatency string

	// Measure DB ping latency
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Printf("Error getting DB: %v", err)
		status = "error"
		dbStatus = "unavailable"
	} else {
		start := time.Now()
		if err := sqlDB.Ping(); err != nil {
			log.Printf("Database ping failed: %v", err)
			status = "error"
			dbStatus = "down"
		} else {
			dbLatency = time.Since(start).String()
		}
	}

	// Memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	resp := fiber.Map{
		"status":       status,
		"server":       "ok",
		"database":     dbStatus,
		"db_latency":   dbLatency,
		"uptime":       database.Uptime(),
		"memory_alloc": m.Alloc,
		"memory_total": m.TotalAlloc,
		"memory_sys":   m.Sys,
		"goroutines":   runtime.NumGoroutine(),
	}

	if status == "error" {
		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}

	return c.JSON(resp)
}
