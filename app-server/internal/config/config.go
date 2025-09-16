package config

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

type Configuration struct {
	AppConfig     fiber.Config
	LoggerConfig  logger.Config
	CorsConfig    cors.Config
	DBUri         string
	GRPCAddress   string
	BrokerAddress string
}

func LoadConfig() *Configuration {
	v := viper.New()

	// Set defaults (only for fallback)
	v.SetDefault("POSTGRES_HOST", "localhost")
	v.SetDefault("POSTGRES_PORT", "5432")
	v.SetDefault("POSTGRES_USER", "dev")
	v.SetDefault("POSTGRES_PASSWORD", "dev")
	v.SetDefault("POSTGRES_DB", "forumdb")
	v.SetDefault("GRPC_ADDRESS", "localhost:50051")
	v.SetDefault("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")

	// Load .env file (environment-specific)
	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Override with environment variables if set
	v.AutomaticEnv()

	// Build Postgres DSN directly from env
	db_uri := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		v.GetString("POSTGRES_HOST"),
		v.GetString("POSTGRES_USER"),
		v.GetString("POSTGRES_PASSWORD"),
		v.GetString("POSTGRES_DB"),
		v.GetString("POSTGRES_PORT"),
	)

	return &Configuration{
		AppConfig: fiber.Config{
			Prefork:       false,
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  "Fiber",
			AppName:       "Forum App Server",
			ReadTimeout:   10 * time.Second,
			WriteTimeout:  10 * time.Second,
		},
		LoggerConfig: logger.Config{
			TimeZone:   "Local",
			TimeFormat: "15:04:05",
			Format:     "[${ip}]:${port} ${locals:latency} ${status} ${method} ${path}\n",
		},
		CorsConfig: cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
			AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		},
		DBUri:         db_uri,
		GRPCAddress:   v.GetString("GRPC_ADDRESS"),
		BrokerAddress: v.GetString("RABBITMQ_URI"),
	}
}
