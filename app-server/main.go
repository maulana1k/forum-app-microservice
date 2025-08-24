package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/maulana1k/forum-app/cmd/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	server.Run()
}
