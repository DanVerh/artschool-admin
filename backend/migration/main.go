package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/DanVerh/artschool-admin/backend/migration/internal/application"
	"github.com/DanVerh/artschool-admin/backend/migration/cmd"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize application configuration
	app := application.Init()
	// Configure and run CLI
	cmd.RootCmd.AddCommand(cmd.UpCommand(app))
    if err := cmd.RootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}