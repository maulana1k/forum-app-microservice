package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/pkg/utils"
	"go.uber.org/zap"
)

func ZapMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	logger := utils.GetLogger()

	// Process request
	err := c.Next()

	// Log the request
	latency := time.Since(start)

	logger.Info("HTTP Request",
		zap.String("ip", c.IP()),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.Int("status", c.Response().StatusCode()),
		zap.Duration("latency", latency),
		zap.String("user_agent", c.Get("User-Agent")),
	)

	if err != nil {
		logger.Error("Request failed",
			zap.String("ip", c.IP()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Error(err),
		)
	}

	return err
}
