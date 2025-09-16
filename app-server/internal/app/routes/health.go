package routes

import (
	"log"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/provider/database"
)

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Checks if the server and database are running
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/health [get]
func HealthCheck(c *fiber.Ctx) error {
	// Default status
	status := "ok"
	dbStatus := "ok"
	var dbLatency string

	// Measure DB ping latency
	latency, err := database.Ping()
	if err != nil {
		status = "error"
		dbStatus = "down"
		log.Printf("Database ping failed: %v", err)
	} else {
		dbLatency = latency.String()
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
