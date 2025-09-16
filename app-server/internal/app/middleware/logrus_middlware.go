package middleware

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/pkg/utils"
	"github.com/sirupsen/logrus"
)

func LogrusMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Generate request ID
	requestID := rand.Intn(1000)

	// Create logger with context
	logger := utils.Logger.WithFields(logrus.Fields{
		"component":  "web-server",
		"method":     c.Method(),
		"path":       c.Path(),
		"request_id": requestID,
	})

	// Skip logging for static assets
	if shouldSkipLog(c.Path()) {
		return c.Next()
	}

	logger.Info("request received")

	err := c.Next()
	latency := time.Since(start)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"latency": latency.Round(time.Millisecond).String(),
			"status":  c.Response().StatusCode(),
		}).Error("response error")
	} else {
		logger.WithFields(logrus.Fields{
			"latency": latency.Round(time.Millisecond).String(),
			"status":  c.Response().StatusCode(),
		}).Info("response success")
	}

	return err
}

func shouldSkipLog(path string) bool {
	skipExtensions := []string{".css", ".js", ".png", ".jpg", ".ico", ".svg"}
	for _, ext := range skipExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}
