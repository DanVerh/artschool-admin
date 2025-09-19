package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/DanVerh/artschool-admin/backend/migration/application"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Start application
	app := application.New()
	app.Start()
}