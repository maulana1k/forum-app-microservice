package configs

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	AppConfig    fiber.Config
	LoggerConfig logger.Config
	CorsConfig   cors.Config
)

func LoadConfig() {
	AppConfig = fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Forum App Server",
	}

	LoggerConfig = logger.Config{
		Output:     nil, // default to os.Stdout
		TimeZone:   "Local",
		TimeFormat: "15:04:05",
		Format:     "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}

	CorsConfig = cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}

}
