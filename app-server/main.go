package main

import (
	"github.com/maulana1k/forum-app/cmd/server"
)

// Package main is the entry point for the Forum App API server.
//
// @title           Forum App API
// @version         1.0
// @description     This is a simple RESTful API for a forum application.
//
//	It allows users to sign in, create posts, read posts, update them, and delete them.
//	JWT authentication is used to secure protected routes.
//
// @termsOfService  http://example.com/terms/
//
// @contact.name    API Support
// @contact.url     https://github.com/maulana1k/forum-app
// @contact.email   support@example.com
//
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
//
// @host            localhost:8080
// @BasePath        /api
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
//
// @schemes         http
func main() {
	// env := os.Getenv("APP_ENV")
	// if env == "" {
	// 	env = "development" // default
	// }

	// envFile := ".env." + env
	// if err := godotenv.Load(envFile); err != nil {
	// 	log.Printf("No %s file found, using system env", envFile)
	// } else {
	// 	log.Printf("Loaded environment from %s", envFile)
	// }

	server.Run()
}
